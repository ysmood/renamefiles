package main

import (
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/ysmood/kit"
	"github.com/ysmood/renamefiles/lib"
)

func main() {
	app := kit.TasksNew("renamefiles", "if the target path doesn't exist it will be auto created")
	app.Version("v0.3.0")
	logPath := app.Flag("file", "the path of the log file").Short('f').Default(".renamefiles.log").String()
	noLog := app.Flag("no-log", "don't generate log file").Short('n').Bool()

	kit.Tasks().App(app).Add(
		kit.Task("do", "use pattern to rename").Init(func(cmd kit.TaskCmd) func() {
			cmd.Default()

			match := cmd.Flag("match", "glob pattern for files to rename").Short('m').Default("*").String()
			regStr := cmd.Flag(
				"key",
				"regex to match the sortable key of the names, trap has priority. "+
					"By default, it will auto analyze files and generate a pattern.",
			).Short('k').Regexp()
			template := cmd.Flag(
				"template",
				"template to move the files to. "+
					"Use function `index` such as {{index 2}} to change the start index",
			).Short('t').Default("{{key}}{{ext}}").String()
			prefix := cmd.Flag("prefix", "prefix to each name").Short('p').Default("").String()
			minPad := cmd.Flag("min-padding", "minimal zero padding for index").Default("2").Int64()
			yes := cmd.Flag("yes", "no confirmation").Bool()

			return func() {
				tasks := plan(*match, *regStr, *template, *prefix, *minPad)

				if len(tasks) == 0 {
					fmt.Println("Nothing to rename")
					return
				}

				var confirm string
				if *yes {
					confirm = "yes"
				} else {
					fmt.Println("Sure to rename? Press enter to continue (CTRL-C to abort)")
					fmt.Scanln()
					confirm = "yes"
				}
				if confirm == "yes" {
					move(tasks)

					if !*noLog {
						log(*logPath, tasks)
					}
				}
			}
		}),
		kit.Task("revert", "revert the batch operation").Run(func() {
			revert(*logPath)
		}),
	).Do()
}

type task struct {
	From string
	Tmp  string
	To   string
}

func plan(match string, reg *regexp.Regexp, template, prefix string, minPad int64) []task {
	list := kit.Walk(match).Sort().MustList()
	tasks := []task{}
	padLen := int64(math.Ceil(math.Log10(float64(len(list)))))

	if padLen < minPad {
		padLen = minPad
	}

	indexFn := genIndexFn(padLen)

	if reg == nil {
		reg = lib.AutoPattern(baseNames(list))
	}

	if reg.String() == "" {
		return tasks
	}

	for _, p := range list {
		m := reg.FindStringSubmatch(filepath.Base(p))

		if m == nil {
			continue
		}

		key := m[0]
		if len(m) > 1 {
			key = m[1]
		}

		index, err := strconv.ParseInt(key, 10, 64)
		kit.E(err)
		key = formatIndex(index, padLen)

		to, err := filepath.Abs(prefix + kit.S(template,
			"key", str(key),
			"ext", str(filepath.Ext(p)),
			"index", indexFn,
		))
		kit.E(err)

		abs, err := filepath.Abs(".")
		kit.E(err)
		relFrom, err := filepath.Rel(abs, p)
		kit.E(err)
		relTo, err := filepath.Rel(abs, to)
		kit.E(err)

		fmt.Println(relFrom, kit.C("->", "cyan"), relTo)
		tasks = append(tasks, task{From: p, Tmp: kit.RandString(16), To: to})
	}

	return tasks
}

func baseNames(list []string) []string {
	names := []string{}
	for _, p := range list {
		names = append(names, filepath.Base(p))
	}
	return names
}

func formatIndex(index, padLen int64) string {
	return fmt.Sprintf("%0"+strconv.FormatInt(padLen, 10)+"d", index)
}

func genIndexFn(padLen int64) func(...int64) string {
	index := int64(1)
	return func(min ...int64) string {
		if len(min) > 0 && index < min[0] {
			index = min[0]
		}
		out := formatIndex(index, padLen)
		index++
		return out
	}
}

func str(s string) func() string {
	return func() string {
		return s
	}
}

func move(tasks []task) {
	for _, t := range tasks {
		kit.E(kit.Move(t.From, t.Tmp, nil))
	}
	for _, t := range tasks {
		kit.E(kit.Move(t.Tmp, t.To, nil))
	}
}

func log(path string, data interface{}) {
	kit.E(kit.OutputFile(path, data, nil))
}

func revert(path string) {
	var tasks []task
	kit.E(kit.ReadJSON(path, &tasks))

	for _, t := range tasks {
		kit.E(kit.Move(t.To, t.From, nil))
	}
}

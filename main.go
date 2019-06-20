package main

import (
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"strconv"

	kit "github.com/ysmood/gokit"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("renamefiles", "if the target path doesn't exist it will be auto created")
	logPath := app.Flag("file", "the path of the log file").Short('f').Default(".renamefiles.log").String()
	noLog := app.Flag("no-log", "don't generate log file").Short('n').Bool()

	kit.Tasks().App(app).Add(
		kit.Task("do", "rename").Init(func(cmd kit.TaskCmd) func() {
			cmd.Default()

			match := cmd.Flag("match", "glob pattern for files to rename").Short('m').Default("*").String()
			regStr := cmd.Flag("key", "regex to match the sortable key of the names").Short('k').Default(`\d+`).String()
			template := cmd.Flag("template", "template to move the files to").Short('t').Default("{{key}}{{ext}}").String()
			prefix := cmd.Flag("prefix", "prefix to each name").Short('p').Default("").String()
			yes := cmd.Flag("yes", "no confirmation").Bool()

			return func() {
				tasks := plan(*match, *regStr, *template, *prefix)

				if len(tasks) == 0 {
					fmt.Println("Nothing to rename")
					return
				}

				var confirm string
				if *yes {
					confirm = "yes"
				} else {
					fmt.Println("Sure to rename: [no]/yes")
					fmt.Scanln(&confirm)
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
	To   string
}

func plan(match, regStr, template, prefix string) []task {
	list := kit.Walk(match).MustList()
	reg := regexp.MustCompile(regStr)
	tasks := []task{}
	padLen := int64(math.Ceil(math.Log10(float64(len(list)))))

	for _, p := range list {
		m := reg.FindStringSubmatch(filepath.Base(p))

		if m == nil {
			continue
		}

		key := m[0]
		if len(m) > 1 {
			key = m[1]
		}

		index, _ := strconv.ParseInt(key, 10, 64)
		key = fmt.Sprintf("%0"+strconv.FormatInt(padLen, 10)+"d", index)

		to, _ := filepath.Abs(prefix + kit.S(template, "key", str(key), "ext", str(filepath.Ext(p))))

		abs, _ := filepath.Abs(".")
		relFrom, _ := filepath.Rel(abs, p)
		relTo, _ := filepath.Rel(abs, to)

		fmt.Println(relFrom, kit.C("->", "cyan"), relTo)
		tasks = append(tasks, task{From: p, To: to})
	}

	return tasks
}

func str(s string) func() string {
	return func() string {
		return s
	}
}

func move(tasks []task) {
	for _, t := range tasks {
		kit.E(kit.Move(t.From, t.To, nil))
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

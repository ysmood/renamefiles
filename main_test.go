package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	kit "github.com/ysmood/gokit"
)

func TestBasic(t *testing.T) {
	p := "tmp/" + kit.RandString(16)

	_ = kit.OutputFile(p+"/test-01-a.txt", "", nil)
	_ = kit.OutputFile(p+"/test-02-a.txt", "", nil)
	_ = kit.OutputFile(p+"/test-03-a.txt", "", nil)

	kit.Exec(
		"go", "run", ".",
		"--yes",
		"-f"+p+"/.log",
		`-k-(\d+)-`, "-m"+p+"/*",
		"-t"+p+"/ok-{{key}}.txt",
	).MustDo()

	assert.True(t, kit.FileExists(p+"/ok-01.txt"))
	assert.True(t, kit.FileExists(p+"/ok-02.txt"))
	assert.True(t, kit.FileExists(p+"/ok-03.txt"))

	kit.Exec(
		"go", "run", ".",
		"revert",
		"-f"+p+"/.log",
	).MustDo()

	assert.True(t, kit.FileExists(p+"/test-01-a.txt"))
	assert.True(t, kit.FileExists(p+"/test-02-a.txt"))
	assert.True(t, kit.FileExists(p+"/test-03-a.txt"))
}

func TestNameShifting(t *testing.T) {
	p := "tmp/" + kit.RandString(16)

	_ = kit.OutputFile(p+"/01", "", nil)
	_ = kit.OutputFile(p+"/02", "", nil)

	kit.Exec(
		"go", "run", ".",
		"--yes",
		"-f"+p+"/.log",
		"-m"+p+"/*",
		"-t"+p+"/{{index 2}}",
	).MustDo()

	assert.True(t, kit.FileExists(p+"/02"))
	assert.True(t, kit.FileExists(p+"/03"))
}

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ysmood/kit"
)

func TestBasic(t *testing.T) {
	p := "tmp/" + kit.RandString(16)

	_ = kit.OutputFile(p+"/test-01-a.txt", "", nil)
	_ = kit.OutputFile(p+"/test-02-a.txt", "", nil)
	_ = kit.OutputFile(p+"/test-03-a.txt", "", nil)

	os.Args = []string{
		"",
		"--yes",
		"-f" + p + "/.log",
		`-k-(\d+)-`, "-m" + p + "/*",
		"-t" + p + "/ok-{{key}}.txt",
	}

	main()

	assert.FileExists(t, p+"/ok-01.txt")
	assert.FileExists(t, p+"/ok-02.txt")
	assert.FileExists(t, p+"/ok-03.txt")

	os.Args = []string{
		"",
		"revert",
		"-f" + p + "/.log",
	}

	main()

	assert.FileExists(t, p+"/test-01-a.txt")
	assert.FileExists(t, p+"/test-02-a.txt")
	assert.FileExists(t, p+"/test-03-a.txt")
}

func TestAutoPattern(t *testing.T) {
	p := "tmp/" + kit.RandString(16)

	_ = kit.OutputFile(p+"/01-01-a.txt", "", nil)
	_ = kit.OutputFile(p+"/01-02-a.txt", "", nil)
	_ = kit.OutputFile(p+"/01-03-a.txt", "", nil)

	os.Args = []string{
		"",
		"--yes",
		"--no-log",
		"-m" + p + "/*",
		"-t" + p + "/ok-{{key}}.txt",
	}

	main()

	assert.FileExists(t, p+"/ok-01.txt")
	assert.FileExists(t, p+"/ok-02.txt")
	assert.FileExists(t, p+"/ok-03.txt")
}

func TestNameShifting(t *testing.T) {
	p := "tmp/" + kit.RandString(16)

	_ = kit.OutputFile(p+"/01", "", nil)
	_ = kit.OutputFile(p+"/02", "", nil)

	os.Args = []string{
		"",
		"--yes",
		"-f" + p + "/.log",
		"-m" + p + "/*",
		"-t" + p + "/{{index 2}}",
	}

	main()

	assert.FileExists(t, p+"/02")
	assert.FileExists(t, p+"/03")
}

func TestNothingTodo(t *testing.T) {
	os.Args = []string{
		"",
	}

	main()

	os.Args = []string{
		"",
		"--key=" + kit.RandString(16),
	}

	main()
}

func TestPrompt(t *testing.T) {
	p := "tmp/" + kit.RandString(16)

	_ = kit.OutputFile(p+"/01", "", nil)

	os.Args = []string{
		"",
		"-m", p + "/*",
		"-t" + p + "/{{key}}",
	}

	main()
}

# renamefiles

[![GoDoc](https://godoc.org/github.com/ysmood/renamefiles?status.svg)](http://godoc.org/github.com/ysmood/renamefiles)
[![Build Status](https://travis-ci.org/ysmood/renamefiles.svg?branch=master)](https://travis-ci.org/ysmood/renamefiles)
[![codecov](https://codecov.io/gh/ysmood/renamefiles/branch/master/graph/badge.svg)](https://codecov.io/gh/ysmood/renamefiles)
[![goreport](https://goreportcard.com/badge/github.com/ysmood/renamefiles)](https://goreportcard.com/report/github.com/ysmood/renamefiles)

A CLI tool to safely batch rename files with sane defaults. Usually, you don't have to pass any arguments to make it work.
The quality of the statistic algorithm is defined by how hard to mock some file names to make the defaults not work as expected.

## Features

- Auto-detect the pattern with a [statistic based algorithm](lib/auto_pattern.go)
- Preview the renames before apply
- Command to revert the renames

## Install

Shell command `curl -L https://git.io/fjaxx | repo=ysmood/renamefiles sh`

or build from source code:

```bash
go get github.com/ysmood/renamefiles
```

For Windows go to the release page and download the binary.

## Usage

```bash
renamefiles --help-long

# usually, command without argument is enough
renamefiles
```

The tool will display a preview and ask if you want to apply the batch operation.

After the rename a log file will be created, you can use it to revert batch operation in case of mistakes.

For example, file list:

```txt
2001-abc-01.ass
2001-abc-02.ass
2001-abc-10.ass
readme.txt
```

The auto rename with `renamefiles --template "2013-{{key}}.ass"` will be:

```txt
2001-abc-01.ass -> 2013-01.ass
2001-abc-02.ass -> 2013-02.ass
2001-abc-10.ass -> 2013-10.ass
```

Because `readme.txt` is not similar to other files, it won't be touched.

## Dev

Release to github:

```bash
go get github.com/ysmood/kit/cmd/godev

godev build -d
```

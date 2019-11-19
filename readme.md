[![GoDoc](https://godoc.org/github.com/ysmood/renamefiles?status.svg)](http://godoc.org/github.com/ysmood/renamefiles)
[![Build Status](https://travis-ci.org/ysmood/renamefiles.svg?branch=master)](https://travis-ci.org/ysmood/renamefiles)
[![codecov](https://codecov.io/gh/ysmood/renamefiles/branch/master/graph/badge.svg)](https://codecov.io/gh/ysmood/renamefiles)
[![goreport](https://goreportcard.com/badge/github.com/ysmood/renamefiles)](https://goreportcard.com/report/github.com/ysmood/renamefiles)

# rename files

A cli tool to safely batch rename files with sane defaults. Most times you don't have to pass any arguments to make it work.
The quality of the statistic algorithm is defined by how hard to mock some file names to make the defaults doesn't work as expected.

## Features

- Auto detect the pattern with a [statistic based algorithm](lib/auto_pattern.go)
- Preview the changes before apply
- Command to revert the changes

## Install

Shell command `curl -L https://git.io/fjaxx | repo=ysmood/renamefiles sh`

or use golang:

```bash
go get github.com/ysmood/renamefiles
```

For Windows goto the release page and download the binary.

## Usage

```bash
renamefiles --help-long

# normally, command without arguments is enough
renamefiles
```

The tool will display a preview and ask if you want to apply the batch operation.

A log file will be created, you can use it to revert batch operation in case mistakes.

## Dev

Release to github:

```bash
go get github.com/ysmood/gokit/cmd/godev

godev build -d
```

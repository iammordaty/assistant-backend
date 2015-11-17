package helper

import (
	"bytes"
	"os/exec"
    "strings"
)

func RunCommand(cmd string, ch chan CommandResult) {
    command := &Command{"/bin/sh", []string{"-c", cmd}, ch}

    go command.Run()
}

type Command struct {
	Command     string
	Args        []string
	Result      chan CommandResult
}

type CommandResult struct {
	Stdout      string
    Stderr      string
    Error       error
}

func (w *Command) Run() {
    var stdout, stderr bytes.Buffer

    defer close(w.Result)

    cmd := exec.Command(w.Command, w.Args...)
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()

    wr := new (CommandResult)
    wr.Stdout = strings.Trim(stdout.String(), "\n ")
    wr.Stderr = strings.Trim(stderr.String(), "\n ")
    wr.Error = err

    w.Result <- *wr
}
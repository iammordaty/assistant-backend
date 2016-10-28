package main

import (
    "bytes"
    "os/exec"
    "strings"
)

// http://stackoverflow.com/questions/40247726/go-execute-a-bash-command-n-times-using-goroutines-and-store-print-its-resul
// http://stackoverflow.com/questions/25283600/go-routine-for-cmd-exec-but-with-errorcode

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

func RunCommand(cmd string, ch chan CommandResult) {
    command := &Command{"/bin/sh", []string{"-c", cmd}, ch}

    go command.Run()
}

func (c *Command) Run() {
    var stdout, stderr bytes.Buffer

    defer close(c.Result)

    cmd := exec.Command(c.Command, c.Args...)
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr

    err := cmd.Run()

    cr := &CommandResult{}
    cr.Stdout = strings.Trim(stdout.String(), "\n ")
    cr.Stderr = strings.Trim(stderr.String(), "\n ")
    cr.Error = err

    c.Result <- *cr
}

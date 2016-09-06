package main

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
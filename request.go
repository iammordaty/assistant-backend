package main

type Payload struct {
    Pathname    string    `json:"pathname"`
    InitialKey  []string  `json:"initial_key"`
    Year        []uint16  `json:"year"`
}
package main

type Payload struct {
    Pathname    string    `json:"pathname"`
    Year        []uint16  `json:"year"`
}
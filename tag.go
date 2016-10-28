package main

type Tag struct {
    Name        string  `json:"name"`
    Probability float64 `json:"probability"`
    // + raw data (jaki typ?)
}

type Tags []*Tag

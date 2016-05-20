package main

import (
    "fmt"
)

type Collection struct {
    Name     string `json:"name"`
    Pathname string `json:"pathname"`
    Tracks   Tracks `json:"tracks,omitempty"`
}

func NewCollection(year uint16, key string) *Collection {
    c := &Collection{};
    c.Name = fmt.Sprintf("collection.%d.%s.musly", year, fmt.Sprintf("%03s", key));
    c.Pathname = fmt.Sprintf("/data/collections/%d/%s", year, c.Name)

    return c;
}

type Collections []*Collection
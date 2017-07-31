package main

import (
    "fmt"
)

type Collection struct {
    Name     string `json:"name"`
    Pathname string `json:"pathname"`
    Tracks   Tracks `json:"tracks,omitempty"`
}

func NewCollection(year uint16) *Collection {
    c := &Collection{};
    c.Name = fmt.Sprintf("collection.%d.musly", year);
    c.Pathname = fmt.Sprintf("/data/collections/%s", c.Name)

    return c;
}

type Collections []*Collection
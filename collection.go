package main

import (
    "fmt"
)

type Collection struct {
    Name            string `json:"name"`
    Pathname        string `json:"pathname"`
    JukeboxPathname string `json:"-"`
    Tracks          Tracks `json:"tracks,omitempty"`
}

func NewCollection() *Collection {
    c := &Collection{};
    c.Name = "collection.musly"
    c.Pathname = fmt.Sprintf("/data/collections/%s", c.Name)
    c.JukeboxPathname = "/data/collections/collection.jbox"

    return c;
}

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
    c.Pathname = fmt.Sprintf("/musly/%s", c.Name)
    c.JukeboxPathname = "/musly/collection.jbox"

    return c;
}

package main

import (
    "fmt"
)

func NewTrack(relativePathname string) *Track {
    t := &Track{}
    t.RelativePathname = relativePathname;
    t.Pathname = fmt.Sprintf("/collection%s", relativePathname);

    return t
}

type Track struct {
    InitialKey       string  `json:"initial_key"`
    Bpm              float64 `json:"bpm"`
    RelativePathname string  `json:"pathname"`
    Pathname         string  `json:"-"`
}

type Tracks []*Track

type SimilarTrack struct {
    *Track
    Similarity float64 `json:"similarity"`
}

type SimilarTracks []SimilarTrack

func (s SimilarTracks) Len() int           { return len(s) }
func (s SimilarTracks) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SimilarTracks) Less(i, j int) bool { return s[i].Similarity > s[j].Similarity }
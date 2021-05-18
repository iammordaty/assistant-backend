package main

func NewTrack(pathname string) *Track {
    t := &Track{}
    t.Pathname = pathname;

    return t
}

type Track struct {
    InitialKey  string  `json:"initial_key"`
    Bpm         float64 `json:"bpm"`
    Pathname    string  `json:"pathname"`
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

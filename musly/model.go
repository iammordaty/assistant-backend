package musly

type Track struct {
    RelativePathname string `json:"pathname"`
    Pathname         string `json:"-"`
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

type Collection struct {
    Name     string `json:"name"`
    Pathname string `json:"pathname"`
    Tracks   Tracks `json:"tracks,omitempty"`
}

type Collections []*Collection

type Payload struct {
    Pathname    string    `json:"pathname"`
    InitialKey  []string  `json:"initial_key"`
    Year        []uint16  `json:"year"`
}
package main

type SimilarTrack struct {
    *Track
    Similarity float64 `json:"similarity"`
}

type SimilarTracks []*SimilarTrack

func NewSimilarTrack(relativePathname string, similarity float64) *SimilarTrack {
    st := &SimilarTrack{}
    st.Track = NewTrack(relativePathname)
    st.Similarity = similarity

    return st
}

func (s SimilarTracks) Len() int           { return len(s) }
func (s SimilarTracks) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SimilarTracks) Less(i, j int) bool { return s[i].Similarity > s[j].Similarity }

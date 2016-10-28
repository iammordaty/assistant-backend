package main

import (
    "log"
    "net/http"

    "github.com/go-zoo/bone"
)

const (
    collectionRoot     = "/collection"
    muslyCollectionDir = "/data/collections"
)

func main() {
    mux := bone.New()

    // TODO: TrackController, CollectionController
    tc := NewTrackController()
    // sc := NewSimilarTracksController()

    mux.GetFunc("/collection/track/:pathname/bpm", tc.CalculateBpm)
    mux.GetFunc("/collection/track/:pathname/key", tc.CalculateKey)
    mux.GetFunc("/collection/track/:pathname/tags", tc.CalculateTags)
    // mux.GetFunc("/collection/track/:pathname/similar", sc.GetSimilarTracks)

    // mux.PostFunc("/collection/track", sc.AddTrackToCollection)

    log.Fatal(http.ListenAndServe(":80", mux))
}

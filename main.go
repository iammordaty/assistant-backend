package main

import (
    "log"
    "net/http"

    "github.com/go-zoo/bone"
)

const (
    collectionRoot      = "/collection"
    muslyCollectionsDir = "/data/collections"
)

func main() {
    mux := bone.New()

    // TODO: TrackController, CollectionController
    tc := NewTrackController()
    mux.GetFunc("/collection/track/:pathname/bpm", tc.CalculateBpm)
    mux.GetFunc("/collection/track/:pathname/key", tc.CalculateKey)
    mux.GetFunc("/collection/track/:pathname/tags", tc.CalculateTags)

    sc := NewCollectionController()
    mux.GetFunc("/collection/track/:pathname/similar", sc.GetSimilarTracks)
    mux.PostFunc("/collection/track", sc.AddTrackToCollection)

    log.Fatal(http.ListenAndServe(":80", mux))
}

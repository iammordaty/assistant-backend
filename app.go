package main

import (
    "log"
    "net/http"
    "runtime"

    "github.com/julienschmidt/httprouter"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    r := httprouter.New()

    ac := NewAudioDataController()
    r.GET("/track/*pathname", ac.CalculateTrackAudioData)

    sc := NewSimilarTracksController()

    r.POST("/musly/collection/tracks", sc.AddTrackToCollection)
    r.GET("/musly/similar/*pathname", sc.GetSimilarTracks)

    log.Fatal(http.ListenAndServe(":80", r))
}

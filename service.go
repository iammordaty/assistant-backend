package main

import (
    "log"
    "net/http"
    "runtime"

    "github.com/iammordaty/assistant-backend/musly"
    "github.com/iammordaty/assistant-backend/track"
    "github.com/julienschmidt/httprouter"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    r := httprouter.New()

    tc := track.NewController()
    r.GET("/track/*pathname", tc.CalculateAudioData)

    mc := musly.NewController()
    r.POST("/musly/collection/tracks", mc.AddTrackToCollection)
    r.GET("/musly/similar/*pathname", mc.GetSimilarTracks)

    log.Fatal(http.ListenAndServe(":80", r))
}

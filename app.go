package main

import (
    "fmt"
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

    r.GET("/ping", func (w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
        fmt.Fprint(w, "pong\n")
    })

    log.Fatal(http.ListenAndServe(":80", r))
}

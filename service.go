package main

import (
    "log"
    "net/http"
    "runtime"

    "github.com/iammordaty/assistant-backend/track"
    "github.com/julienschmidt/httprouter"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    r := httprouter.New()
    tc := track.NewController()

    r.GET("/track/*pathname", tc.GetInfo)

    log.Fatal(http.ListenAndServe(":80", r))
}

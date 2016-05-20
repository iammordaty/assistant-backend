package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strconv"

    "github.com/julienschmidt/httprouter"
)

type SimilarTracksController struct {}

func NewSimilarTracksController() *SimilarTracksController {
    return &SimilarTracksController{}
}

// POST /musly/collection/tracks
// { "initial_key": [ "4A", "5A" ], "year": [ 2015, 2014 ], "pathname": "pathname" }
// musly -x mp3 -a "track_pathname" -c "collection_pathname"
func (c SimilarTracksController) AddTrackToCollection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    defer r.Body.Close()

    var payload Payload;

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        RenderJson(w, &ErrorResponse{"Request body is not valid JSON string"}, http.StatusBadRequest)
        fmt.Println(err)

        return
    }

    if payload.Pathname == "" {
        RenderJson(w, &ErrorResponse{"Field \"pathname\" is required."}, http.StatusBadRequest)
        return
    }

    if len(payload.InitialKey) == 0 {
        RenderJson(w, &ErrorResponse{"Field \"initial_key\" is required."}, http.StatusBadRequest)
        return
    }

    if len(payload.Year) == 0 {
        RenderJson(w, &ErrorResponse{"Field \"year\" is required."}, http.StatusBadRequest)
        return
    }

    track := NewTrack(payload.Pathname)

    if filepath.Ext(track.Pathname) != ".mp3" {
        RenderJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file"}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        RenderJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
        return
    }

    collections := Collections{}

    for _, year := range payload.Year {
        for _, key := range payload.InitialKey {
            collections = append(collections, NewCollection(year, key))
        }
    }

    if err := EnsureCollections(collections); err != nil {
        RenderJson(w, &ErrorResponse{fmt.Sprint(err)}, http.StatusInternalServerError)
        return
    }

    if err := AddTrackToCollections(track, collections); err != nil {
        RenderJson(w, &ErrorResponse{fmt.Sprint(err)}, http.StatusInternalServerError)
        return
    }

    RenderJson(w, collections, http.StatusOK)
}

// GET /musly/track/track_pathname/similar?initial_key=4A&initial_key=5A&year=2015&year=2016
// musly -p "track_pathname" -k 100 -c "collection_pathname"
func (c SimilarTracksController) GetSimilarTracks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    defer r.Body.Close()

    payload := Payload{};
    payload.Pathname = p.ByName("pathname")
    payload.InitialKey = r.URL.Query()["initial_key"]

    for _, v := range r.URL.Query()["year"] {
        if s, err := strconv.ParseUint(v, 10, 16); err == nil {
            payload.Year = append(payload.Year, uint16(s))
        }
    }

    track := NewTrack(p.ByName("pathname"))

    if filepath.Ext(track.Pathname) != ".mp3" {
        RenderJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file"}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        RenderJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
        return
    }

    collections := Collections{}

    for _, year := range payload.Year {
        for _, key := range payload.InitialKey {
            collections = append(collections, NewCollection(year, key))
        }
    }

    similarTracks, err := GetSimilarTracks(track, collections)

    if err != nil {
        RenderJson(w, &ErrorResponse{fmt.Sprint(err)}, http.StatusInternalServerError)
        return
    }

    RenderJson(w, similarTracks, http.StatusOK)
}

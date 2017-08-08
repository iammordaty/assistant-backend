package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "path/filepath"

    "github.com/julienschmidt/httprouter"
)

type SimilarTracksController struct {}

func NewSimilarTracksController() *SimilarTracksController {
    return &SimilarTracksController{}
}

// POST /musly/collection/tracks
// { "pathname": "track_pathname" }
// musly -x mp3 -a "track_pathname" -c "collection_pathname"
func (c SimilarTracksController) AddTrackToCollection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    defer r.Body.Close()

    var payload map[string]string

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        RenderJson(w, &ErrorResponse{"Request body is not valid JSON string"}, http.StatusBadRequest)
        fmt.Println(err)

        return
    }

    track := NewTrack(payload["pathname"])

    if filepath.Ext(track.Pathname) != ".mp3" {
        RenderJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file"}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        RenderJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
        return
    }

    collection := NewCollection()

    if err := EnsureCollection(collection); err != nil {
        RenderJson(w, &ErrorResponse{fmt.Sprint(err)}, http.StatusInternalServerError)
        return
    }

    if err := AddTrackToCollection(track, collection); err != nil {
        RenderJson(w, &ErrorResponse{fmt.Sprint(err)}, http.StatusInternalServerError)
        return
    }

    RenderJson(w, collection, http.StatusOK)
}

// GET /musly/track/track_pathname/similar
// musly -p "track_pathname" -k 100 -c "collection_pathname"
func (c SimilarTracksController) GetSimilarTracks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    defer r.Body.Close()

    track := NewTrack(p.ByName("pathname"))

    if filepath.Ext(track.Pathname) != ".mp3" {
        RenderJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file"}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        RenderJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
        return
    }

    similarTracks, err := GetSimilarTracks(track, NewCollection())

    if err != nil {
        RenderJson(w, &ErrorResponse{fmt.Sprint(err)}, http.StatusInternalServerError)
        return
    }

    RenderJson(w, similarTracks, http.StatusOK)
}

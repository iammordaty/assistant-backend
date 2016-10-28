package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strconv"

    "github.com/go-zoo/bone"
)

type Collection struct {
    Name     string `json:"name"`
    Pathname string `json:"pathname"`
    Tracks   Tracks `json:"tracks,omitempty"`
}

type Collections []*Collection

func NewCollection(year uint16, key string) *Collection {
    c := &Collection{};
    c.Name = fmt.Sprintf("collection.%d.%s.musly", year, fmt.Sprintf("%03s", key));
    c.Pathname = fmt.Sprintf("%s/%d/%s", muslyCollectionsDir, year, c.Name)

    return c;
}

type CollectionController struct {}

func NewCollectionController() *CollectionController {
    return &CollectionController{}
}

func (c CollectionController) AddTrackToCollection(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    var payload Payload;

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        WriteJson(w, &ErrorResponse{"Request body is not valid JSON string"}, http.StatusBadRequest)
        fmt.Println(err)

        return
    }

    if payload.Pathname == "" {
        WriteJson(w, &ErrorResponse{"Field \"pathname\" is required."}, http.StatusBadRequest)
        return
    }

    if len(payload.InitialKey) == 0 {
        WriteJson(w, &ErrorResponse{"Field \"initial_key\" is required."}, http.StatusBadRequest)
        return
    }

    if len(payload.Year) == 0 {
        WriteJson(w, &ErrorResponse{"Field \"year\" is required."}, http.StatusBadRequest)
        return
    }

    track := NewTrack(payload.Pathname)

    if filepath.Ext(track.Pathname) != ".mp3" {
        WriteJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file"}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        WriteJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
        return
    }

    collections := Collections{}

    for _, year := range payload.Year {
        for _, key := range payload.InitialKey {
            collections = append(collections, NewCollection(year, key))
        }
    }

    if err := EnsureCollections(collections); err != nil {
        WriteJson(w, &ErrorResponse{fmt.Sprint(err)}, http.StatusInternalServerError)
        return
    }

    if err := AddTrackToCollections(track, collections); err != nil {
        WriteJson(w, &ErrorResponse{fmt.Sprint(err)}, http.StatusInternalServerError)
        return
    }

    WriteJson(w, collections, http.StatusOK)
}

func (c CollectionController) GetSimilarTracks(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()

    pathname, _ := url.QueryUnescape(bone.GetValue(r, "pathname")) // TODO: do NewTrack

    payload := Payload{};
    payload.Pathname = pathname
    payload.InitialKey = r.URL.Query()["initial_key"]

    for _, v := range r.URL.Query()["year"] {
        if s, err := strconv.ParseUint(v, 10, 16); err == nil {
            payload.Year = append(payload.Year, uint16(s))
        }
    }

    track := NewTrack(pathname)

    if filepath.Ext(track.Pathname) != ".mp3" {
        WriteJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file"}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        // WriteJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
        // return
    }

    collections := Collections{}

    for _, year := range payload.Year {
        for _, key := range payload.InitialKey {
            collections = append(collections, NewCollection(year, key))
        }
    }

    similarTracks, err := GetSimilarTracks(track, collections)

    if err != nil {
        WriteJson(w, &ErrorResponse{fmt.Sprint(err)}, http.StatusInternalServerError)
        return
    }

    WriteJson(w, similarTracks, http.StatusOK)
}

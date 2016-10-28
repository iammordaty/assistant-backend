package main

import (
    "fmt"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "strconv"

    "github.com/go-zoo/bone"
)

type Track struct {
    InitialKey       string  `json:"initial_key,omitempty"`
    Bpm              float64 `json:"bpm,omitempty"`
    Tags             Tags    `json:"tags,omitempty"`
    RelativePathname string  `json:"pathname"`
    Pathname         string  `json:"-"`
}

type Tracks []*Track

func NewTrack(relativePathname string) *Track {
    t := &Track{}
    t.RelativePathname = relativePathname
    t.Pathname = fmt.Sprintf("%s%s", collectionRoot, relativePathname)

    return t
}

type TrackController struct {}

func NewTrackController() *TrackController {
    return &TrackController{}
}

func (c TrackController) CalculateBpm(w http.ResponseWriter, r *http.Request) {
    pathname, _ := url.QueryUnescape(bone.GetValue(r, "pathname")) // TODO: do NewTrack
    track := NewTrack(pathname)

    if filepath.Ext(track.Pathname) != ".mp3" {
        WriteJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file"}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        WriteJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
        return
    }

    ch := make(chan CommandResult)
    RunCommand(fmt.Sprintf("sox \"%s\" -t raw -r 44100 -e floating-point -c 2 --norm -G - | bpm -f \"%%.1f\"", track.Pathname), ch)
    cr := <- ch

    track.Bpm, _ = strconv.ParseFloat(cr.Stdout, 64)

    if track.Bpm <= 100 {
        ch = make(chan CommandResult)
        RunCommand(fmt.Sprintf("sox \"%s\" -t raw -r 44100 -e floating-point -c 1 --norm -G - | bpm -f \"%%.1f\"", track.Pathname), ch)
        cr = <- ch

        track.Bpm, _ = strconv.ParseFloat(cr.Stdout, 64)
    }

    if cr.Stderr != "" {
        WriteJson(w, &CommandErrorResponse{cr.Stderr, "bpm"}, http.StatusInternalServerError)
        return
    }

    WriteJson(w, track, http.StatusOK)
}

func (c TrackController) CalculateKey(w http.ResponseWriter, r *http.Request) {
    pathname, _ := url.QueryUnescape(bone.GetValue(r, "pathname")) // TODO: do NewTrack
    track := NewTrack(pathname)

    if filepath.Ext(track.Pathname) != ".mp3" {
        WriteJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file"}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        WriteJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
        return
    }

    ch := make(chan CommandResult)
    RunCommand(fmt.Sprintf("keyfinder-cli -n camelot \"%s\"", track.Pathname), ch)

    cr := <- ch

    if cr.Error != nil {
        WriteJson(w, &CommandErrorResponse{cr.Stderr, "keyfinder-cli"}, http.StatusInternalServerError)
        return
    }

    track.InitialKey = cr.Stdout

    WriteJson(w, track, http.StatusOK)
}

// essentia_streaming_extractor_music Finally.mp3 Finally.mp3.json profile.in
func (c TrackController) CalculateTags(w http.ResponseWriter, r *http.Request) {
    pathname, _ := url.QueryUnescape(bone.GetValue(r, "pathname")) // TODO: do NewTrack
    track := NewTrack(pathname)

    if filepath.Ext(track.Pathname) != ".mp3" {
        WriteJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file"}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        WriteJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
        return
    }
}


package main

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strconv"

    "github.com/julienschmidt/httprouter"
)

type AudioDataController struct {}

func NewAudioDataController() *AudioDataController {
    return &AudioDataController{}
}

func (c AudioDataController) CalculateTrackAudioData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    track := NewTrack(p.ByName("pathname"))

    if filepath.Ext(track.Pathname) != ".mp3" {
        RenderJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file"}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        RenderJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound)}, http.StatusNotFound)
        return
    }

    kfch := make(chan CommandResult)
    RunCommand(fmt.Sprintf("keyfinder-cli -n camelot \"%s\"", track.Pathname), kfch)

    bpmch := make(chan CommandResult)
    RunCommand(fmt.Sprintf("sox -V1 \"%s\" -t raw -r 44100 -e floating-point -c 2 -G - | bpm -x 155 -f \"%%.1f\"", track.Pathname), bpmch)

    kfr := <- kfch
    bpmr := <- bpmch

    if kfr.Error != nil {
        RenderJson(w, &CommandErrorResponse{kfr.Stderr, "keyfinder-cli"}, http.StatusInternalServerError)
        return
    }

    track.InitialKey = kfr.Stdout
    track.Bpm, _ = strconv.ParseFloat(bpmr.Stdout, 64)

    if track.Bpm <= 100 {
        bpmch = make(chan CommandResult)
        RunCommand(fmt.Sprintf("sox -V1 \"%s\" -t raw -r 44100 -e floating-point -c 1 -G - | bpm -x 155 -f \"%%.1f\"", track.Pathname), bpmch)
        bpmr = <- bpmch

        track.Bpm, _ = strconv.ParseFloat(bpmr.Stdout, 64)
    }

    if bpmr.Error != nil || bpmr.Stderr != "" {
        RenderJson(w, &CommandErrorResponse{bpmr.Stderr, "bpm"}, http.StatusInternalServerError)
        return
    }

    RenderJson(w, track, http.StatusOK)
}

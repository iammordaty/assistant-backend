package track

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strconv"

    "github.com/iammordaty/assistant-backend/helper"
    "github.com/julienschmidt/httprouter"
)

type Controller struct {}

func NewController() *Controller {
    return &Controller{}
}

func (c Controller) GetInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    pathname := fmt.Sprintf("/collection%s", p.ByName("pathname"))

    if filepath.Ext(pathname) != ".mp3" {
        helper.RenderJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file", pathname}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(pathname); os.IsNotExist(err) {
        helper.RenderJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound), pathname}, http.StatusNotFound)
        return
    }

    kfch := make(chan helper.CommandResult)
    helper.RunCommand(fmt.Sprintf("keyfinder-cli -n camelot \"%s\"", pathname), kfch)

    bpmch := make(chan helper.CommandResult)
    helper.RunCommand(fmt.Sprintf("sox \"%s\" -t raw -r 44100 -e float -c 1 -G - | bpm -f \"%%.1f\"", pathname), bpmch)

    kfr := <- kfch
    bpmr := <- bpmch

    if kfr.Error != nil {
        helper.RenderJson(w, &CommandErrorResponse{ErrorResponse{kfr.Stderr, pathname}, "keyfinder-cli"}, http.StatusInternalServerError)
        return
    }

    if bpmr.Error != nil || bpmr.Stderr != "" {
        helper.RenderJson(w, &CommandErrorResponse{ErrorResponse{bpmr.Stderr, pathname}, "bpm"}, http.StatusInternalServerError)
        return
    }

    sr := &SuccessResponse{}
    sr.InitialKey = kfr.Stdout
    sr.Bpm, _ = strconv.ParseFloat(bpmr.Stdout, 64)

    helper.RenderJson(w, sr, http.StatusOK)
}

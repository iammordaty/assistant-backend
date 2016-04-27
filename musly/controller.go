package musly

import (
    "encoding/json"
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

func NewTrack(relativePathname string) *Track {
    t := &Track{}
    t.RelativePathname = relativePathname;
    t.Pathname = fmt.Sprintf("/collection%s", relativePathname);

    return t
}

func NewCollection(year uint16, key string) *Collection {
    c := &Collection{};
    c.Name = fmt.Sprintf("collection.%d.%s.musly", year, fmt.Sprintf("%03s", key));
    c.Pathname = fmt.Sprintf("/data/collections/%d/%s", year, c.Name)

    return c;
}

// POST /musly/collection/tracks
// { "initial_key": [ "4A", "5A" ], "year": [ 2015, 2014 ], "pathname": "pathname" }
// musly -x mp3 -a "track_pathname" -c "collection_pathname"
func (c Controller) AddTrackToCollection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    defer r.Body.Close()

    var payload Payload;

    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        helper.RenderJson(w, &ErrorResponse{"Request body is not valid JSON string", ""}, http.StatusBadRequest)
        return
    }

    if payload.Pathname == "" {
        helper.RenderJson(w, &ErrorResponse{"Field \"pathname\" is required.", ""}, http.StatusBadRequest)
        return
    }

    if len(payload.InitialKey) == 0 {
        helper.RenderJson(w, &ErrorResponse{"Field \"initial_key\" is required.", ""}, http.StatusBadRequest)
        return
    }

    if len(payload.Year) == 0 {
        helper.RenderJson(w, &ErrorResponse{"Field \"year\" is required.", ""}, http.StatusBadRequest)
        return
    }

    track := NewTrack(payload.Pathname)

    if filepath.Ext(track.Pathname) != ".mp3" {
        helper.RenderJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file", track.Pathname}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        helper.RenderJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound), track.Pathname}, http.StatusNotFound)
        return
    }

    collections := Collections{}

    for _, year := range payload.Year {
        for _, key := range payload.InitialKey {
            collections = append(collections, NewCollection(year, key))
        }
    }

    if err := ensureCollections(collections); err != nil {
        helper.RenderJson(w, ErrorResponse{fmt.Sprint(err), ""}, http.StatusInternalServerError)
        return
    }

    if err := addTrackToCollections(track, collections); err != nil {
        helper.RenderJson(w, ErrorResponse{fmt.Sprint(err), ""}, http.StatusInternalServerError)
        return
    }

    helper.RenderJson(w, collections, http.StatusOK)
}

// GET /musly/track/track_pathname/similar?initial_key=4A&initial_key=5A&year=2015&year=2016
// musly -p "track_pathname" -k 100 -c "collection_pathname"
func (c Controller) GetSimilarTracks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
        helper.RenderJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file", track.Pathname}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        helper.RenderJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound), track.Pathname}, http.StatusNotFound)
        return
    }

    collections := Collections{}

    for _, year := range payload.Year {
        for _, key := range payload.InitialKey {
            collections = append(collections, NewCollection(year, key))
        }
    }

    similarTracks, err := getSimilarTracks(track, collections)

    if err != nil {
        helper.RenderJson(w, ErrorResponse{fmt.Sprint(err), ""}, http.StatusInternalServerError)
        return
    }

    helper.RenderJson(w, similarTracks, http.StatusOK)
}

package musly

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"

    "github.com/iammordaty/assistant-backend/helper"
    "github.com/julienschmidt/httprouter"
    "github.com/tv42/slug"
)

type Controller struct {}

func NewController() *Controller {
    return &Controller{}
}

// POST /musly/collection/tracks
// { "genre": "Techno", "year": 2015, "pathname": "track_pathname" }
// musly -x mp3 -a "track_pathname" -c "collection_pathname"
func (c Controller) AddTrackToCollection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    var b map[string]interface{}

    defer r.Body.Close()

    if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
        helper.RenderJson(w, &ErrorResponse{"Request body is not valid JSON string", ""}, http.StatusBadRequest)
        return
    }

    genre, ok := b["genre"].(string)

    if ok == false {
        helper.RenderJson(w, &ErrorResponse{"Field \"genre\" is required.", ""}, http.StatusBadRequest)
        return
    }

    year, ok := b["year"].(float64)

    if ok == false {
        helper.RenderJson(w, &ErrorResponse{"Field \"year\" is required.", ""}, http.StatusBadRequest)
        return
    }

    tPathname, ok := b["pathname"].(string)

    if ok == false {
        helper.RenderJson(w, &ErrorResponse{"Field \"pathname\" is required.", ""}, http.StatusBadRequest)
        return
    }

    track := &Track{}
    track.Pathname = fmt.Sprintf("/collection%s", tPathname)

    if filepath.Ext(track.Pathname) != ".mp3" {
        helper.RenderJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file", track.Pathname}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(track.Pathname); os.IsNotExist(err) {
        helper.RenderJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound), track.Pathname}, http.StatusNotFound)
        return
    }

    collection := &Collection{}
    collection.Name = fmt.Sprintf("collection.%d.%s.musly", int(year), slug.Slug(genre));
    collection.Pathname = fmt.Sprintf("/data/%s", collection.Name)

    if mr := c.ensureCollection(collection); mr.Error != nil {
        helper.RenderJson(w, ErrorResponse{mr.Stderr, collection.Pathname}, http.StatusInternalServerError)
        return
    }

    ch := make(chan helper.CommandResult)
    helper.RunCommand(fmt.Sprintf("musly -a \"%s\" -c \"%s\"", track.Pathname, collection.Pathname), ch)

    mr := <- ch

    if mr.Error != nil {
        helper.RenderJson(w, ErrorResponse{mr.Stderr, track.Pathname}, http.StatusInternalServerError)
        return
    }

    if strings.Contains(mr.Stdout, "[FAILED]") {
        helper.RenderJson(w, ErrorResponse{"Track can not be added to collection", track.Pathname}, http.StatusInternalServerError)
        return
    }

    collection.Tracks = Tracks{track}

    helper.RenderJson(w, collection, http.StatusOK)
}

// GET /musly/collection/%collection_pathname%/tracks
// musly -l -c /data/collection.2015.tech-house.musly
func (c Controller) GetCollectionTracks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    var line string

    pathname := fmt.Sprintf("/data/%s", p.ByName("collection_pathname"))

    if _, err := os.Stat(pathname); os.IsNotExist(err) {
        helper.RenderJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound), pathname}, http.StatusNotFound)
        return
    }

    mch := make(chan helper.CommandResult)
    helper.RunCommand(fmt.Sprintf("musly -l -c \"%s\"", pathname), mch)

    mr := <- mch

    if mr.Error != nil {
        helper.RenderJson(w, ErrorResponse{mr.Stderr, pathname}, http.StatusInternalServerError)
        return
    }

    collection := &Collection{}
    collection.Name = p.ByName("collection_pathname")
    collection.Pathname = pathname

    scanner := bufio.NewScanner(strings.NewReader(mr.Stdout))

    for scanner.Scan() {
        line = scanner.Text()

        if strings.HasPrefix(line, "track-id") {
            collection.Tracks = append(collection.Tracks, &Track{strings.SplitAfter(line, "track-origin: ")[1]} )
        }
    }

    helper.RenderJson(w, collection, http.StatusOK)
}

// GET /musly/track/track_pathname/similar?genre=Techno&genre=Tech-House&year=2015&year=2016
// musly -p "track_pathname" -k 100 -c "collection_pathname"
func (c Controller) GetSimilarTracks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    pathname := fmt.Sprintf("/collection%s", p.ByName("pathname"))

    if filepath.Ext(pathname) != ".mp3" {
        helper.RenderJson(w, &ErrorResponse{"Pathname does not seems to be an mp3 file", pathname}, http.StatusBadRequest)
        return
    }

    if _, err := os.Stat(pathname); os.IsNotExist(err) {
        helper.RenderJson(w, &ErrorResponse{http.StatusText(http.StatusNotFound), pathname}, http.StatusNotFound)
        return
    }

    var collections []string
    var collection string

    for _, year := range r.URL.Query()["year"] {
        for _, genre := range r.URL.Query()["genre"] {
            collection = fmt.Sprintf("/data/collection.%s.%s.musly", year, slug.Slug(genre));

            if _, err := os.Stat(collection); err == nil {
                collections = append(collections, collection)
            }
        }
    }

	var chans = []chan helper.CommandResult{}

	for i := 0; i < len(collections); i++ {
		mch := make(chan helper.CommandResult)
		helper.RunCommand(fmt.Sprintf("musly -p \"%s\" -k10 -c \"%s\"", pathname, collections[i]), mch)

		chans = append(chans, mch)
	}

    var buffer bytes.Buffer
    var line string

	for i := 0; i < len(chans); i++ {
		mr := <- chans[i]

        buffer.WriteString(mr.Stdout + "\n")
    }

    st := SimilarTracks{}
    scanner := bufio.NewScanner(strings.NewReader(buffer.String()))

    for scanner.Scan() {
        line = scanner.Text()

        if strings.HasPrefix(line, "track-id") {
            pathname := strings.SplitAfter(line, "track-origin: ")[1]
            similarity, _ := strconv.ParseFloat(strings.Split(strings.SplitAfter(line, "track-similarity: ")[1], ", ")[0], 64)
            similarity, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", 100 - (similarity * 100)), 64)

            st = append(st, SimilarTrack{Track{pathname}, similarity})
        }
    }

    sort.Sort(st)

    helper.RenderJson(w, st, http.StatusOK)
}

func (c Controller) ensureCollection(collection *Collection) helper.CommandResult {
    if _, err := os.Stat(collection.Pathname); err == nil {
        return helper.CommandResult{}
    }

    ch := make(chan helper.CommandResult)
    helper.RunCommand(fmt.Sprintf("musly -n timbre -c \"%s\"", collection.Pathname), ch)

    r := <- ch

    return r
}
package main

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
)

// Ensures that collections are exists has been initialized
func EnsureCollection(collection *Collection) (err error) {
    if _, err = os.Stat(collection.Pathname); err == nil {
        return
    }

    dir := filepath.Dir(collection.Pathname)

    if _, err := os.Stat(dir); os.IsNotExist(err) {
        os.MkdirAll(dir, 0777)
    }

    ch := make(chan CommandResult)
    RunCommand(fmt.Sprintf("musly -n timbre -c \"%s\"", collection.Pathname), ch)

    cr := <- ch

    if cr.Error != nil {
        err = errors.New(fmt.Sprintf("An error occurred when initializing collection: %s.", cr.Stderr))
    }

    return
}

// Adds track to collection
func AddTrackToCollection(track *Track, collection *Collection) (err error) {
    ch := make(chan CommandResult)
    RunCommand(fmt.Sprintf("musly -a \"%s\" -c \"%s\"", track.Pathname, collection.Pathname), ch)

    cr := <- ch

    if cr.Error != nil {
        err = errors.New(fmt.Sprintf("An error occurred when adding track to collection: %s.", cr.Stderr))

        return
    }

    if strings.Contains(cr.Stdout, "[FAILED]") {
        err = errors.New("An error occurred when adding track to collection: failed.")

        return
    }

    return
}

// Returns similar tracks
func GetSimilarTracks(track *Track, collection *Collection) (similarTracks SimilarTracks, err error) {
    ch := make(chan CommandResult)
    RunCommand(fmt.Sprintf("musly -p \"%s\" -k 200 -c \"%s\" -j \"%s\"", track.Pathname, collection.Pathname, collection.JukeboxPathname), ch)

    cr := <- ch

    similarTracks = SimilarTracks{}
    scanner := bufio.NewScanner(strings.NewReader(cr.Stdout))

    for scanner.Scan() {
        line := scanner.Text()

        if strings.HasPrefix(line, "track-id") == false {
            continue
        }

        pathname := strings.Replace(strings.SplitAfter(line, "track-origin: ")[1], "/collection", "", 1)

        distance, _ := strconv.ParseFloat(strings.Split(strings.SplitAfter(line, "track-similarity: ")[1], ", ")[0], 64)
        similarity, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", 100 - (distance * 100)), 64)

        similarTracks = append(similarTracks, SimilarTrack{NewTrack(pathname), similarity})
    }

    sort.Sort(similarTracks)

    return
}

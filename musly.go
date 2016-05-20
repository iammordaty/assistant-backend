package main

import (
    "bufio"
    "bytes"
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"
)

// Ensures that collections are exists has been initialized
func EnsureCollections(collections Collections) (err error) {
    var chans = []chan CommandResult{}

    for i := 0; i < len(collections); i++ {
        if _, err := os.Stat(collections[i].Pathname); err == nil {
            continue
        }

        dir := filepath.Dir(collections[i].Pathname)

        if _, err := os.Stat(dir); os.IsNotExist(err) {
            os.MkdirAll(dir, 0777)
        }

        ch := make(chan CommandResult)
        chans = append(chans, ch)

        RunCommand(fmt.Sprintf("musly -n timbre -c \"%s\"", collections[i].Pathname), ch)
    }

    if len(chans) == 0 {
        return
    }

    var crs = []CommandResult{}

    for i := 0; i < len(chans); i++ {
        cr := <- chans[i]

        crs = append(crs, cr)
    }

    for i := 0; i < len(crs); i++ {
        if crs[i].Error != nil {
            err = errors.New(fmt.Sprintf("An error occurred when initializing collection: %s.", crs[i].Stderr))
            break
        }
    }

    return
}

// Adds track to collections
func AddTrackToCollections(track *Track, collections Collections) (err error) {
    var chans = []chan CommandResult{}

    for i := 0; i < len(collections); i++ {
        ch := make(chan CommandResult)
        chans = append(chans, ch)

        RunCommand(fmt.Sprintf("musly -a \"%s\" -c \"%s\"", track.Pathname, collections[i].Pathname), ch)
    }

    var crs = []CommandResult{}

    for i := 0; i < len(chans); i++ {
        cr := <- chans[i]

        crs = append(crs, cr)
    }

    for i := 0; i < len(crs); i++ {
        if crs[i].Error != nil {
            err = errors.New(fmt.Sprintf("An error occurred when adding track to collection: %s.", crs[i].Stderr))
            break
        }

        if strings.Contains(crs[i].Stdout, "[FAILED]") {
            err = errors.New("An error occurred when adding track to collection: failed.")
            break
        }
    }

    return
}

// Returns similar tracks
func GetSimilarTracks(track *Track, collections Collections) (similarTracks SimilarTracks, err error) {
    var chans = []chan CommandResult{}

    for i := 0; i < len(collections); i++ {
        ch := make(chan CommandResult)
        RunCommand(fmt.Sprintf("musly -p \"%s\" -k20 -c \"%s\"", track.Pathname, collections[i].Pathname), ch)

        chans = append(chans, ch)
    }

    var crs = []CommandResult{}

    for i := 0; i < len(chans); i++ {
        mr := <- chans[i]

        crs = append(crs, mr)
    }

    var stdouts bytes.Buffer

    for i := 0; i < len(crs); i++ {
        stdouts.WriteString(crs[i].Stdout + "\n")
    }

    occuriences := map[string]float64{}
    similaritySum := map[string]float64{}

    scanner := bufio.NewScanner(strings.NewReader(stdouts.String()))

    for scanner.Scan() {
        line := scanner.Text()

        if strings.HasPrefix(line, "track-id") == false {
            continue
        }

        pathname := strings.Replace(strings.SplitAfter(line, "track-origin: ")[1], "/collection", "", 1)
        similarity, _ := strconv.ParseFloat(strings.Split(strings.SplitAfter(line, "track-similarity: ")[1], ", ")[0], 64)

        similaritySum[pathname] += similarity
        occuriences[pathname]++
    }

    similarTracks = SimilarTracks{}

    for pathname, similarity := range similaritySum {
        similarity, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", 100 - (similarity / occuriences[pathname] * 100)), 64)

        similarTracks = append(similarTracks, SimilarTrack{NewTrack(pathname), similarity})
    }

    sort.Sort(similarTracks)

    return
}

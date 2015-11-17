package helper

import (
	"encoding/json"
	"net/http"
)

func RenderJson(w http.ResponseWriter, obj interface{}, s int) {
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(s)

    json.NewEncoder(w).Encode(obj)
}
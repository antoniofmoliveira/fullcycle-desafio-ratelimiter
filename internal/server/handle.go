package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

var availableHandles = []handle{getHandleHello, postHandleHello}

type handle struct {
	path    string
	handler http.Handler
}

// GET /hello
var getHandleHello = handle{
	path: "GET /hello",
	handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}),
}

// POST /hello
var postHandleHello = handle{
	path: "POST /hello",
	handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Can't read body", "error", err)
			http.Error(w, "Can't read body", http.StatusInternalServerError)
			return
		}
		var message struct {
			Message string `json:"message"`
		}
		err = json.Unmarshal(body, &message)
		if err != nil {
			slog.Error("Can't unmarshal json", "error", err)
			http.Error(w, "Can't unmarshal json", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Hello, %s!", message.Message)
	}),
}

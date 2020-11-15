package handlers

import (
	"net/http"
	"log"
)

// Goodbye does
type Goodbye struct {
	l *log.Logger
}

// NewGoodbye takes a logger and returns a Goodbye
// struct containing it
func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (g *Goodbye) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("goodbye world"))

}
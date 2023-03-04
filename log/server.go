package log

import (
	"io"
	stlog "log"
	"net/http"
	"os"
)

var log *stlog.Logger

type fileLog string

func (fl fileLog) Write(p []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(p)
}

func Run(description string) {
	log = stlog.New(fileLog(description), "Ahriknow: ", stlog.LstdFlags)
}

func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			msg, err := io.ReadAll(r.Body)
			if err != nil || len(msg) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(msg))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}

func write(msg string) {
	log.Printf("%v\n", msg)
}

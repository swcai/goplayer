package main

import (
	"math/rand"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

type Entry struct {
	Name  string
}

var (
	addr = flag.String("http", ":8080", "http listen address")
	root = flag.String("root", "/var/music", "music root")
	entries []Entry
)

func buildPlayList(path string) []Entry{
	defer func() {
		if _, ok := recover().(error); ok {
			log.Println("failed to build a playlist")
		}
	}()
	d, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	files, err := d.Readdir(-1)
	if err != nil {
		panic(err)
	}

	entries := make([]Entry, 0, len(files))
	dirs := make([]string, 0)
	for k := range files {
		if !files[k].IsDir() {
			entry := Entry{path + "/" + files[k].Name()}
			log.Print(entry.Name)
			entries = append(entries, entry)
		} else {
			dirs = append(dirs, files[k].Name())
		}
	}

	if len(dirs) != 0 {
		for _, path := range dirs {
			entries = append(entries, buildPlayList(*root + "/" + path)...)
		}
	}
	return entries
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	entries = buildPlayList(*root)
	http.HandleFunc("/", index)
	http.HandleFunc("/random", randomFile)
	http.ListenAndServe(*addr, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
	log.Print("index called")
}

func randomFile(w http.ResponseWriter, r *http.Request) {
	index := rand.Int() % len(entries)
	_, err := os.Stat(entries[index].Name)
	log.Print("File called: ", entries[index])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, entries[index].Name)
}

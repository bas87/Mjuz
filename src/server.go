// The MIT License (MIT)
//
// Copyright (c) 2013 ushi <ushi@honkgong.info>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CON

// This package implements simple webserver for control mocp

package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"os"
	"os/exec"
	"net/http"
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"./m3u"
)

const PORT int = 3000
const MUSIC_DIR string = "/media/KINGSTON/"
const PLAY_LIST string = "/home/pi/.moc/playlist.m3u"
const SHELL string = "/bin/sh" 

// todo: kouknout jestli to je treba 
const STATIC_URL string = "/web/"
const STATIC_ROOT string = "../web/"

// DTO
type State struct {
  State string
  File string
  Title string
  Artist string
  SongTitle string
  Album string
  TotalTime string
  TimeLeft string
  TotalSec string
  CurrentTime string
  CurrentSec string
  Bitrate string
  Rate string
}

type Track struct {
  Path  string
  Title string
  Time  int64
  FileExt string
}

type Fail struct {
  Code  int
  Msg string
}

type Tracklist []Track

// Helpers
func mocp(param... string) string {

	// Curr. dir
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Prepend
	param = append([]string{wd + "/mocp.sh"}, param...)

	// Exec cmd
	out, err := exec.Command(SHELL, param...).Output()
	if err != nil {
		log.Fatal(err)
	}

	// prepare aoutput
	resp := string(out)
	resp = strings.Trim(resp, "\n")
	resp = strings.Trim(resp, "\r")
	resp = strings.Trim(resp, " ")

	return resp
}

func parse_state(key string, data string) string {
	var result string;
	res := regexp.MustCompile(key + ": (.*)\n").FindStringSubmatch(data)
	
	if (len(res) > 1) {
		result = res[1]
	} else {
		result = ""
	}

	return result
}


func youtube_dl(url string) {

	// Curr. dir
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Exec cmd
	cmd_err := exec.Command(SHELL, wd + "/youtube-dl.sh", url, MUSIC_DIR, wd + "/../tmp/").Run()

	if cmd_err != nil {
		log.Fatal(err)
	}
}

// API's
func Play(w http.ResponseWriter, req *http.Request) {
	mocp("-p")
}

func Stop(w http.ResponseWriter, req *http.Request) {
	mocp("-s")
}

func Pause(w http.ResponseWriter, req *http.Request) {
	mocp("-P")
}

func TogglePause(w http.ResponseWriter, req *http.Request) {
	mocp("-G")
}

func UnPause(w http.ResponseWriter, req *http.Request) {
	mocp("-U")
}

func Next(w http.ResponseWriter, req *http.Request) {
	mocp("-f")
}

func Prev(w http.ResponseWriter, req *http.Request) {
	mocp("-r")
}

func Clear(w http.ResponseWriter, req *http.Request) {
	mocp("-c")
}

func Append(w http.ResponseWriter, req *http.Request) {
	data, err := base64.StdEncoding.DecodeString(string(mux.Vars(req)["path"]))
	if err != nil {
		log.Fatal(err)
		return
	}

	mocp("-a", string(data))
}

func Download(w http.ResponseWriter, req *http.Request) {
	data, err := base64.StdEncoding.DecodeString(string(mux.Vars(req)["url"]))
	if err != nil {
		log.Fatal(err)
		return
	}

	youtube_dl(string(data))
}

func Info(w http.ResponseWriter, req *http.Request) {
	data := mocp("-i")

	state := State {
		parse_state("State", data),
		parse_state("File", data),
		parse_state("Title", data),
		parse_state("Artist", data),
		parse_state("SongTitle", data),
		parse_state("Album", data),
		parse_state("TotalTime", data),
		parse_state("TimeLeft", data),
		parse_state("TotalSec", data),
		parse_state("CurrentTime", data),
		parse_state("Bitrate", data),
		parse_state("AvgBitrate", data),
		parse_state("Rate", data),
	}
	JsonRender(w, state)
}

func List(w http.ResponseWriter, req *http.Request) {
	tl := Tracklist{}
	files, _ := filepath.Glob(MUSIC_DIR + "*.mp3") 
    for _, f := range files {

        tl = append(tl, Track{
        	Path: f,
        	Title: "",
        	Time: -1,
        	FileExt: filepath.Ext(f)[1:],
        })

    }

    JsonRender(w, tl)
}

func PlayList(w http.ResponseWriter, req *http.Request) {
	f, err := os.Open(PLAY_LIST)

	if err != nil {
	    JsonRender(w, Fail { Code: 500, Msg: "PlayList not exists."})
	    return
	}
	defer f.Close()

	pl, err := m3u.Parse(f)

	if err != nil {
	    log.Fatal(err)
	}

	JsonRender(w, pl)
}

func JsonRender(w http.ResponseWriter, data interface{}) {
  js, err := json.Marshal(data)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func main() {
	r := mux.NewRouter()

	// Base Page
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, STATIC_ROOT + "index.html")
	})

	// Api
	r.HandleFunc("/api/play/", Play).Methods("GET")
	r.HandleFunc("/api/stop/", Stop).Methods("GET")
	r.HandleFunc("/api/prev/", Prev).Methods("GET")
	r.HandleFunc("/api/next/", Next).Methods("GET")
	r.HandleFunc("/api/pause/", Pause).Methods("GET")
	r.HandleFunc("/api/un-pause/", UnPause).Methods("GET")
	r.HandleFunc("/api/toogle-pause/", TogglePause).Methods("GET")
	r.HandleFunc("/api/info/", Info).Methods("GET")
	r.HandleFunc("/api/list/", List).Methods("GET")
	r.HandleFunc("/api/play-list/", PlayList).Methods("GET")
	r.HandleFunc("/api/clear/", Clear).Methods("GET")
	r.HandleFunc("/api/append/", Append).Queries("path", "{path}").Methods("GET")
	r.HandleFunc("/api/download/", Download).Queries("url", "{url}").Methods("GET")

	// Static content
	fileHandler := http.FileServer(http.Dir(STATIC_ROOT))
	r.PathPrefix(STATIC_URL).Handler(http.StripPrefix(STATIC_URL, fileHandler))


	log.Println("Serving on:", strconv.Itoa(PORT))
	http.ListenAndServe(":" + strconv.Itoa(PORT), context.ClearHandler(r))
}

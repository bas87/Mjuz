// The MIT License (MIT)
//
// Copyright (c) 2015 bas <toman@devzone.cz>
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
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// This package implements a default tolerant m3u parser
//
// The spec can be found at http://www.scvi.net/pls.htm
package m3u

import (
  "bufio"
  "fmt"
  "io"
  "strconv"
  "strings"
  "path/filepath"
)

// Represents a list of tracks.
type Playlist []Track

// Represents a single track.
type Track struct {
  Path  string // path to the file
  Title string // title of the track
  Time  int64  // duration of the track
  FileExt string // type of file, f.e: mp3, ogg ...
}

// Parses simple and extended m3u files. Returns the playlist.
func Parse(r io.Reader) (Playlist, error) {
  br := bufio.NewReader(r)
  pl := Playlist{}

  for {
    line, err := br.ReadString('\n')

    if err != nil {
      if err == io.EOF {
        return pl, nil
      }
      return pl, err
    }
    line = line[:len(line)-1]

    if len(line) > 0 && line[0] != '#' {
      pl = append(pl, Track{
        Path: line,
        Title: "",
        Time: -1,
        FileExt: filepath.Ext(line)[1:],
        })
      continue
    }

    if len(line) > 8 && line[:8] == "#EXTINF:" {
      i := strings.Index(line[8:], ",")

      if i < 0 {
        return pl, fmt.Errorf("unexpected line: %q", line)
      }
      time, err := strconv.ParseInt(line[8:i+8], 10, 64)

      if err != nil {
        return pl, err
      }
      path, err := br.ReadString('\n')

      if err != nil {
        return pl, err
      } else {
        path =  path[:len(path)-1]
      }

      pl = append(pl, Track{
        Path: path, 
        Title: line[i+9:], 
        Time: time, 
        FileExt: filepath.Ext(path)[1:],
        })
    }
  }
}
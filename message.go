package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	
	"path/filepath"
	// "sort"
	// "strconv"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var exPath string

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath = filepath.Dir(ex)
}

// handleMessages handles messages
func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "explore":
		// Unmarshal payload
		var path string
		if len(m.Payload) > 0 {
			// Unmarshal payload
			if err = json.Unmarshal(m.Payload, &path); err != nil {
				payload = err.Error()
				return
			}
		}

		// Explore
		if payload, err = explore(path); err != nil {
			payload = err.Error()
			return
		}
	}

	return
}

// Exploration represents the results of an exploration
type Exploration struct {
	Dirs       []Dir              `json:"dirs"`
	Path       string             `json:"path"`
	RenderHTML string             `json:"renderHTML"`
}

// PayloadDir represents a dir payload
type Dir struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// explore explores a path.
// If path is empty, it explores the user's home directory
func explore(path string) (e Exploration, err error) {

	if filepath.Ext(path) == ".md" {
		var b []byte
		b, err = ioutil.ReadFile(path) // just pass the file name
		if err != nil {
			return
		}

		output := blackfriday.Run(b)

		return Exploration{
			Dirs: []Dir{},
			Path: path,
			RenderHTML: string(output),
		}, nil
	}

	// If no path is provided, get current path
	if len(path) == 0 {
		path = exPath
	}

	// Read dir
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(path); err != nil {
		return
	}

	// Init exploration
	e = Exploration{
		Dirs: []Dir{},
		Path: path,
	}

	if path != exPath && filepath.Dir(path) != path {
		e.Dirs = append(e.Dirs, Dir{
			Name: "..",
			Path: filepath.Dir(path),
		})
	}

	// Loop through files
	var fileArr []Dir

	for _, f := range files {
		if f.IsDir() {
			e.Dirs = append(e.Dirs, Dir{
				Name: f.Name(),
				Path: filepath.Join(path, f.Name()),
			})
		} else if filepath.Ext(f.Name()) == ".md" {
			fileArr = append(fileArr, Dir{
				Name: f.Name(),
				Path: filepath.Join(path, f.Name()),
			})

		}
	}

	e.Dirs = append(e.Dirs, fileArr...)

	return
}

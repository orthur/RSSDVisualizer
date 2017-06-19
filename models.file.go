package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type file struct {
	ID     int         `json:"id"`
	Name   string      `json:"name"`
	Path   string      `json:"path"`
	Points [][]float64 `json:"points"`
}

const BASEAPTH = "SampleFiles"

var fileList = []file{}

// Fetch all files
func getAllFiles() []file {
	return fileList
}

// Delete a file based on the ID supplied
func deleteFileByID(id int) (*file, error) {
	for _, a := range fileList {
		if a.ID == id {
			if err := os.Remove(a.Path); err == nil {
				return &a, nil
			} else {
				return nil, err
			}
		}
	}
	return nil, errors.New("File not found")
}

// Fetch a file based on the ID supplied
func getFileByID(id int) (*file, error) {
	for i, a := range fileList {
		if a.ID == id {
			return &fileList[i], nil
		}
	}
	return nil, errors.New("File not found")
}

// Upload a new file
func uploadNewFile(newFile multipart.File, handler *multipart.FileHeader) (*file, error) {
	f, err := os.OpenFile(handler.Filename+".tmp", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	io.Copy(f, newFile)
	err = os.Rename(handler.Filename+".tmp", BASEAPTH+"/"+handler.Filename)
	if err != nil {
		return nil, err
	}
	return &file{}, nil
}

// Parse file
func parseFile(fileToParse *file) ([][]float64, error) {
	var pts [][]float64
	f, err := os.Open(fileToParse.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 0 && !strings.HasPrefix(line, "#") {
			splited := strings.Fields(line)
			if 2 == len(splited) {
				var time, val float64
				time, err = strconv.ParseFloat(splited[0], 64)
				if err != nil {
					continue
				}
				val, err = strconv.ParseFloat(splited[1], 64)
				if err != nil {
					continue
				}
				pts = append(pts, []float64{time, val})
			}
		}
	}
	if len(pts) == 0 {
		return nil, errors.New("Empty file or parsing error")
	}
	return pts, nil
}

// Watch directory
func watchDir(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR START WATCHER", err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				_, fileName := path.Split(event.Name)
				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("new file:", event.Name)
					a := file{ID: len(fileList) + 1, Name: fileName, Path: BASEAPTH + "/" + fileName}
					fileList = append(fileList, a)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Println("removed file:", event.Name)
					for i, l := range fileList {
						if l.Name == fileName {
							fileList = append(fileList[:i], fileList[i+1:]...)
							break
						}
					}
				} else if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("edited file:", event.Name)
				}
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()
	if err := watcher.Add(dir); err != nil {
		fmt.Println("ERROR", err)
	}
	<-done
}

package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

var watchingFolders string
var scavangingInterval string

func main() {
	ReadEnv()

	Run(watchingFolders)
}

func Run(path string) {

	isRunning := false

	duration, err := time.ParseDuration(scavangingInterval)
	if err != nil {
		log.Fatalf("parsing failed: %v", err)
	}
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			if !isRunning {
				isRunning = true
				amount := MoveFiles(path)
				if amount != 0 {
					log.Printf("moved %d files", amount)
				}
				isRunning = false
			}
		}
	}
}

func MoveFiles(path string) int {
	folderItems, _ := os.ReadDir(path)
	var movingCount int

	for _, item := range folderItems {
		if item.IsDir() {
			continue
		}

		itemInfo, err := item.Info()
		if err != nil {
			log.Printf("failed to fetch fileInfo: %s. skip moving", item.Name())
			continue
		}
		monthString := itemInfo.ModTime().Format("01")
		yearString := itemInfo.ModTime().Format("2006")

		monthPath := filepath.Join(path, yearString, monthString)
		err = os.MkdirAll(monthPath, os.ModePerm)
		if err != nil {
			log.Printf("failed to create folder %s", monthPath)
			continue
		}

		sourcePath := filepath.Join(path, item.Name())
		targetPath := filepath.Join(monthPath, item.Name())
		err = move(sourcePath, targetPath)
		if err != nil {
			log.Printf("failed to move file: %s", sourcePath)
			continue
		}
		movingCount++
	}
	return movingCount
}

func FolderExists(subFolders []fs.DirEntry, folderName string) bool {
	for _, folder := range subFolders {
		if folder.IsDir() && folder.Name() == folderName {
			return true
		}
	}
	return false
}

func move(filePath string, destPath string) error {
	log.Printf("moving file '%s' to '%s'", filePath, destPath)

	_, err := copy(filePath, destPath)
	if err != nil {
		log.Printf("copy file '%s' to '%s' failed: %v", filePath, destPath, err)
		return err
	}
	err = os.Remove(filePath)
	if err != nil {
		log.Printf("remove file '%s' failed: %v", filePath, err)
		return err
	}
	return nil
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func ReadEnv() {
	watchingFolders = os.Getenv("WATCH")
	scavangingInterval = os.Getenv("RUNNING_INTERVAL")

	if watchingFolders == "" {
		log.Fatalf("no folders to watch specified. set the WATCH environment variable. quit programm")
	}
	log.Printf("found WATCH folder: %s", watchingFolders)

	if scavangingInterval == "" {
		scavangingInterval = "1m"
		log.Printf("no interval specified. start moving files every minute\n")
	}
	log.Printf("found RUNNING_INTERVAL: %s", scavangingInterval)
}

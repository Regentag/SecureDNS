package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/jimlawless/whereami"
)

func WriteErrorLog(err error) {
	if err != nil {
		log.Printf("[ERROR] %s, %v", whereami.WhereAmI(2), err)
	}
}

func WriteErrorLogMsg(msg string, err error) {
	log.Printf("[ERROR] %s", msg)
	if err != nil {
		log.Printf("[ERROR] %s, %v", whereami.WhereAmI(2), err)
	}
}

func WriteErrorLogF(err error) {
	if err != nil {
		log.Fatalf("[ERROR] %s, %v", whereami.WhereAmI(2), err)
	}
}

func WriteErrorLogMsgF(msg string, err error) {
	log.Printf("[ERROR] %s", msg)
	if err != nil {
		log.Fatalf("[ERROR] %s, %v", whereami.WhereAmI(2), err)
	}
}

func GetMyDataPath(filename string) string {
	ex, err := os.Executable()
	if err != nil {
		ex = "." // current working directory
	}

	exPath := filepath.Dir(ex)
	return filepath.Join(exPath, filename)
}

func CheckFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

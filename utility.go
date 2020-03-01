package main

import (
	"log"

	"github.com/jimlawless/whereami"
)

func WriteErrorLog(err error) {
	if err != nil {
		log.Printf("[ERROR] %s, %v", whereami.WhereAmI(2), err)
	}
}

func WriteErrorLogMsg(msg string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s", msg)
		log.Printf("[ERROR] %s, %v", whereami.WhereAmI(2), err)
	}
}

func WriteErrorLogF(err error) {
	if err != nil {
		log.Fatalf("[ERROR] %s, %v", whereami.WhereAmI(2), err)
	}
}

func WriteErrorLogMsgF(msg string, err error) {
	if err != nil {
		log.Printf("[ERROR] %s", msg)
		log.Fatalf("[ERROR] %s, %v", whereami.WhereAmI(2), err)
	}
}

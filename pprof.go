package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		if err := http.ListenAndServe(":1234", nil); err != nil {
			log.Fatal("E! " + err.Error())
		}
	}()
}

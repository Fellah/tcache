// +build pprof

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Println(err)
		}
	}()
}

package main

import (
	"fmt"

	"github.com/fellah/stop"

	"github.com/fellah/tcache/db"
	"github.com/fellah/tcache/jobs"
)

func main() {
	go jobs.Start()

	<-stop.Ch

	db.Close()
	jobs.Stop()

	fmt.Println("done")
}

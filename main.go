package main

import (
	"fmt"

	"github.com/fellah/stop"

	"github.com/fellah/tcache/jobs"
)

func main() {
	jobs.Pipe()

	<-stop.Ch

	fmt.Println("done")
}
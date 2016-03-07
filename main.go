package main

import (
	"fmt"

	"github.com/fellah/stop"
	"github.com/fellah/tcache/jobs"
)

func main() {
	jobs.GetPacketList()

	<-stop.Ch

	fmt.Println("done")
}
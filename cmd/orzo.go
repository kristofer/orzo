package main

import (
	"log"
	"os"

	"github.com/kristofer/orzo"
)

func main() {
	// log setup....
	// f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	log.Fatalf("error opening file: %v", err)
	// }
	// defer f.Close()

	// log.SetOutput(f)
	// f.Truncate(0)
	// log.Println("Start of Log...")

	argv := os.Args // filename to edit
	argc := len(argv)
	if argc != 2 {
		log.Printf("Usage: orzo <filename>\n")
		return
	}

	Editor := &orzo.Orzo{}
	Editor.Start(argv[1])
}

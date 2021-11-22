package main

import (
	"log"
	"os"

	"github.com/carlescere/scheduler"

	"github.com/lukasdietrich/dockerized/updater/internal/meta"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	m, err := meta.Open()
	if err != nil {
		log.Fatalf("could not open repository: %v", err)
	}

	scheduler.
		Every().Day().At(os.Getenv("DOCKERIZED_SCHEDULE")).
		Run(func() {
			if err := m.UpdateVersions(); err != nil {
				log.Printf("error during update: %v", err)
			}
		})

	<-make(chan int)
}

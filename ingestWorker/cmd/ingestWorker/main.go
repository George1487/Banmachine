package main

import (
	"log"

	"ingestWorker/internal/config"
	"ingestWorker/internal/worker"
)

func main() {
	if err := worker.Run(config.MustConfig()); err != nil {
		log.Fatal(err)
	}
}

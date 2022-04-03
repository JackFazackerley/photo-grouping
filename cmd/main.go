package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/JackFazackerley/photo-grouping/internal/categoriser"
	"github.com/JackFazackerley/photo-grouping/internal/consumer"
	"github.com/JackFazackerley/photo-grouping/internal/heap"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"googlemaps.github.io/maps"
)

var (
	csvPath string
	apiKey  string
)

func init() {
	flag.StringVar(&csvPath, "csvPath", "", "path to csv")
	flag.StringVar(&apiKey, "apiKey", "", "apiKey required for Google's Reverse Geocoding API")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	reader, err := consumer.NewReader(csvPath)
	if err != nil {
		log.WithError(err).Fatal("creating new reader")
	}
	defer reader.Close()

	photosChan := reader.ReadCSV(ctx)

	consumer, err := consumer.NewConsumer(apiKey, maps.WithRateLimit(50))
	if err != nil {
		log.WithError(err).Fatal("creating consumer")
	}

	go func() {
		c := make(chan os.Signal, 1)

		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	photoHeap := heap.New()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go consumer.Run(ctx, photoHeap, photosChan, wg)
	}

	wg.Wait()

	locations := categoriser.Group(photoHeap)

	for _, location := range locations {
		suggestions := location.GenerateTitles()
		for _, suggestion := range suggestions {
			log.Println(suggestion)
		}
	}
}

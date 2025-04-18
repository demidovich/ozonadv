package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"ozonadv/internal/app"
	"ozonadv/internal/ui"
	"syscall"
)

func main() {
	log.SetFlags(0)

	app := app.New(os.Stdout)
	defer app.Shutdown()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		app.Shutdown()
		os.Exit(1)
	}()

	err := ui.Home(app)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
}

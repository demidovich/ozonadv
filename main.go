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

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		app.Shutdown()
		os.Exit(1)
	}()

	err := ui.Home(app)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
}

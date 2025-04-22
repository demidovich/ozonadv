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

	appInstance := app.New(os.Stdout)
	defer appInstance.Shutdown()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		appInstance.Shutdown()
		os.Exit(1)
	}()

	err := ui.Home(appInstance)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("")
}

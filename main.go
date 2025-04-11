package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"ozonadv/internal/app"
	"ozonadv/pkg/prompts"
	"syscall"
)

func main() {
	log.SetFlags(0)
	defer fmt.Println("")

	app := app.New(os.Stdout)
	defer app.Shutdown()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		app.Shutdown()
		os.Exit(1)
	}()

	action := mainMenuAction()

	fmt.Println(action)
}

func mainMenuAction() string {
	options := map[string]string{
		"Кабинеты":   "cabinets",
		"Статистика": "stats",
		"Выход":      "quit",
	}

	value, _ := prompts.List("---", options)

	return value
}

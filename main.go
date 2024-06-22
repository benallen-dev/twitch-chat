package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

func main() {
	// Hide cursor and clear screen
	fmt.Print("\033[?25l")
	fmt.Print("\033[2J\033[H")

	// Default to my twitch channel
	var channel = "stevemacawesome"
	if len(os.Args) > 1 {
		channel = os.Args[1]
	}

	client := twitch.NewAnonymousClient()
	client.Join(channel)

	client.OnConnect(func() {
		fmt.Println("Connected to", channel)
	})

	// Private message is a misnomer, it's actually a chat message
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		usr := message.User.DisplayName
		msg := message.Message
		col := message.User.Color
		r, g, b := int64(49), int64(154), int64(36) // Same default as Twitch
		if col != "" {
			r, _ = strconv.ParseInt(col[1:3], 16, 64)
			g, _ = strconv.ParseInt(col[3:5], 16, 64)
			b, _ = strconv.ParseInt(col[5:7], 16, 64)
		}

		colorCode := fmt.Sprintf("\033[0m\033[38;2;%d;%d;%dm", r, g, b)

		fmt.Printf("%s%s\033[0m: %s\n", colorCode, usr, msg)
	})

	// Restore cursor on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print("\033[?25h") // Show cursor
		os.Exit(1)
	}()

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

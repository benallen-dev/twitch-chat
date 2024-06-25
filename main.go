package main

/*
https://twitchapps.com/tmi/
*/

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"golang.org/x/term"

	twitch "github.com/gempir/go-twitch-irc/v4"
	"github.com/muesli/reflow/wordwrap"
)

func displayName(User twitch.User) string {
	name := User.DisplayName
	col := User.Color

	r, g, b := int64(185), int64(163), int64(227)
	if col != "" {
		r, _ = strconv.ParseInt(col[1:3], 16, 64)
		g, _ = strconv.ParseInt(col[3:5], 16, 64)
		b, _ = strconv.ParseInt(col[5:7], 16, 64)
	}

	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[0m", r, g, b, name)
}

func getTermWidth() int {

	if !term.IsTerminal(0) {
		return -1
    }

	width, _, err := term.GetSize(0)
    if err != nil {
        return -1
    }

	return width
}

func main() {
	// Hide cursor and clear screen
	fmt.Print("\033[?25l\033[2J\033[H")

	// Default to my twitch channel
	var channel = "stevemacawesome"
	if len(os.Args) > 1 {
		channel = os.Args[1]
	}

	user := os.Getenv("TWITCH_USER")
	oauth := os.Getenv("TWITCH_OAUTH")

	//client := twitch.NewAnonymousClient()
	client := twitch.NewClient(user, oauth)
	client.Join(channel)

	client.OnSelfJoinMessage(func(message twitch.UserJoinMessage) {
		fmt.Printf("%s: Joined %s\n", message.User, message.Channel)
	})

	// Private message is a misnomer, it's actually a chat message but it's how Twitch IRC works
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		s := displayName(message.User) + ": " + message.Message
		s = wordwrap.String(s, getTermWidth())
		fmt.Println(s)
	})

	// Restore cursor on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		client.Disconnect()
		fmt.Print("\033[?25h")
		os.Exit(1)
	}()

	fmt.Println("Attempting to connect...")
	err := client.Connect()
	if err != nil {
		panic(err)
	}

	// The application will not reach here until it is terminated
	// because client.Connect() is a blocking call

}

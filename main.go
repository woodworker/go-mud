package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"github.com/woodworker/go-mud/game"
)

func main() {
	workingdir, _ := os.Getwd()

	log.Printf("Leveldir %s", workingdir + "/static/levels/")

	server := game.NewServer("berlin-mud", workingdir)
	server.LoadLevels()
	log.Printf("%v", server)

	ln, err := net.Listen("tcp", ":1337")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msgchan := make(chan string)
	addchan := make(chan game.Client)
	rmchan := make(chan game.Client)

	go handleMessages(msgchan, addchan, rmchan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn, msgchan, addchan, rmchan, server)
	}
}


func promptNick(c net.Conn, bufc *bufio.Reader) string {
	io.WriteString(c, "What is your nick? ")
	nick, _, _ := bufc.ReadLine()
	return string(nick)
}

func promptMessage(c net.Conn, bufc *bufio.Reader, message string) string {
	io.WriteString(c, message)
	nick, _, _ := bufc.ReadLine()
	return string(nick)
}

func handleConnection(c net.Conn, msgchan chan <- string, addchan chan <- game.Client, rmchan chan <- game.Client, server *game.Server) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	io.WriteString(c, fmt.Sprintf("\033[1;30;41mWelcome to the Go-Mud Server %s!\033[0m\n\r", server.GetName()))

	var nickname string
	for {
		nickname = promptNick(c, bufc)
		ok := server.LoadPlayer(nickname)

		if ok == false {
			io.WriteString(c, fmt.Sprintf("Username %s does not exists. Should i create it?", nickname))

		}

		if ok == true {
			break
		}
	}

	player, playerLoaded := server.GetPlayerByNick(nickname)

	if !playerLoaded {
		log.Println("problem getting user object")
		io.WriteString(c, "Problem getting user object\n")
		return
	}

	client := game.NewClient(c, player)

	if strings.TrimSpace(client.Nickname) == "" {
		log.Println("invalid username")
		io.WriteString(c, "Invalid Username\n")
		return
	}

	// Register user
	addchan <- client
	defer func() {
		msgchan <- fmt.Sprintf("User %s left the chat room.\n\r", client.Nickname)
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmchan <- client
	}()
	io.WriteString(c, fmt.Sprintf("Welcome, %s!\n\n\r", client.Nickname))

	location, locationLoaded := server.GetRoom(client.Player.Position);

	if locationLoaded {
		location.OnEnterRoom(server, client)
	}

	//msgchan <- fmt.Sprintf("New user %s has joined the chat room.\n\r", client.Nickname)

	// I/O
	go client.ReadLinesInto(msgchan, server)
	client.WriteLinesFrom(client.Ch)
}

func handleMessages(msgchan <-chan string, addchan <-chan game.Client, rmchan <-chan game.Client) {
	clients := make(map[net.Conn]chan <- string)

	for {
		select {
		case msg := <-msgchan:
			log.Printf("New message: %s", msg)
			for _, ch := range clients {
				go func(mch chan <- string) { mch <- "\033[1;33;40m" + msg + "\033[m\n\r\n\r" }(ch)
			}
		case client := <-addchan:
			log.Printf("New client: %v\n\r\n\r", client.Conn)
			clients[client.Conn] = client.Ch
		case client := <-rmchan:
			log.Printf("Client disconnects: %v\n\r\n\r", client.Conn)
			delete(clients, client.Conn)
		}
	}
}

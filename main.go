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

type Client struct {
	conn     	 net.Conn
	nickname	 string
	player		 game.Player
	ch       	 chan string
}

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
	addchan := make(chan Client)
	rmchan := make(chan Client)

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

func (c Client) ReadLinesInto(ch chan <- string, server *game.Server) {
	bufc := bufio.NewReader(c.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}

		userLine := strings.TrimSpace(line)

		if(userLine==""){
			continue
		}

		//io.WriteString(c.conn, fmt.Sprintf("You wrote: %s\n\r", userLine))
		lineParts := strings.SplitN(userLine, " ", 2)

		var command, commandText string
		if(len(lineParts)>0){
			command = lineParts[0]
		}
		if(len(lineParts)>1){
			commandText = lineParts[1]
		}

		log.Printf("Command: %s  -  %s", command, commandText)

		switch command {
		case "watch":
			place, ok := server.GetRoom(c.player.Position)
			if ok {
				io.WriteString(c.conn, fmt.Sprintf("You are at \033[1;30;41m%s\033[0m\n\r", place.Name))
				for _,oneDirection := range place.Directions {
					place, ok := server.GetRoom(oneDirection.Station)
					if ok {
						io.WriteString(c.conn, fmt.Sprintf("When you look %s you see %s\n\r", oneDirection.Direction, place.Name))
					}
				}
			}
		case "go":
			place, ok := server.GetRoom(c.player.Position)
			if ok {
				for _,oneDirection := range place.Directions {
					if(oneDirection.Direction == commandText){
						place, ok := server.GetRoom(oneDirection.Station)
						if ok {
							io.WriteString(c.conn, fmt.Sprintf("You are now at \033[1;30;41m%s\033[0m\n\r", place.Name))
							c.player.Position = place.Key
							server.SavePlayer(c.player)
						}
					}
				}
			}
		case "say":
			ch <- fmt.Sprintf("%s: %s", c.player.Gamename, commandText)
		}
	}
}

func (c Client) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		_, err := io.WriteString(c.conn, msg)
		if err != nil {
			return
		}
	}
}

func promptNick(c net.Conn, bufc *bufio.Reader) string {
	io.WriteString(c, "What is your nick? ")
	nick, _, _ := bufc.ReadLine()
	return string(nick)
}

func handleConnection(c net.Conn, msgchan chan <- string, addchan chan <- Client, rmchan chan <- Client, server *game.Server) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	io.WriteString(c, fmt.Sprintf("\033[1;30;41mWelcome to the Go-Mud Server %s!\033[0m\n\r", server.GetName()))

	var nickname string
	for {
		nickname = promptNick(c, bufc)
		ok := server.LoadPlayer(nickname)
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

	client := Client{
		conn:     c,
		nickname: player.Nickname,
		player:   player,
		ch:       make(chan string),
	}

	if strings.TrimSpace(client.nickname) == "" {
		log.Println("invalid username")
		io.WriteString(c, "Invalid Username\n")
		return
	}

	// Register user
	addchan <- client
	defer func() {
		msgchan <- fmt.Sprintf("User %s left the chat room.\n\r", client.nickname)
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmchan <- client
	}()
	io.WriteString(c, fmt.Sprintf("Welcome, %s!\n\n\r", client.nickname))

	location, locationLoaded := server.GetRoom(client.player.Position);

	if locationLoaded {
		io.WriteString(c, fmt.Sprintf("You are at: \033[1;33;40m%s\033[m\n\n\r", location.Name))
	}

	msgchan <- fmt.Sprintf("New user %s has joined the chat room.\n\r", client.nickname)

	// I/O
	go client.ReadLinesInto(msgchan, server)
	client.WriteLinesFrom(client.ch)
}

func handleMessages(msgchan <-chan string, addchan <-chan Client, rmchan <-chan Client) {
	clients := make(map[net.Conn]chan <- string)

	for {
		select {
		case msg := <-msgchan:
			log.Printf("New message: %s", msg)
			for _, ch := range clients {
				go func(mch chan <- string) { mch <- "\033[1;33;40m" + msg + "\033[m\n\r\n\r" }(ch)
			}
		case client := <-addchan:
			log.Printf("New client: %v\n\r\n\r", client.conn)
			clients[client.conn] = client.ch
		case client := <-rmchan:
			log.Printf("Client disconnects: %v\n\r\n\r", client.conn)
			delete(clients, client.conn)
		}
	}
}

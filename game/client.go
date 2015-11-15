package game

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Client struct {
	Conn     net.Conn
	Nickname string
	Player   Player
	Ch       chan string
}

func NewClient(c net.Conn, player Player) Client {
	return Client{
		Conn:     c,
		Nickname: player.Nickname,
		Player:   player,
		Ch:       make(chan string),
	}
}

func (c Client) WriteToUser(msg string) {
	io.WriteString(c.Conn, msg)
}

func (c Client) WriteLineToUser(msg string) {
	io.WriteString(c.Conn, msg + "\n\r")
}

func (c Client) ReadLinesInto(ch chan<- string, server *Server) {
	bufc := bufio.NewReader(c.Conn)

	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}

		userLine := strings.TrimSpace(line)

		if userLine == "" {
			continue
		}

		//io.WriteString(c.Conn, fmt.Sprintf("You wrote: %s\n\r", userLine))
		lineParts := strings.SplitN(userLine, " ", 2)

		var command, commandText string
		if len(lineParts) > 0 {
			command = lineParts[0]
		}
		if len(lineParts) > 1 {
			commandText = lineParts[1]
		}

		log.Printf("Command by %s: %s  -  %s", c.Player.Nickname, command, commandText)

		switch command {
		case "look":
			fallthrough
		case "watch":
			place, ok := server.GetRoom(c.Player.Position)
			if ok {
				c.WriteLineToUser(fmt.Sprintf("You are at \033[1;30;41m%s\033[0m", place.Name))
				for _, direction := range place.Directions {
					place, ok := server.GetRoom(direction.Station)
					if (ok && place.CanSeeDirection(direction, c.Player, commandText)) {
						c.WriteLineToUser(fmt.Sprintf(" When you look %s you see %s", direction.Direction, place.Name))
					}
				}
			}
		case "go":
			place, ok := server.GetRoom(c.Player.Position)
			if ok {
				for _, oneDirection := range place.Directions {
					if strings.ToLower(oneDirection.Direction) == strings.ToLower(commandText) {
						place, ok := server.GetRoom(oneDirection.Station)
						if ok {
							canEnter, message := place.CanGoDirection(oneDirection, c.Player)
							if !canEnter {
								c.WriteLineToUser(" >" + message)
							} else {
								place.OnEnterRoom(server, c)
								c.Player.Position = string(place.Key)
								log.Println(c.Player)
								c.Player.LogAction(place.Key)
								server.SavePlayer(c.Player)
							}
						} else {
							c.WriteToUser("\n")
						}
					}
				}
			}
		case "say":
			// TODO: implement channel wide communication
			io.WriteString(c.Conn, "\033[F") // up one line so we overwrite the say command typed with the result
			ch <- fmt.Sprintf("%s: %s", c.Player.Gamename, commandText)
		case "quit":
			fallthrough
		case "leave":
			fallthrough
		case "exit":
			server.OnExit(c)
			c.Conn.Close()
		default:
			place, ok := server.GetRoom(c.Player.Position)
			if ok {
				action, ok := place.GetRoomAction(command)
				if ok {
					actionName := place.GetRoomActionName(action)
					if action.Dependency == "" || c.Player.HasAction(action.Dependency) {
						c.Player.LogAction(actionName)
						server.SavePlayer(c.Player)
						c.WriteLineToUser(fmt.Sprintf(" > %s", action.Answer))
						continue
					}
				}
			}
		}
	}
}

func (c Client) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		_, err := io.WriteString(c.Conn, msg)
		if err != nil {
			return
		}
	}
}

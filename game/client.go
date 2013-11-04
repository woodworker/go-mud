package game

import (
	"net"
	"io"
	"fmt"
	"strings"
	"log"
	"bufio"
)


type Client struct {
	Conn     	 net.Conn
	Nickname	 string
	Player		 Player
	Ch       	 chan string
}

func NewClient(c net.Conn, player Player) Client {
	return Client{
		Conn:     c,
		Nickname: player.Nickname,
		Player:   player,
		Ch:       make(chan string),
	}
}


func (c Client) ReadLinesInto(ch chan <- string, server *Server) {
	bufc := bufio.NewReader(c.Conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}

		userLine := strings.TrimSpace(line)

		if(userLine==""){
			continue
		}

		//io.WriteString(c.Conn, fmt.Sprintf("You wrote: %s\n\r", userLine))
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
		case "look": fallthrough
		case "watch":
			place, ok := server.GetRoom(c.Player.Position)
			if ok {
				io.WriteString(c.Conn, fmt.Sprintf("You are at \033[1;30;41m%s\033[0m\n\r", place.Name))
				for _,oneDirection := range place.Directions {
					place, ok := server.GetRoom(oneDirection.Station)
					if ok && ((oneDirection.Hidden == "" && commandText == "") || strings.ToLower(oneDirection.Direction) == strings.ToLower(commandText)) {
						io.WriteString(c.Conn, fmt.Sprintf("When you look %s you see %s\n\r", oneDirection.Direction, place.Name))
					}
				}
			}
		case "go":
			place, ok := server.GetRoom(c.Player.Position)
			if ok {
				for _,oneDirection := range place.Directions {
					if(strings.ToLower(oneDirection.Direction) == strings.ToLower(commandText)){
						place, ok := server.GetRoom(oneDirection.Station)
						if ok {
							place.OnEnterRoom(server, c)
							c.Player.Position = string(place.Key)
							log.Println(c.Player)
							server.SavePlayer(c.Player)
						}
					}
				}
			}
		case "say":
			ch <- fmt.Sprintf("%s: %s", c.Player.Gamename, commandText)
		case "exit":
			server.OnExit(c)
			c.Conn.Close()
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

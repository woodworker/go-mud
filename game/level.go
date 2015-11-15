package game

import (
	"io"
	"fmt"
)

type Level struct {
	Key			string		`xml:"key,attr"`
	Tag			string		`xml:"tag,attr"`
	Name		string		`xml:"name"`
	Directions	[]Direction	`xml:"directions>direction"`
	Intro		string		`xml:"intro"`
}

type Direction struct {
	Station			string		`xml:"station,attr"`
	Hidden			string		`xml:"hidden,attr"`
	Direction		string      `xml:",chardata"`
}

func (l *Level) OnEnterRoom(s *Server, c Client) {
	io.WriteString(c.Conn, fmt.Sprintf("You are now at \033[1;30;41m%s\033[0m\n\r", l.Name))
	if l.Intro != "" {
		io.WriteString(c.Conn, fmt.Sprintf(" > %s\n\r", l.Intro))
	}
}

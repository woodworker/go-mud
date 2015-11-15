package game

import (
	"fmt"
	"log"
)

type Level struct {
	Key			string		`xml:"key,attr"`
	Tag			string		`xml:"tag,attr"`
	Name		string		`xml:"name"`
	Directions	[]Direction	`xml:"directions>direction"`
	Actions	    []Action    `xml:"actions>action"`
	Intro		string		`xml:"intro"`
}

type Action struct {
	Name		string		`xml:"name,attr"`
	Hidden		string		`xml:"hidden,attr"`
	Answer		string      `xml:",chardata"`
}

type Direction struct {
	Station			string		`xml:"station,attr"`
	Hidden			string		`xml:"hidden,attr"`
	Direction		string      `xml:",chardata"`
}

func (l *Level) OnEnterRoom(s *Server, c Client) {
	c.WriteToUser(fmt.Sprintf("You are now at \033[1;30;41m%s\033[0m\n\r", l.Name))
	if l.Intro != "" {
		c.WriteToUser(fmt.Sprintf(" > %s\n\r", l.Intro))
	}
}

func (l *Level) GetRoomAction(command string) (Action, bool) {

	if len(l.Actions) > 0 {
		for _, a := range l.Actions {
			if a.Name == command {
				log.Println(fmt.Sprintf("Found Action: %", a.Name))
				return a, true;
			}
		}
	}

	return Action{}, false;
}

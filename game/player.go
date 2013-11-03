package game

import "encoding/xml"

type Player struct {
	XMLName 	xml.Name	`xml:"player"`
	Nickname	string		`xml:"nickname,attr"`
	Gamename	string      `xml:"name"`
	Position	string		`xml:"position,attr"`
	PlayerType	string		`xml:"type"`
	Ch       	chan string
}

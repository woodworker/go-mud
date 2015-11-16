package game

import (
	"encoding/xml"
	"strings"
)

type Player struct {
	XMLName    xml.Name    `xml:"player"`
	Nickname   string      `xml:"nickname,attr"`
	Gamename   string      `xml:"name"`
	Position   string      `xml:"position,attr"`
	PlayerType string      `xml:"type"`
	Ch         chan string `xml:"-"`
	ActionLog  []string    `xml:"actions>action"`
}

func (p *Player) LogAction(action string) {
	if !p.HasAction(action) {
		p.ActionLog = append(p.ActionLog, strings.ToLower(action))
	}
}

func (p *Player) HasAction(action string) bool {
	for _, a := range p.ActionLog {
		if strings.ToLower(a) == strings.ToLower(action) {
			return true
		}
	}
	return false
}

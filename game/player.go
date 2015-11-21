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
	Attributes []Attribute `xml:"attributes>attribute"`
}

type Attribute struct {
	Name   string   `xml:"name,attr"`
	Value  int64    `xml:"value"`
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

func (p *Player) GetAttribute(name string) int64 {
	for _, a := range p.Attributes {
		if strings.ToLower(a.Name) == strings.ToLower(name) {
			return a.Value
		}
	}
	return 0
}

func (p *Player) UpdateAttribute(name string, update int64) {
	for i := range p.Attributes {
		if strings.ToLower(p.Attributes[i].Name) == strings.ToLower(name) {
			p.Attributes[i].Value = p.Attributes[i].Value + update
			if p.Attributes[i].Value <= 0 {
				p.Attributes[i].Value = 0
			}
			return
		}
	}

	p.Attributes = append(p.Attributes, Attribute{
		Name: name,
		Value: update,
	})
}

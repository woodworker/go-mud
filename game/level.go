package game

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Level struct {
	Key         string      `xml:"key,attr"`
	Tag         string      `xml:"tag,attr"`
	Name        string      `xml:"name"`
	Directions  []Direction `xml:"directions>direction"`
	Actions     []Action    `xml:"actions>action"`
	Intro       string      `xml:"intro"`
	Asciimation Asciimation `xml:"asciimation"`
}

type Asciimation struct {
	Frames []Frame `xml:"frame"`
}

type Frame struct {
	Id       int      `xml:"id,attr"`
	Duration int      `xml:"duration,attr"`
	Lines    []string `xml:"line"`
}

type Action struct {
	Name         string       `xml:"name,attr"`
	Hidden       string       `xml:"hidden,attr"`
	Dependencies []Dependency `xml:"dependency"`
	Answer       string       `xml:"answer"`
}

type Direction struct {
	Station      string       `xml:"station"`
	Hidden       bool         `xml:"hidden,attr"`
	Dependencies []Dependency `xml:"dependency"`
	Direction    string       `xml:"name"`
}

type Dependency struct {
	Key         string `xml:"key,attr"`
	OkMessage   string `xml:"okMessage"`
	FailMessage string `xml:"failMessage"`
}

func (l *Level) OnEnterRoom(s *Server, c Client) {

	c.WriteToUser("┌────────────");
	runeLen := len([]rune(l.Name))
	for i := 0; i < runeLen; i++ {
		c.WriteToUser("─");
	}
	c.WriteToUser("─┐\n\r");
	c.WriteLineToUser(fmt.Sprintf("│ You are at \033[1;30;41m%s\033[0m │", l.Name))
	c.WriteToUser("└────────────")
	for i := 0; i < runeLen; i++ {
		c.WriteToUser("─");
	}
	c.WriteToUser("─┘\n\r");

	if len(l.Asciimation.Frames) > 0 {
		l.Asciimation.Play(c)
	}

	if l.Intro != "" {
		c.WriteToUser(fmt.Sprintf(" > %s\n\r", l.Intro))
	}
}

func (a *Asciimation) Play(c Client) {

	lineCount := 0
	frameCounter := 0
	for _, f := range a.Frames {
		frameCounter++
		if lineCount == 0 {
			lineCount = len(f.Lines)
		}

		i := 0
		for _, l := range f.Lines {
			i++
			if frameCounter > 1 && i == 1 {
				c.WriteToUser(fmt.Sprintf("\033[%dF\033[K%s\n\r", lineCount, l))
			} else {
				c.WriteToUser(fmt.Sprintf("\033[K%s\n\r", l))
			}
			if i == lineCount {
				break
			}
		}
		if frameCounter < len(a.Frames) {
			duration := time.Duration(f.Duration) * time.Millisecond
			time.Sleep(duration)
		}
	}

}

func (l *Level) GetRoomAction(command string) (Action, bool) {

	if len(l.Actions) > 0 {
		for _, a := range l.Actions {
			if a.Name == command {
				log.Println(fmt.Sprintf("Found Action: %", a.Name))
				return a, true
			}
		}
	}

	return Action{}, false
}

func (l *Level) GetRoomActionName(action Action) string {
	return fmt.Sprintf("%s:%s", l.Key, action.Name)
}

func (l *Level) CanDoAction(action Action, player Player) (bool, string) {
	if len(action.Dependencies) == 0 {
		return true, action.Answer
	}

	lastOkMessage := "";
	for _, d := range action.Dependencies {
		if !player.HasAction(d.Key) {
			return false, d.FailMessage
		} else {
			lastOkMessage = d.OkMessage
		}
	}
	if action.Answer != "" {
		lastOkMessage = action.Answer
	}
	return true, lastOkMessage
}

func (l *Level) CanSeeDirection(direction Direction, player Player, viewDirection string) bool {
	if viewDirection != "" {
		return strings.ToLower(direction.Direction) == strings.ToLower(viewDirection)
	}
	if player.HasAction(direction.Station) {
		return true
	}
	if direction.Hidden {
		return false
	}
	return true
}

func (l *Level) CanGoDirection(direction Direction, player Player) (bool, string) {
	if len(direction.Dependencies) == 0 {
		return true, ""
	}

	for _, d := range direction.Dependencies {
		if !player.HasAction(d.Key) {
			return false, d.FailMessage
		}
	}
	return true, ""
}

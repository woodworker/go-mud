package game

import (
	"fmt"
	"log"
	"strings"
	"time"
	"strconv"
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
	Type        string `xml:"type,attr"`
	MinValue    string `xml:"minValue"`
	MaxValue    string `xml:"maxValue"`
	OkMessage   string `xml:"okMessage"`
	FailMessage string `xml:"failMessage"`
}

func (l *Level) OnEnterRoom(s *Server, c Client) {

	c.WriteToUser("┌────────────")
	runeLen := len([]rune(l.Name))
	for i := 0; i < runeLen; i++ {
		c.WriteToUser("─")
	}
	c.WriteToUser("─┐\n\r")
	c.WriteLineToUser(fmt.Sprintf("│ You are at \033[1;30;41m%s\033[0m │", l.Name))
	c.WriteToUser("└────────────")
	for i := 0; i < runeLen; i++ {
		c.WriteToUser("─")
	}
	c.WriteToUser("─┘\n\r")

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
				log.Println(fmt.Sprintf("Found Action: %s", a.Name))
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
	return CheckDependencies(action.Dependencies, player, action.Answer)
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
	return CheckDependencies(direction.Dependencies, player, "")
}

func CheckDependencies(dependencies []Dependency, player Player, defaultAnswer string) (bool, string) {
	if len(dependencies) == 0 {
		return true, defaultAnswer
	}

	lastOkMessage := ""
	for _, d := range dependencies {
		switch d.Type {
		case "":
			fallthrough
		case "action":
			if !player.HasAction(d.Key) {
				return false, d.FailMessage
			}

			lastOkMessage = d.OkMessage
		case "attribute":
			playerAttribute := player.GetAttribute(d.Key)

			minValue, errMin := strconv.ParseInt(d.MinValue, 10, 64)
			if errMin==nil && d.MinValue != "" && playerAttribute < minValue {
				return false, d.FailMessage
			}

			maxValue, errMax := strconv.ParseInt(d.MinValue, 10, 64)
			if errMax==nil && d.MaxValue != "" && playerAttribute > maxValue {
				return false, d.FailMessage
			}

			lastOkMessage = d.OkMessage
		case "time":
			if strings.Count(d.MinValue, ":")==1 && strings.Count(d.MaxValue, ":")==1 {
				minParts := strings.SplitN(d.MinValue, ":", 2)
				maxParts := strings.SplitN(d.MaxValue, ":", 2)

				nowHour := int64(time.Now().Hour());
				nowMinute := int64(time.Now().Minute())

				minHour, minHourErr := strconv.ParseInt(minParts[0], 10, 64)
				minMinute, minMinuteErr := strconv.ParseInt(minParts[1], 10, 64)

				if minHourErr != nil || minMinuteErr != nil || minHour > 23 || minMinute > 59 {
					return false, d.FailMessage
				}

				if minHour > nowHour || (minHour == nowHour && minMinute > nowMinute) {
					return false, d.FailMessage
				}

				maxHour, maxHourErr := strconv.ParseInt(maxParts[0], 10, 64)
				maxMinute, maxMinuteErr := strconv.ParseInt(maxParts[1], 10, 64)

				if maxHourErr != nil || maxMinuteErr != nil || maxHour > 23 || maxMinute > 59 {
					return false, d.FailMessage
				}

				if maxHour < nowHour || (maxHour == nowHour && maxMinute < nowMinute) {
					return false, d.FailMessage
				}
				return true, d.OkMessage
			}
			return false, d.FailMessage
		case "date":
			from, fromError := time.Parse("2006-01-02 15:04:05", d.MinValue+" 00:00:00")
			to, toError := time.Parse("2006-01-02 15:04:05", d.MaxValue+" 23:59:59")
			now := time.Now();
			if (fromError == nil && toError == nil) {
				diffFrom := now.Sub(from)
				diffTo := to.Sub(now)

				if diffFrom.Seconds() < 0 {
					return false, d.FailMessage
				}
				if diffTo.Seconds() < 0 {
					return false, d.FailMessage
				}
				return true, d.OkMessage
			} else {
				return false, d.FailMessage
			}
		}
	}
	if defaultAnswer != "" {
		lastOkMessage = defaultAnswer
	}
	return true, lastOkMessage
}

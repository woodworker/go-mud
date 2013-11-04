package game

import (
	"encoding/xml"
	"path/filepath"
	"os"
	"io"
	"fmt"
	"io/ioutil"
	"log"
)

type Server struct {
	name		string
	players		map[string]Player
	levels		map[string]Level
	workingdir	string
}

func NewServer(servername string, serverdir string) *Server {
	return &Server{
		name: servername,
		players: make(map[string]Player),
		levels: make(map[string]Level),
		workingdir: serverdir,
	}
}

func (s *Server) LoadLevels() error {
	log.Println("Loading levels ...")
	levelWalker := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileContent, fileIoErr := ioutil.ReadFile(path)
		if fileIoErr != nil {
			log.Printf("\n")
			log.Printf("File %s could not be loaded\n", path)
			log.Printf("%v",fileIoErr)
			//return fileIoErr
			return nil
		}
		level := Level{}
		if xmlerr := xml.Unmarshal(fileContent, &level); xmlerr != nil {
			log.Printf("\n")
			log.Printf("File %s could not be Unmarshaled\n", path, xmlerr)
			log.Printf("%v",xmlerr)
			//return xmlerr
			return nil
		}
		log.Printf(" loaded: %s", info.Name())
		s.addLevel(level)
		return nil
	}

	filepath.Walk(s.workingdir+"/static/levels/", levelWalker)
	return nil
}

func (s *Server) getPlayerFileName(playerName string) string {
	return s.workingdir+"/static/player/"+playerName+".player"
}

func (s *Server) LoadPlayer(playerName string) bool {
	playerFileName := s.getPlayerFileName(playerName)

	log.Println("Loading player %s", playerFileName)

	fileContent, fileIoErr := ioutil.ReadFile(playerFileName)
	if fileIoErr != nil {
		log.Printf("\n")
		log.Printf("File %s could not be loaded\n", playerFileName)
		log.Printf("%v",fileIoErr)
		//return fileIoErr
		return false
	}

	player := Player{}
	if xmlerr := xml.Unmarshal(fileContent, &player); xmlerr != nil {
		log.Printf("\n")
		log.Printf("File %s could not be Unmarshaled\n", playerFileName, xmlerr)
		log.Printf("%v",xmlerr)
		//return xmlerr
		return false
	}
	log.Printf(" loaded: %s", player.Gamename)
	s.addPlayer(player)

	return true
}

func (s *Server) addLevel(level Level) error {
	s.levels[level.Key] = level
	return nil
}

func (s *Server) addPlayer(player Player) error {
	s.players[player.Nickname] = player
	return nil
}

func (s *Server) GetPlayerByNick(nickname string) (Player, bool) {
	player,ok := s.players[nickname]
	return player, ok
}

func (s *Server) GetRoom(key string) (Level, bool) {
	level,ok := s.levels[key]
	return level, ok
}

func (s *Server) GetName() string {
	return s.name
}

func (s *Server) SavePlayer(player Player) (bool) {
	data, err := xml.Marshal(player)
	if err == nil {
		playerFileName := s.getPlayerFileName(player.Nickname)
		if ioerror := ioutil.WriteFile(playerFileName, data, 0666); ioerror != nil {
			log.Println(ioerror)
			return true
		}
	} else {
		log.Println(err)
	}
	return false
}

func (s *Server) OnExit(client Client) {
	s.SavePlayer(client.Player)
	io.WriteString(client.Conn, fmt.Sprintf("Good bye %s", client.Player.Nickname))
}

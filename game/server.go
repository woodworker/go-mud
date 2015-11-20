package game

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type Server struct {
	players      map[string]Player
	levels       map[string]Level
	workingdir   string
	DefaultLevel Level
	Config       ServerConfig
}

type ServerConfig struct {
	Name      string `xml:"name"`
	Interface string `xml:"interface"`
	Motd      string `xml:"motd"`
}

func (s *Server) HasDefaultLevel() bool {
	return s.DefaultLevel.Key != ""
}

func NewServer(serverdir string) *Server {
	server := &Server{
		players:    make(map[string]Player),
		levels:     make(map[string]Level),
		workingdir: serverdir,
	}

	server.LoadConfig()

	return server
}

func (s *Server) LoadConfig() error {
	log.Println("Loading config ...")
	configFileName := s.workingdir + "/static/server.xml"
	fileContent, fileIoErr := ioutil.ReadFile(configFileName)
	if fileIoErr != nil {
		log.Printf("\n")
		log.Printf("File %s could not be loaded\n", configFileName)
		log.Printf("%v", fileIoErr)
		return fileIoErr
	}
	config := ServerConfig{}
	if xmlerr := xml.Unmarshal(fileContent, &config); xmlerr != nil {
		log.Printf("\n")
		log.Printf("File %s could not be Unmarshaled\n", configFileName, xmlerr)
		log.Printf("%v", xmlerr)
		return xmlerr
	}
	s.Config = config
	log.Println(" config loaded")
	return nil
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
			log.Printf("%v", fileIoErr)
			return fileIoErr
		}
		level := Level{}
		if xmlerr := xml.Unmarshal(fileContent, &level); xmlerr != nil {
			log.Printf("\n")
			log.Printf("File %s could not be Unmarshaled\n", path, xmlerr)
			log.Printf("%v", xmlerr)
			return xmlerr
		}
		log.Printf(" loaded: %s\n", info.Name())
		s.addLevel(level)
		return nil
	}

	return filepath.Walk(s.workingdir+"/static/levels/", levelWalker)
}

func (s *Server) getPlayerFileName(playerName string) (bool, string) {
	if !s.IsValidUsername(playerName) {
		return false, ""
	}
	return true, s.workingdir + "/static/player/" + playerName + ".player"
}

func (s *Server) IsValidUsername(playerName string) bool {
	r, err := regexp.Compile(`^[a-zA-Z0-9_-]{1,40}$`)
	if err != nil {
		return false
	}
	if !r.MatchString(playerName) {
		return false
	}
	return true
}

func (s *Server) LoadPlayer(playerName string) bool {
	ok, playerFileName := s.getPlayerFileName(playerName)
	if !ok {
		return false
	}
	log.Println("Loading player %s", playerFileName)

	fileContent, fileIoErr := ioutil.ReadFile(playerFileName)
	if fileIoErr != nil {
		log.Printf("\n")
		log.Printf("File %s could not be loaded\n", playerFileName)
		log.Printf("%v", fileIoErr)
		//return fileIoErr
		return false
	}

	player := Player{}
	if xmlerr := xml.Unmarshal(fileContent, &player); xmlerr != nil {
		log.Printf("\n")
		log.Printf("File %s could not be Unmarshaled\n", playerFileName, xmlerr)
		log.Printf("%v", xmlerr)
		//return xmlerr
		return false
	}
	log.Printf(" loaded: %s", player.Gamename)
	s.addPlayer(player)

	return true
}

func (s *Server) addLevel(level Level) error {
	if level.Tag == "default" {
		log.Printf("default level loaded: %s\n", level.Key)
		s.defaultLevel = level
	}
	s.levels[level.Key] = level
	return nil
}

func (s *Server) addPlayer(player Player) error {
	s.players[player.Nickname] = player
	return nil
}

func (s *Server) GetPlayerByNick(nickname string) (Player, bool) {
	player, ok := s.players[nickname]
	return player, ok
}

func (s *Server) GetRoom(key string) (Level, bool) {
	level, ok := s.levels[key]
	return level, ok
}

func (s *Server) GetName() string {
	return s.Config.Name
}

func (s *Server) CreatePlayer(nick string, name string, playerType string) {
	ok, playerFileName := s.getPlayerFileName(nick)
	if !ok {
		return
	}
	if _, err := os.Stat(playerFileName); err == nil {
		s.LoadPlayer(nick)
		fmt.Printf("Player %s does already exists", nick)
		return
	}
	player := Player{
		Gamename:   name,
		Nickname:   nick,
		PlayerType: playerType,
		Position:   s.defaultLevel.Key,
	}
	s.addPlayer(player)
}

func (s *Server) SavePlayer(player Player) bool {
	data, err := xml.MarshalIndent(player, "", "    ")
	if err == nil {
		ok, playerFileName := s.getPlayerFileName(player.Nickname)
		if !ok {
			return false
		}

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
	client.WriteLineToUser(fmt.Sprintf("Good bye %s", client.Player.Gamename))
}

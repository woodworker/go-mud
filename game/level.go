package game

type Level struct {
	Key			string		`xml:"key,attr"`
	Name		string		`xml:"name"`
	Directions	[]Direction	`xml:"directions>direction"`
}

type Direction struct {
	Station			string		`xml:"station,attr"`
	Direction		string      `xml:",chardata"`
}

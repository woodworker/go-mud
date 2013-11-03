package game

type Level struct {
	Key			string		`xml:"key,attr"`
	Name		string		`xml:"name"`
	Directions	[]Direction	`xml:"directions>direction"`
	Intro		string		`xml:"intro"`
}

type Direction struct {
	Station			string		`xml:"station,attr"`
	Hidden			string		`xml:"hidden,attr"`
	Direction		string      `xml:",chardata"`
}

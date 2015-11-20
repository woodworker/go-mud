package game

import (
	"testing"
)

func MakeTestLevel(name string, tag string) Level {
	lvl := Level{
		Key:   name,
		Name:  "Room"+name,
		Intro: "Room"+name+"Intro",
		Tag: tag,
	}
	return lvl
}

func MakeTestDirection(name string, station string) Direction {
	return Direction{
		Station:   name,
		Direction: station,
	}
}

func TestLevelGoDirectionWithoutDependencies(t *testing.T) {
	p := Player{}

	lvl1 := MakeTestLevel("A", "default")
	dir := MakeTestDirection("North", "b")

	ok, message := lvl1.CanGoDirection(dir, p)

	if !ok {
		t.Error("Not allowed to go direction")
	}
	if message != "" {
		t.Error("Should not return message")
	}
}

func TestLevelGoDirectionWithMissingDependencies(t *testing.T) {
	p := Player{}

	lvl1 := MakeTestLevel("A", "default")

	var dependencies []Dependency
	dependencies = append(dependencies, Dependency{
		Key:"test",
		OkMessage:"OK",
		FailMessage:"FAIL",
	})

	dir := MakeTestDirection("North", "b")
	dir.Dependencies = dependencies

	ok, message := lvl1.CanGoDirection(dir, p)

	if ok {
		t.Error("Should not be allowed to go direction")
	}
	if message != "FAIL" {
		t.Error("Should get FAIL as message")
	}
}

func TestLevelGoDirectionWithDependencies(t *testing.T) {
	p := Player{}

	lvl1 := MakeTestLevel("A", "default")

	var dependencies []Dependency

	dependency := Dependency{
		Key:"test",
		OkMessage:"OK",
		FailMessage:"FAIL",
	}
	dependencies = append(dependencies, dependency)

	p.LogAction(dependency.Key)

	dir := MakeTestDirection("North", "b")
	dir.Dependencies = dependencies

	ok, message := lvl1.CanGoDirection(dir, p)

	if !ok {
		t.Error("Should not be allowed to go direction")
	}
	if message != "OK" {
		t.Error("Should get OK as message")
	}
}

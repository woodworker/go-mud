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

func TestLevelCanDoAction(t *testing.T) {
	p := Player{}
	lvl1 := MakeTestLevel("A", "default")

	action := Action{
		Name:"testaction",
	}

	var dependencies []Dependency
	dependency := Dependency{
		Key:"test",
		OkMessage:"OK",
		FailMessage:"FAIL",
	}
	dependencies = append(dependencies, dependency)
	action.Dependencies = dependencies

	ok, message := lvl1.CanDoAction(action, p)

	if ok {
		t.Error("Should not be allowed to do action")
	}
	if message != "FAIL" {
		t.Error("Should get FAIL as message")
	}

	p.LogAction("test")

	ok, message = lvl1.CanDoAction(action, p)

	if !ok {
		t.Error("Should now be allowed to do action")
	}
	if message != "OK" {
		t.Error("Should get OK as message")
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

func TestCheckDependenciesAttribute(t *testing.T) {
	p := Player{}

	dependency := Dependency{
		Key:"test",
		Type:"attribute",
		MinValue:10,
		MaxValue:15,
		OkMessage:"OK",
		FailMessage:"FAIL",
	}
	var dependencies []Dependency
	dependencies = append(dependencies, dependency)

	ok, msg := CheckDependencies(dependencies, p, "YES")
	if ok {
		t.Error("Should get not OK")
	}
	if msg != "FAIL" {
		t.Error("Should get not FAIL as message")
	}

	p.UpdateAttribute("test", 10)
	ok, msg = CheckDependencies(dependencies, p, "YES")
	if !ok {
		t.Error("Should get OK")
	}
	if msg != "YES" {
		t.Error("Should get YES as message")
	}

	ok, msg = CheckDependencies(dependencies, p, "")
	if !ok {
		t.Error("Should get OK")
	}
	if msg != "OK" {
		t.Error("Should get OK as message")
	}

	p.UpdateAttribute("test", 10)
	ok, msg = CheckDependencies(dependencies, p, "YES")
	if ok {
		t.Error("Should get no OK")
	}
	if msg != "FAIL" {
		t.Error("Should get FAIL as message")
	}
}


func TestLevelCanSeeDirection(t *testing.T) {
	p := Player{}

	lvl1 := MakeTestLevel("A", "default")

	dir := MakeTestDirection("North", "b")

	// simple check
	ok := lvl1.CanSeeDirection(dir, p, "")
	if !ok {
		t.Error("Should be allowed to see direction")
	}

	// hidden check
	dir.Hidden = true
	ok = lvl1.CanSeeDirection(dir, p, "")
	if ok {
		t.Error("Should not be allowed to see direction")
	}

	// hidden direction check
	ok = lvl1.CanSeeDirection(dir, p, "north")
	if ok {
		t.Error("Should not be allowed to see direction")
	}

	// hidden but already visited check
	p.LogAction(dir.Station)
	ok = lvl1.CanSeeDirection(dir, p, "")
	if !ok {
		t.Error("Should be allowed to see direction")
	}
}

func TestLevelGetRoomAction(t *testing.T) {

	lvl1 := MakeTestLevel("A", "default")

	// unkown action
	_, ok := lvl1.GetRoomAction("unkown")
	if ok {
		t.Error("Should get false on unkwon action")
	}

	// knwon action
	var actions []Action
	action := Action{
		Name:"testaction",
	}
	actions = append(actions, action)
	lvl1.Actions = actions

	var testAction Action
	testAction, ok = lvl1.GetRoomAction("testaction")
	if !ok {
		t.Error("Should get true on unkwon action")
	}
	if testAction.Name != "testaction" {
		t.Error("Should get false on unkwon action")
	}
}
package game

import (
	"testing"
)

func TestAttributeUpdate(t *testing.T) {
	player := Player{}

	player.UpdateAttribute("test", 10);

	if len(player.Attributes) != 1 {
		t.Error("Player should have exact one attribute")
	}

	if player.GetAttribute("test") != 10 {
		t.Error("Player attribute should be exact 10")
	}

	if player.GetAttribute("random") != 0 {
		t.Error("Not existing Player attribute should be exact 0")
	}

	player.UpdateAttribute("test", 10);

	if len(player.Attributes) != 1 {
		t.Error("Player should have exact one attribute")
	}

	if player.GetAttribute("test") != 20 {
		t.Error("Updated Player attribute should be exact 20")
	}

	if player.GetAttribute("random") != 0 {
		t.Error("Not existing Player attribute should be exact 0")
	}

	player.UpdateAttribute("test", -30);

	if len(player.Attributes) != 1 {
		t.Error("Player should have exact one attribute")
	}

	if player.GetAttribute("test") != 0 {
		t.Error("Updated Player attribute should be exact 0")
	}

	if player.GetAttribute("random") != 0 {
		t.Error("Not existing Player attribute should be exact 0")
	}
}

func TestActions(t *testing.T) {
	player := Player{}

	player.LogAction("test");

	if len(player.ActionLog) != 1 {
		t.Error("Player should have exact one action log entry")
	}

	if !player.HasAction("test") {
		t.Error("Player player should have test action")
	}

	if player.HasAction("foo") {
		t.Error("Player player should not have foo action")
	}
}
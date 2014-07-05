// Package player defines the Player interface.
package player

import (
    "dr2w.com/hf/model/state"
    "dr2w.com/hf/model/action"
)

// Player defines the basic interface needed for an entity (human or AI) to play the game.
type Player interface {
    Play(state state.State, message action.Message) []int
}

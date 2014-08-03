// Package ai describes a framework for building AI Players from components.
package ai

import (
    "dr2w.com/hf/model/action"
    "dr2w.com/hf/model/state"
)

// Decider is a function which takes a State and Specific Type of Message and returns
// the selected options from the Message.
type Decider func(s state.State, m action.Message) []int

// AIPlayer describes a set of behaviors for a given AI Player.
type AIPlayer struct {
    Name string
    Deciders map[action.Type]Decider
}

func (p AIPlayer) String() string {
    return p.Name
}

// Play implements the player.Player interface.
func (p AIPlayer) Play(s state.State, m action.Message) []int {
    if decider, ok := p.Deciders[m.Type]; ok {
        return decider(s,m) 
    }
    return nil
}

// Update implements the player.Player interface.
func (p AIPlayer) Update(s state.State, t action.Type) {
    return
}

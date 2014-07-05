// Package action defines player inputs into the game.
package action

import (
	"fmt"

	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/state"
)

// Type specifies the type of action to be performed next on a State.
type Type int

const (
	Deal Type = iota
	Bid
	Trump
	Discard
    ReDeal
	Play
	ThrowIn
    Score
)

// typeNames map types to their string names
var typeNames = map[Type]string{
	Deal:    "Deal",
	Bid:     "Bid",
	Trump:   "Trump",
    ReDeal:  "ReDeal",
	Discard: "Discard",
	Play:    "Play",
	ThrowIn: "ThrowIn",
    Score:   "Score",
}

// String returns the string representation of a Type.
func (t Type) String() string {
	return typeNames[t]
}

// SelectionRange returns a slice of Selections including all between start and
// end, not including end.
func SelectionRange(start, end int) (s []int) {
	for i := start; i < end; i++ {
		s = append(s, i)
	}
	return s
}

// Message represents a Message or Message to/from the player for a choice of next play.
// Options represent valid options when sent to the player, and a subset of these options
// should be returned by the player to indicate selection.
type Message struct {
	Type    Type
	Seat    seat.Seat
	Options []int
}

// String returns a human readable representation of the Message.
func (m Message) String() string {
	return fmt.Sprintf("Message - %s(%s) %v", m.Type, m.Seat, m.Options)
}

// Selection returns a single value from the Options array, or error if more
// than one value or no values exist.
func (m Message) Selection() (int, error) {
	if len(m.Options) != 1 {
		return 0, fmt.Errorf("expected single option in message (%s).", m)
	}
	return m.Options[0], nil
}

// Func is a function that performs a specific action on a State.
// It then returns the new resulting State and next Action.
type Func func(s state.State, m Message) (state.State, Message, error)

// actionMap is a static mapping of Types to their corresponding Funcs.
var actionMap = map[Type]Func{
	Deal:  deal,
	Bid:   bids,
	Trump: trump,
    ReDeal: redeal,
    Discard: discard,
    Play: play,
    Score: score,
}

// NextState converts a State and an Action into the next State and Action.
func NextState(s state.State, m Message) (state.State, Message, error) {
	return actionMap[m.Type](s, m)
}

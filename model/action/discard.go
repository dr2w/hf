package action

import (
	"fmt"

	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/state"
)

// omit takes a hand of cards and a list of card indices to omit. It returns
// a copy of the original Hand with the selected cards omitted.
func omit(h hand.Hand, indices []int) *hand.Hand {
	m := make(map[int]bool)
	for _, i := range indices {
		m[i] = true
	}
	var newHand hand.Hand
	for i, c := range h {
		if _, ok := m[i]; !ok {
			newHand = append(newHand, c)
		}
	}
	return &newHand
}

// discard takes the results of a discard action by the bid winner and
// modifies the game state to a valid state for starting play.
func discard(s state.State, m Message) (state.State, Message, error) {
	hand := *s.Hands[m.Seat]
	toDiscard := len(hand) - state.HandSize
	if toDiscard != len(m.Options) {
		return state.State{}, Message{}, fmt.Errorf("Invalid discard selection. Selected %d cards instead of %d", len(m.Options), toDiscard)
	}
	s.Hands[m.Seat] = omit(hand, m.Options)
	return s, Message{
		Type:    Play,
		Seat:    m.Seat,
		Options: SelectionRange(0, state.HandSize),
	}, nil
}

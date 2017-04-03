package action

import (
	"fmt"

	"dr2w.com/hf/model/card"
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

// discard takes the results of a discard action and removes the selected cards from
// the given player's hand. It then either calls for another discard action or triggers
// the redeal if all non-winners have discarded. It assumes that the first discard call
// is made to winner.Next().
// If called with the winner's seat, discard assumes that the redeal has occurred, applies
// the winner's discard, then modifies the game state to a valid state for starting play.
func discard(s state.State, m Message) (state.State, Message, error) {
	winner, _ := s.WinningBid()
	hand := *s.Hands[m.Seat]
	newHand := omit(hand, m.Options)
	trump := card.Set(hand).TrumpCards(s.Trump)
	if m.Seat != winner && newHand.Length() != trump.Length() {
		return state.State{}, Message{}, fmt.Errorf("Invalid discard selection. Must select all non-trump. Hand had %d remaining cards, expected %d", newHand.Length(), trump.Length())
	}
	if m.Seat == winner && newHand.ExtraCards() != 0 {
		return state.State{}, Message{}, fmt.Errorf("Invalid discard selection. Should have no extra cards reminaing, but found %d", hand.ExtraCards())
	}
	s.Hands[m.Seat] = newHand

	if m.Seat.Next() == winner {
		// Redeal
		return s, Message{
			Type:    ReDeal,
			Seat:    s.Dealer,
			Options: []int{0},
			Expect: 1,
		}, nil
	}
	if m.Seat == winner {
		if len(s.Deck) > 0 {
			return state.State{}, Message{}, fmt.Errorf("Invalid state reached. Discard called for winner with non-empty deck: %v", s.Deck) 
		}
		// Start Play.
		return s, Message{
			Type:    Play,
			Seat:    winner,
			Options: SelectionRange(0, s.Hands[winner].Length()),
			Expect: 1,
		}, nil
	}
	// Continue discarding
	nextToDiscard := m.Seat.Next()
	return s, Message{
			Type:    Discard,
			Seat:    nextToDiscard,
			Options: s.Hands[nextToDiscard].Discards(s.Trump),
			Expect: s.Hands[nextToDiscard].NumToDiscard(s.Trump),
	}, nil
}


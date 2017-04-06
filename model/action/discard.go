package action

import (
	"fmt"

	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/seat"
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
	newHand := omit(*s.Hands[m.Seat], m.Options)
	if err := validateNewHand(s, m, newHand); err != nil {
		return state.State{}, Message{}, err
	}
	s.Hands[m.Seat] = newHand
	return s, nextMessage(s, m.Seat), nil
}

func validateNewHand(s state.State, m Message, h *hand.Hand) error {
	winner, _ := s.WinningBid()
	trump := card.Set(*s.Hands[m.Seat]).TrumpCards(s.Trump)
	if m.Seat != winner && h.Length() != trump.Length() {
		return fmt.Errorf("Invalid discard selection. Must select all non-trump. Hand had %d remaining cards, expected %d", h.Length(), trump.Length())
	}
	if m.Seat == winner && h.ExtraCards() != 0 {
		return fmt.Errorf("Invalid discard selection. Should have no extra cards reminaing, but found %d", h.ExtraCards())
	}
	if m.Seat == winner && len(s.Deck) > 0 {
		return fmt.Errorf("Invalid state reached. Discard called for winner with non-empty deck: %v", s.Deck) 
	}
	return nil
}

func nextMessage(s state.State, st seat.Seat) Message {
	winner, _ := s.WinningBid()
	if st == winner {
		return Message{
			Type:    Play,
			Seat:    winner,
			Options: SelectionRange(0, s.Hands[winner].Length()),
			Expect: 1,
		}
	}
	nextToDiscard := st.Next()
	if nextToDiscard == winner {
		return Message{
			Type:    ReDeal,
			Seat:    s.Dealer,
			Options: []int{0},
			Expect: 1,
		}
	}
	return Message{
			Type:    Discard,
			Seat:    nextToDiscard,
			Options: s.Hands[nextToDiscard].Discards(s.Trump),
			Expect: s.Hands[nextToDiscard].NumToDiscard(s.Trump),
	}
}


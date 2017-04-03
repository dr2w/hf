package action

import (
	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/state"
)

// trump takes a State and a selected trump and returns the next State and Action.
func trump(s state.State, m Message) (state.State, Message, error) {
	sel, err := m.Selection()
	if err != nil {
		return state.State{}, Message{}, err
	}

	// Set Trump and adjust cards to match trump suit.
	s.Trump = card.Suits[sel]
	joker := s.FindCard(card.Card{card.Joker, card.NoSuit})
	joker.Suit = s.Trump
	five := s.FindCard(card.Card{card.Five, card.SameColorSuit(s.Trump)})
	five.Suit = s.Trump
	five.Value = card.OffFive

    	winner, _ := s.WinningBid()
	firstToDiscard := winner.Next()
	return s, Message{
		Type:    Discard,
		Seat:    firstToDiscard,
		Options: s.Hands[firstToDiscard].Discards(s.Trump),
		Expect: s.Hands[firstToDiscard].NumToDiscard(s.Trump),
	}, nil
}

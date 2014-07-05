package action

import (
	"fmt"

	"dr2w.com/hf/model/bid"
	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/state"
)

const (
	// cardsPerDeal is the number of cards drawn from the deck at once during a deal.
	cardsPerDeal = 3

	// dealsPerHand is the number of rounds of dealing to each Hand.
	dealsPerHand = 3
)

// deal takes a State with a full Deck and returns a State with a full Hand dealt
// to each Seat. The Selection is ignored.
func deal(s state.State, _ Message) (state.State, Message, error) {
	return dealWithPattern(s, cardsPerDeal, dealsPerHand)
}

// dealWithPattern deals cards from the Deck to the Hands according to the specified
// pattern.
func dealWithPattern(s state.State, cpd int, dph int) (state.State, Message, error) {
	s.Hands = make(map[seat.Seat]*hand.Hand)
	s.Dealer = s.Dealer.Next()
	for i := 0; i < dph; i++ {
		st := s.Dealer.Next()
		for j := 0; j < len(seat.Order); j++ {
			if _, ok := s.Hands[st]; !ok {
				s.Hands[st] = &hand.Hand{}
			}
			cards, err := s.Deck.Deal(cpd)
			if err != nil {
				return state.State{}, Message{}, fmt.Errorf("Unable to deal %d cards from %v: %s", cpd, s.Deck, err)
			}
			s.Hands[st].Add(cards...)
			st = st.Next()
		}
	}
	r := Message{
		Type:    Bid,
		Seat:    s.Dealer.Next(),
		Options: SelectionRange(0, len(bid.Values)),
	}
	return s, r, nil
}

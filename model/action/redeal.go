
package action

import (
    "fmt"
    "math/rand"

	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/state"
)

// redeal takes a State with bidding complete, trump decided, and non-winners
// discarded down to only trump. It deals the three
// non-bid-winning hands up to handSize and deals the rest of the deck to the
// bid-winner.
func redeal(s state.State, _ Message) (state.State, Message, error) {
    winner, _ := s.WinningBid()
    st := s.Dealer.Next()
    for i := 0; i < len(seat.Order); i++ {
        if st == winner {
	    st = st.Next()
            continue
        }
	// TODO(drw): Fix now that we discard non-trump before this.
	toDeal := -s.Hands[st].ExtraCards()
	fmt.Printf("ToDeal: %d", toDeal)
        if toDeal > 0 {
            cards, err := s.Deck.Deal(toDeal)
            if err != nil {
                return state.State{}, Message{}, err
            }
            s.Hands[st].Add(cards...)
        } else if toDeal < 0 {
	    for s.Hands[st].ExtraCards() > 0 {
            	d := s.Hands[st].Discards(s.Trump)
            	if len(d) < s.Hands[st].ExtraCards() {
                	return state.State{}, Message{}, fmt.Errorf("unable to discard/redeal from hand %v with trump %v", s.Hands[st], s.Trump)
                }
                // TODO(drw): make this selection more reasonable
                // TODO(drw): add s.Reveal(card) for this and for bid winner reveal
                s.Hands[st].Remove(d[rand.Intn(len(d))])
	    }
        }
        st = st.Next()
    }
    cards, err := s.Deck.Deal(len(s.Deck))
    if err != nil {
        return state.State{}, Message{}, err
    }
    s.Hands[winner].Add(cards...)
    // We enforce sorting on players' hands everywhere we add to them.
    for st := range s.Hands {
	card.Set(*s.Hands[st]).Sort()
    }
	r := Message{
		Type:    Discard,
		Seat:    winner,
		Options: s.Hands[winner].Discards(s.Trump),
		Expect: s.Hands[winner].ExtraCards(),
	}
	return s, r, nil
}

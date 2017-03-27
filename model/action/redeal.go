
package action

import (
    "fmt"

	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/state"
)

// handSize is the final hand size after redealing.
const handSize = 6

// discards takes in a trump suit and a Hand, and returns the indices of
// the cards which can be discarded.
func discards(trump card.Suit, h *hand.Hand) (options []int) {
	pointCards := (card.Set(*h)).PointCards(trump)
	trumpCards := (card.Set(*h)).TrumpCards(trump)
	for i, c := range *h {
		if c.Suit != trump ||
			len(trumpCards) > state.HandSize && c.Points(trump) == 0 ||
			len(pointCards) > state.HandSize && c.Value == card.Deuce && c.Suit == trump {
			options = append(options, i)
		}
	}
	fmt.Printf("\n\n%v\n%v\n\n", card.Set(*h), options)
	return options
}

// redeal takes a State with bidding complete and trump decided. It deals the three
// non-bid-winning hands up to handSize and deals the rest of the deck to the
// bid-winner.
func redeal(s state.State, _ Message) (state.State, Message, error) {
    winner, _ := s.WinningBid()
    st := s.Dealer.Next()
    for i := 0; i < len(seat.Order); i++ {
	fmt.Printf("Seat: %s\n", st.String())
	fmt.Printf("Winner: %s\n", winner.String())
        if st == winner {
	    st = st.Next()
            continue
        }
        newHand := &hand.Hand{}
        for _, c := range *s.Hands[st] {
            if c.Suit == s.Trump {
                newHand.Add(c)
            }
        }
        toDeal := handSize - newHand.Length()
	fmt.Printf("\nNew Hand: %s\n", newHand.String())
	fmt.Printf("ToDeal: %d", toDeal)
        if toDeal > 0 {
            cards, err := s.Deck.Deal(toDeal)
            if err != nil {
                return state.State{}, Message{}, err
            }
            newHand.Add(cards...)
        } else if toDeal < 0 {
            d := discards(s.Trump, newHand)
            if newHand.Length() - len(d) > handSize {
                return state.State{}, Message{}, fmt.Errorf("unable to discard/redeal from hand %v with trump %v", newHand, s.Trump)
            }
            // TODO(drw): make this selection more reasonable
            // TODO(drw): add s.Reveal(card) for this and for bid winner reveals
            // TODO(drw): this is broken. I'm precomputing the indices to remove
            // then removing them one at a time (so every one after the first is
            // wrong)
            for i := 0; newHand.Length() > handSize; i++ {
                *newHand = append((*newHand)[:d[i]], (*newHand)[d[i]+1:]...)
            }
        }
        s.Hands[st] = newHand
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
		Options: discards(s.Trump, s.Hands[winner]),
	}
	return s, r, nil
}

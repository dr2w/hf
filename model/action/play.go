package action

import (
	"fmt"

	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/state"
	"dr2w.com/hf/model/trick"
)

// playCard takes a State and plays the given card on the last (or next) trick,
// returning the resulting slice of played tricks. Returns an error if this
// appears to be an invalid play.
func playCard(s state.State, st seat.Seat, c card.Card) ([]trick.Trick, error) {
	last := s.LastPlayed()
	if last.Full() || last.Empty() {
		return append(s.Played, trick.Trick{Cards: map[seat.Seat]card.Card{st: c}, First: st}), nil
	}
	if _, ok := last.Cards[st]; ok {
		return []trick.Trick{}, fmt.Errorf("attempting to play %v:%v on Trick %v in State:\n%v", st, c, last, s)
	}
	s.Played[len(s.Played)-1].Cards[st] = c
	return s.Played, nil
}

// validCards takes a Hand, Trump, and Last Trick (complete or incomplete) and returns
// an integer slice of the valid plays from the Hand.
func validCards(h *hand.Hand, trump card.Suit, trick trick.Trick) (options []int) {
	for i, c := range *h {
		if c.Suit == trump ||
			c.Suit == trick.SuitLead() ||
			trick.Empty() ||
			trick.Full() ||
			!h.HasSuit(trick.SuitLead()) {
			options = append(options, i)
		}
	}
	return options
}

// next takes the current State of the game, waiting for a card to be played, and
// returns the Message requesting that card or a Message indicating end of hand.
func next(s state.State) (Message, error) {
	st, h := s.ToPlay()
	if h == nil {
		//return Message{}, fmt.Errorf("next unable to locate hand for %v", st)
        return Message{
            Type:  Score,
            Seat:  seat.None,
            Options: []int{0},
	    Expect: 1,
        }, nil
	}
	return Message{
		Type:    Play,
		Seat:    st,
		Options: validCards(h, s.Trump, s.LastPlayed()),
		Expect: 1,
	}, nil
}

// advanceEmptyHands advances the given state through any empty Hands until it
// finds the next non-empty Hand. Returns the final state with all empty hands
// advanced.
func advanceEmptyHands(s state.State) (state.State, error) {
	st, h := s.ToPlay()
	if h != nil && len(*h) == 0 {
		played, err := playCard(s, st, card.Card{})
		if err != nil {
			return state.State{}, err
		}
		s.Played = played
		ns, err := advanceEmptyHands(s)
		if err != nil {
			return state.State{}, err
		}
		return ns, nil
	}
	return s, nil
}

// play takes a State and a selected Card to play and returns the resulting
// state and a request for the next play (doesn't detect end-of-game).
func play(s state.State, m Message) (state.State, Message, error) {
    hand, ok := s.Hands[m.Seat]
    if !ok {
        return state.State{}, Message{}, fmt.Errorf("Seat %v has no Hand (%v)", m.Seat, s)
    }
	index, err := m.Selection()
	if err != nil {
		return state.State{}, Message{}, err
	}
	card, err := hand.Remove(index)
	if err != nil {
		return state.State{}, Message{}, err
	}
	s.Played, err = playCard(s, m.Seat, card)
	if err != nil {
		return state.State{}, Message{}, err
	}
	s, err = advanceEmptyHands(s)
	if err != nil {
		return state.State{}, Message{}, err
	}
	msg, err := next(s)
	if err != nil {
		return state.State{}, Message{}, err
	}
	return s, msg, nil
}

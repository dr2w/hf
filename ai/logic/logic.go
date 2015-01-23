// Package logic wraps state.State and implements many high level
// logic/state related functions on top of the state.
package logic

import "dr2w.com/hf/model/card"
import "dr2w.com/hf/model/seat"
import "dr2w.com/hf/model/state"

// Logix wraps a state.State from a given Perspective and provides
// several high level logical computations.
type Logic struct {
	State       state.State
	Perspective seat.Seat
}

func (l Logic) played() card.Set {
	s := card.Set{}
	for _, t := range l.State.Played {
		for _, c := range t.AsCardSet() {
			s = append(s, c)
		}
	}
	return s.AsTrump(l.State.Trump)
}

func (l Logic) TrumpPlayed() card.Set {
	return l.played().TrumpCards(l.State.Trump)
}

func (l Logic) TrumpOut() card.Set {
	out := card.Set{}
	played := l.TrumpPlayed()
	for _, v := range card.Values {
		c := card.Card{v, l.State.Trump}
		if !played.Contains(c) {
			out = append(out, c)
		}
	}
	return out
}

func (l Logic) MyTrump() card.Set {
	myHand := *l.State.Hands[l.Perspective]
	return card.Set(myHand).AsTrump(l.State.Trump).TrumpCards(l.State.Trump)
}

func (l Logic) TopTrumpOut() card.Card {
	return max(l.TrumpOut())
}

func (l Logic) MyTopTrump() card.Card {
	return max(l.MyTrump())
}

func max(s card.Set) card.Card {
	maxCard := card.Card{}
	for _, c := range s {
		if c.Value > maxCard.Value {
			maxCard = c
		}
	}
	return maxCard
}

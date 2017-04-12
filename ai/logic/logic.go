// Package logic wraps state.State and implements many high level
// logic/state related functions on top of the state.
package logic

import (
	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/state"
)

// Logic wraps a state.State from a given Perspective and provides
// several high level logical computations.
type Logic struct {
	State       state.State
	Perspective seat.Seat
}

// played returns all cards played on *previous* tricks.
func (l Logic) played() card.Set {
	s := card.Set{}
	for _, t := range l.State.Played {
		if !t.Full() {
			continue
		}
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
	return l.MyHand().AsTrump(l.State.Trump).TrumpCards(l.State.Trump)
}

func (l Logic) TopTrumpOut() card.Card {
	return max(l.TrumpOut())
}

func (l Logic) MyTopTrump() card.Card {
	return max(l.MyTrump())
}

func (l Logic) IHaveHighCardOut() bool {
	return l.IHave(l.TopTrumpOut())
}

func (l Logic) IAmLeading() bool {
	return l.State.LastPlayed().Empty()
}

func (l Logic) OffsuitLead() bool {
	return false
}

func (l Logic) MyHand() card.Set {
	return card.Set(*l.State.Hands[l.Perspective])
}

func (l Logic) IHave(c card.Card) bool {
	return l.MyHand().Contains(c)
}

func (l Logic) IHaveAFive() bool {
	return l.IHave(card.Card{Suit: l.State.Trump, Value: card.Five}) ||
		l.IHave(card.Card{Suit: l.State.Trump, Value: card.OffFive})
}

func (l Logic) ICanCoverAFive() bool {
    for _, c := range l.MyHand() {
        if c.TrumpValue(l.State.Trump) > card.Five {
            return true
        }
    }
    return false
}

func (l Logic) AFiveIsOut() bool {
    return l.TrumpOut().Contains(card.Card{card.Five, l.State.Trump}) ||
            l.TrumpOut().Contains(card.Card{card.OffFive, l.State.Trump})
}

func (l Logic) IAmLast() bool {
	// TODO(drw): Add support for players throwing in.
	return len(l.State.LastPlayed().Cards) == 3
}

func (l Logic) IHaveHighCard() bool {
    for _, c := range l.MyHand() {
        if l.State.LastPlayed().WouldTakeLead(c, l.State.Trump) {
            return true
        }
    }
    return false
}

func (l Logic) PointsAreShowing() bool {
    points := 0
    for _, c := range l.State.LastPlayed().AsCardSet() {
        points += c.Points(l.State.Trump)
    }
    return points > 0
}

func (l Logic) NextPlayerIsLast() bool {
	// TODO(drw): Add support for players throwing in.
    return len(l.State.LastPlayed().Cards) == 2
}

func (l Logic) TrickHasAFive() bool {
    return l.State.LastPlayed().AsCardSet().Contains(card.Card{card.Five, l.State.Trump}) ||
            l.State.LastPlayed().AsCardSet().Contains(card.Card{card.OffFive, l.State.Trump})
}

func (l Logic) PartnerPlayedHighCard() bool {
    seat, _ := l.State.LastPlayed().Winner(l.State.Trump)
    return seat == l.Perspective.Partner()
}

func (l Logic) PartnerPlayedHighCardOut() bool {
    c, ok := l.State.LastPlayed().Cards[l.Perspective.Partner()]
    return ok && c == l.TopTrumpOut() 
}

func (l Logic) PartnerToPlay() bool {
    _, ok := l.State.LastPlayed().Cards[l.Perspective.Partner()]
    return !ok
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

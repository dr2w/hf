package bidding

import (
	"dr2w.com/hf/model/action"
	"dr2w.com/hf/model/bid"
	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/state"
)

type forSuit func(cards card.Set) (min, max bid.Bid)
type isSix func(cards card.Set) bool

type bidder func(s state.State, m action.Message) ([]int, card.Suit)
type decider func(s state.State, m action.Message) []int

var (
	DRWValue = value(fromBidders(drwSuitBid, drwSixBid))
	DRWSuit  = suit(fromBidders(drwSuitBid, drwSixBid))
)

// Value converts a bidder into a decider that returns the
// value of the bid (as opposed to the suit)
func value(f bidder) decider {
	return func(s state.State, m action.Message) []int {
		v, _ := f(s, m)
		return v
	}
}

// Suit convers a bidder into a decider that returns the
// suit of the bid (as opposed to the value)
func suit(f bidder) decider {
	return func(s state.State, m action.Message) []int {
		_, suit := f(s, m)
		for i, st := range card.Suits {
			if st == suit {
				return []int{i}
			}
		}
		// else likely no suit, return 0 for now.
		return []int{0}
	}
}

// fromBidders takes in a forSuit and an isSix
// and returns a Decider that will bid based on these two
// functions. If no max suit bid is >= 8, then it first
// checks for a valid six bid. If there is not a valid
// six bid, then it bids according to the min and max
// of the suit with the highest max bid.
// TODO(drw): consider adding in aggressiveness and
// randomness as adjustments to min and max.
// TODO(drw): deal with the no mans land between 10 and 14/28.
func fromBidders(forSuit forSuit, isSix isSix) func(s state.State, m action.Message) ([]int, card.Suit) {
	return func(s state.State, m action.Message) ([]int, card.Suit) {
		h := *s.Hands[m.Seat]
		minSuit, maxSuit := bid.Pass, bid.Pass
		bestSuit := card.NoSuit
		for _, suit := range card.Suits {
			suitCards := card.Set(h).AsTrump(suit).TrumpCards(suit)
			min, max := forSuit(suitCards)
			if max > maxSuit {
				minSuit = min
				maxSuit = max
				bestSuit = suit
			}
		}
		_, currentBid := s.WinningBid()
		if maxSuit < bid.B8 && currentBid == bid.Pass &&
			isSix(card.Set(h)) {
			// TODO(drw): handle the odd case where this wins.
			// For now we just pick an arbitrary suit
			return []int{int(bid.B6)}, card.NoSuit
		}
		for b := minSuit; b <= maxSuit; b++ {
			if b > currentBid {
				return []int{int(b)}, bestSuit
			}
		}
		return []int{int(bid.Pass)}, card.NoSuit
	}
}

// drwSixBid implements a basic version of the logic drw uses
// as an isSix function.
func drwSixBid(cards card.Set) bool {
	suitCovered := func(s card.Suit) bool {
		if cards.Contains(card.Card{card.Five, s}) ||
			cards.Contains(card.Card{card.Five, card.SameColorSuit(s)}) {
			return true
		}
		for _, c := range cards {
			if c.Suit == s && c.Value > card.Joker {
				return true
			}
			if c.Value == card.Joker &&
				len(cards.AsTrump(s).TrumpCards(s)) > 2 {
				return true
			}
		}
		return false
	}
	for _, suit := range card.Suits {
		if !suitCovered(suit) {
			return false
		}
	}
	return true
}

// drwSuitBid implements a basic version of the per-suit
// bidding logic drw uses as a forSuit function:
func drwSuitBid(cards card.Set) (min, max bid.Bid) {

	// hasN is a helper function that pattern matches
	// the hand to help choose a bid.
	hasN := func(n int, values string) bool {
		count := 0
		for _, v := range card.ValuesFromShorthand(values) {
			for _, c := range cards {
				if c.Value == v {
					count++
				}
			}
		}
		return count >= n
	}
	has := func(v string) bool { return hasN(1, v) }

	switch {
	case hasN(5, "AKQJj2"):
		return bid.B1530, bid.B1530
	case hasN(4, "AKQ2"):
		return bid.B1428, bid.B1530
	case hasN(3, "AKQ"):
		return bid.B1428, bid.B1428
	case hasN(2, "AK") &&
		hasN(2, "JjT"):
		return bid.B10, bid.B1428
	case hasN(2, "AK"):
		return bid.B9, bid.B10
	case has("A") &&
		hasN(2, "KQJj"):
		return bid.B8, bid.B10
	case has("A") &&
		hasN(1, "KQJjT"):
		return bid.B8, bid.B9
	case has("A"):
		return bid.B8, bid.B8
	case hasN(1, "AKQJ") &&
		hasN(1, "5F") &&
		len(cards) > 4:
		return bid.B8, bid.B8
	}
	return bid.Pass, bid.Pass
}

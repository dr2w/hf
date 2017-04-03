package ai

import (
	"math/rand"

	"dr2w.com/hf/model/action"
	"dr2w.com/hf/model/state"
)

// simpleDiscard chooses random non-trump cards and the lowest valued
// non-point trump cards.
func simpleDiscard(s state.State, m action.Message) (discards []int) {
	hand := s.Hands[m.Seat]
	winner, _ := s.WinningBid()
	if m.Seat != winner {
		return hand.Discards(s.Trump)
	}
	n := hand.ExtraCards()
	d := hand.Discards(s.Trump)
	if len(m.Options) < n || len(d) < n {
		return []int{} // Triggers downstream error
	}
	discards = make([]int, n)
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		discards[i] = d[perm[i]]
	}
	return discards
}

package ai

import (
    "dr2w.com/hf/model/action"
    "dr2w.com/hf/model/state"
)

// simpleDiscard chooses random non-trump cards and the lowest valued
// non-point trump cards.
func simpleDiscard(s state.State, m action.Message) (discards []int) {
    hand := s.Hands[m.Seat]
    n := hand.Length() - 6 // TODO(drw): factor this (and others) out into a constants file.
    if len(m.Options) < n {
        return []int{} // Triggers downstream error
    }
    var (
        onlyTrump bool
        index int
    )
    for len(discards) < n {
        if !onlyTrump && (*hand)[index].Suit != s.Trump ||
            onlyTrump && (*hand)[index].Suit == s.Trump {
            discards = append(discards, index)
        }
        index++
        if index >= hand.Length() {
            onlyTrump = true
            index = 0
        }
    }
    return discards
}

// Dumb chooses randomly or always chooses the same option.
var Dumb = AIPlayer {
    Name: "Dumb",
    Deciders: map[action.Type]Decider{
        action.Deal: first,
        action.Bid: second,
        action.Trump: rand1,
        action.Discard: simpleDiscard,
        action.Play: first,
    },
}

package ai

import (
	"dr2w.com/hf/ai/bidding"
	"dr2w.com/hf/ai/playing"
	"dr2w.com/hf/model/action"
)

// Dumb chooses randomly or always chooses the same option.
var Dumb = AIPlayer{
	Name: "Dumb",
	Deciders: map[action.Type]Decider{
		action.Deal:    first,
		action.Bid:     second,
		action.Trump:   rand1,
		action.Discard: simpleDiscard,
		action.Play:    first,
	},
}

// DRW plays using a set of rudimentary heuristics.
var DRW = AIPlayer{
	Name: "DRW",
	Deciders: map[action.Type]Decider{
		action.Deal:    first,
		action.Bid:     Decider(bidding.DRWValue),
		action.Trump:   Decider(bidding.DRWSuit),
		action.Discard: simpleDiscard,
		action.Play:    Decider(playing.InconsistentPlayer),
	},
}

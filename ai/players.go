package ai

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
		action.Bid:     bidding.DRWValue,
		action.Trump:   bidding.DRWSuit,
		action.Discard: simpleDiscard,
		action.Play:    playing.DRWPlayer,
	},
}

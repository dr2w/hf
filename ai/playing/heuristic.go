// Package playing implements intelligent players of hands
// (not bidding). Given a bit, trump suit, hand, etc. this
// package decides what card to play next.
package playing

import "log"
import "math/rand"
import "sort"

import "dr2w.com/hf/ai/logic"
import "dr2w.com/hf/model/action"
import "dr2w.com/hf/model/card"
import "dr2w.com/hf/model/state"

var (
	InconsistentPlayer = inconsistently(0.2, scorerFromDT(basicTree))
	NoisyPlayer        = noisily(0.2, scorerFromDT(basicTree))
)

var basicTree = &tree{}

type decider func(s state.State, m action.Message) []int

type scorer func(s state.State, m action.Message, c card.Card) float64

type jointSort struct {
	sorter []float64
	sortee []int
}

func (js jointSort) Len() int           { return len(js.sorter) }
func (js jointSort) Less(i, j int) bool { return js.sorter[i] < js.sorter[j] }
func (js jointSort) Swap(i, j int) {
	js.sorter[i], js.sorter[j] = js.sorter[j], js.sorter[i]
	js.sortee[i], js.sortee[j] = js.sortee[j], js.sortee[i]
}

// inconsistently returns a decider that plays with some level
// of inconsistency. It takes in a rate of inconsistency
// (between 0 and 1) and a score to score each available
// play and returns a decider that will play well (1-'rate') of the
// time, and play suboptimally 'rate' of the time.
func inconsistently(rate float64, score scorer) decider {
	if rate < 0 || rate >= 0.99 {
		log.Printf("ERROR: played.inconsistently received a bad value for rate (%.4f). Defaulting to 0.0", rate)
		rate = 0.0
	}
	return func(s state.State, m action.Message) []int {
		var scores []float64
		for _, option := range m.Options {
			c := (*s.Hands[m.Seat])[option]
			scores = append(scores, score(s, m, c))
		}
		sort.Sort(sort.Reverse(jointSort{scores, m.Options}))
		i := 0
		for rand.Float64() < rate {
			i = (i + 1) % len(m.Options)
		}
		return []int{m.Options[i]}
	}
}

// noisily returns a decider that plays with a bit of noise.
// This decider is unlikely to make a terrible decision, but
// will play suboptimally often, as the scores are observed
// only after noise is injected. Rate determines the level
// of noise added to each score (0 - rate).
func noisily(rate float64, score scorer) decider {
	if rate < 0.01 || rate > 1.0 {
		log.Printf("ERROR: played.noisily received a bad value for rate (%.4f). Defaulting to 0.0", rate)
		rate = 0.0
	}
	return func(s state.State, m action.Message) []int {
		var scores []float64
		for _, option := range m.Options {
			c := (*s.Hands[m.Seat])[option]
			adjustment := rand.Float64()*rate*2 - 1 // [-1,-1] -> [-1,1]
			scores = append(scores, score(s, m, c)+adjustment)
		}
		sort.Sort(sort.Reverse(jointSort{scores, m.Options}))
		return []int{m.Options[0]}
	}
}

// tree provided structure for a basic decision tree.
type tree struct {
	left   *tree
	right  *tree
	goLeft func(l logic.Logic, c card.Card) bool
	score  float64
}

// evaluate traverses the decision tree using the given logic and card,
// returning the score of the leaf we end up at.
func (t *tree) evaluate(l logic.Logic, c card.Card) float64 {
	if t.left == nil && t.right == nil {
		return t.score
	}
	if t.right == nil || t.goLeft(l, c) {
		return t.left.evaluate(l, c)
	}
	return t.right.evaluate(l, c)
}

// scorerFromDT builds and returns a scorer from the given decision tree.
func scorerFromDT(t *tree) scorer {
	return func(s state.State, m action.Message, c card.Card) float64 {
		return t.evaluate(logic.Logic{s, m.Seat}, c)
	}
}

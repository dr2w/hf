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
	InconsistentPlayer = inconsistently(0.2, scorerFromDT(initialTree))
	NoisyPlayer        = noisily(0.2, scorerFromDT(initialTree))
)

type scoreFn func(c card.Card, t card.Suit) float64

const scoreMultiplier = 10.0

func combine(more, less float64) float64 {
	return (more*scoreMultiplier + less) / scoreMultiplier
}

func byValue(c card.Card, t card.Suit) float64 {
	return float64(c.TrumpValue(t)) / float64(card.MaxTrumpValue)
}

func byNegValue(c card.Card, t card.Suit) float64 {
	return 1.0 - byValue(c, t)
}

func byPoints(c card.Card, t card.Suit) float64 {
	return float64(c.Points(t)) / float64(card.MaxPoints)
}

func byNegPoints(c card.Card, t card.Suit) float64 {
    return 1.0 - byPoints(c, t)
}

func byNegPointsThenNegValue(c card.Card, t card.Suit) float64 {
	return combine(byNegPoints(c,t), byNegValue(c,t))
}

func forValuesAbove(v card.Value, f scoreFn) scoreFn {
    return func(c card.Card, t card.Suit) float64 {
        if c.TrumpValue(t) <= v {
            return 0.0
        }
        return f(c, t)
    }
}

func byFivesThenNegValue(c card.Card, t card.Suit) float64 {
	five := 0.0
	if c.Suit == t && (c.Value == card.Five || c.Value == card.OffFive) {
		five = 1.0
	}
	return combine(five, byNegValue(c, t))
}

var byValueTree = &tree{
	logic: "ScoreByValue",
	score: byValue,
}

var byNegValueTree = &tree{
	logic: "ScoreByNegValue",
	score: byNegValue,
}
var byFivesThenNegValueTree = &tree{
	logic: "ScoreByFivesThenNegValue",
	score: byFivesThenNegValue,
}

var byNegPointsThenNegValueTree = &tree{
	logic: "ScoreByNegPoints",
	score: byNegPointsThenNegValue,
}

// Basic Logic:
var initialTree = &tree{
	logic: "If I'm leading",
	goLeft: logic.Logic.IAmLeading,
	left: leadingTree,
	right: notLeadingTree,
}

var leadingTree = &tree{
	logic: "If I have the high card",
	goLeft: logic.Logic.IHaveHighCardOut,
	left: byValueTree,
	right: byNegPointsThenNegValueTree,
}

var notLeadingTree = &tree{
	logic: "If Offsuit",
	goLeft: logic.Logic.OffsuitLead,
	left: offsuitLeadTree,
	right: trumpLeadTree,
}

var offsuitLeadTree = &tree{
	logic: "If I have a 5 and I am last",
	goLeft: func(l logic.Logic) bool {
		return l.IHaveAFive() && l.IAmLast()
	},
	left: byFivesThenNegValueTree,
	right: offsuitLeadNoFiveTree,
}

var trumpLeadTree = &tree{
	logic: "If there's a 5",
	goLeft: logic.Logic.TrickHasAFive,
	left: byValueTree,
	right: trumpLeadNoFiveTree,
}

var offsuitLeadNoFiveTree = &tree{
	logic: "If I am second to last and there's a 5 out I can cover",
	goLeft: func(l logic.Logic) bool {
		return l.NextPlayerIsLast() && l.AFiveIsOut() && !l.IHaveAFive() && l.ICanCoverAFive()
	},
	left: &tree{
		logic: "Score inverse of value for value above 5",
		score: forValuesAbove(card.Five, byNegValue),
	},
	right: byNegValueTree,
}

var trumpLeadNoFiveTree = &tree{
	logic: "If partner played the high card out",
	goLeft: logic.Logic.PartnerPlayedHighCardOut,
	left: byFivesThenNegValueTree,
	right: partnerNotHighCardOutTree,
}

var partnerNotHighCardOutTree = &tree{
	logic: "If partner is winning",
	goLeft: logic.Logic.PartnerPlayedHighCard,
	left: partnerHighCardNotHighCardOutTree,
	right: partnerNotHighCardTree,
}

var partnerHighCardNotHighCardOutTree = &tree{
	logic: "If I am last player",
	goLeft: logic.Logic.IAmLast,
	left: byFivesThenNegValueTree,
	right: byNegValueTree,
}

var partnerNotHighCardTree = &tree{
	logic: "If partner hasn't played yet",
	goLeft: logic.Logic.PartnerToPlay,
	left: partnerToPlayTree,
	right: parterPlayedNotWinningTree,
}

var partnerToPlayTree = &tree{
	logic: "If I have high card out",
	goLeft: logic.Logic.IHaveHighCardOut,
	left: byValueTree,
	right: partnerToPlayIDontHaveHighCardOutTree,
}

var partnerToPlayIDontHaveHighCardOutTree = &tree{
	logic: "If I can take lead",
	goLeft: logic.Logic.IHaveHighCard,
	left: &tree{
		logic: "Score inverse value where I can take lead",
		score: byValue, // TODO(drw)
	},
	right: byNegPointsThenNegValueTree,
}

var parterPlayedNotWinningTree = &tree{
	logic: "If there are more than 0 point showing",
	goLeft: logic.Logic.PointsAreShowing,
	left: pointsShowingParterFailedTree,
	right: noPointsPartnerFailedTree,
}

var pointsShowingParterFailedTree = &tree{
	logic: "If I can take lead",
	goLeft: logic.Logic.IHaveHighCard,
	left: byValueTree,
	right: byNegValueTree,
}

var noPointsPartnerFailedTree = &tree{
	logic: "If I am last",
	goLeft: logic.Logic.IAmLast,
	left: byNegPointsThenNegValueTree,
	right: noPointsPartnerFailedImNotLastTree,
}

var noPointsPartnerFailedImNotLastTree = &tree {
	logic: "If I can take lead",
	goLeft: logic.Logic.IHaveHighCard,
	left: byValueTree,
	right: byNegValueTree,
}

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
	logic   string
	left   *tree
	right  *tree
	goLeft func(l logic.Logic) bool
	score  func(c card.Card, t card.Suit) float64
}

// evaluate traverses the decision tree using the given logic and card,
// returning the score of the leaf we end up at. Also returns a string
func (t *tree) evaluate(l logic.Logic, c card.Card) (float64, string) {
	if t.left == nil && t.right == nil {
		return t.score(c, l.State.Trump), t.logic
	}
	if t.right == nil || t.goLeft(l) {
		value, s := t.left.evaluate(l, c)
		return value, t.logic + " <Y " + s
	}
	value, s := t.right.evaluate(l, c)
	return value, t.logic + " N> " + s
}

// scorerFromDT builds and returns a scorer from the given decision tree.
func scorerFromDT(t *tree) scorer {
	return func(s state.State, m action.Message, c card.Card) float64 {
		value, _ := t.evaluate(logic.Logic{s, m.Seat}, c)
		return value
	}
}

package playing

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"dr2w.com/hf/ai/logic"
	"dr2w.com/hf/model/action"
	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/hand"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/state"
)

const (
	posnReplicates = 100
	posnHandSize   = 10
)

var inconsistentlyTests = []struct {
	name string
	rate float64
	want []int
}{
	{
		"Rate 0",
		0.0,
		[]int{100},
	},
	{
		"Rate 0.5",
		0.5,
		[]int{45, 32, 12, 4, 5, 1, 0, 1},
	},
	{
		"Rate 0.2",
		0.2,
		[]int{81, 19},
	},
}

func count(ints []int) []int {
	counts := []int{}
	sort.Ints(ints)
	prev := -1
	for _, v := range ints {
		for v > prev {
			counts = append(counts, 0)
			prev++
		}
		counts[len(counts)-1] += 1
	}
	return counts
}

func posnScore(s state.State, m action.Message, c card.Card) float64 {
	for i := range *s.Hands[m.Seat] {
		return float64(i)
	}
	return -1
}

var posnMessage = action.Message{
	Seat:    seat.North,
	Options: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
}

var posnState = state.State{
	Hands: map[seat.Seat]*hand.Hand{
		seat.North: &hand.Hand{
			card.Card{Value: 0},
			card.Card{Value: 1},
			card.Card{Value: 2},
			card.Card{Value: 3},
			card.Card{Value: 4},
			card.Card{Value: 5},
			card.Card{Value: 6},
			card.Card{Value: 7},
			card.Card{Value: 8},
			card.Card{Value: 9},
		},
	},
}

func TestInconsistently(t *testing.T) {
	for _, test := range inconsistentlyTests {
		rand.Seed(0)
		d := inconsistently(test.rate, posnScore)
		results := make([]int, posnReplicates)
		for i := 0; i < posnReplicates; i++ {
			results[i] = d(posnState, posnMessage)[0]
		}
		got := count(results)
		gotStr, wantStr := fmt.Sprintf("%v", got), fmt.Sprintf("%v", test.want)
		if gotStr != wantStr {
			t.Errorf("%s: got %s, want %s", test.name, gotStr, wantStr)
		}
	}
}

var evaluateTests = []struct {
	name  string
	tree  *tree
	logic logic.Logic
	card  card.Card
	want  float64
}{
	{
		"Trivial Tree",
		&tree{score: 0.5},
		logic.Logic{},
		card.Card{},
		0.5,
	},
	{
		"Left Tree",
		&tree{
			left:  &tree{score: 0.1},
			right: &tree{score: 0.9},
			goLeft: func(i logic.Logic, c card.Card) bool {
				return true
			},
		},
		logic.Logic{},
		card.Card{},
		0.1,
	},
	{
		"Reasonably Complex",
		&tree{
			left: &tree{
				left:  &tree{score: 0.1},
				right: &tree{score: 0.2},
				goLeft: func(l logic.Logic, c card.Card) bool {
					return false
				},
			},
			right: &tree{
				left:  &tree{score: 0.8},
				right: &tree{score: 0.9},
				goLeft: func(l logic.Logic, c card.Card) bool {
					return true
				},
			},
			goLeft: func(l logic.Logic, c card.Card) bool {
				return false
			},
		},
		logic.Logic{},
		card.Card{},
		0.8,
	},
}

func TestEvaluate(t *testing.T) {
	for _, test := range evaluateTests {
		got := test.tree.evaluate(test.logic, test.card)
		if got != test.want {
			t.Errorf("%s: got %f, want %f", test.name, got, test.want)
		}
	}
}

var noisilyTests = []struct {
	name string
	rate float64
	want []int
}{
	{
		"Rate 0",
		0.0,
		[]int{100},
	},
	{
		"Rate 0.5",
		0.5,
		[]int{45, 32, 12, 4, 5, 1, 0, 1},
	},
	{
		"Rate 0.2",
		0.2,
		[]int{81, 19},
	},
}

// TODO(drw): modify from inconsistently to noisily.
func TestNoisily(t *testing.T) {
	for _, test := range noisilyTests {
		rand.Seed(0)
		d := noisily(test.rate, posnScore)
		results := make([]int, posnReplicates)
		for i := 0; i < posnReplicates; i++ {
			results[i] = d(posnState, posnMessage)[0]
		}
		got := count(results)
		gotStr, wantStr := fmt.Sprintf("%v", got), fmt.Sprintf("%v", test.want)
		if gotStr != wantStr {
			t.Errorf("%s: got %s, want %s", test.name, gotStr, wantStr)
		}
	}
}

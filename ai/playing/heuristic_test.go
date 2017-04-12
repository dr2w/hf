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
	"dr2w.com/hf/model/trick"
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
		[]int{0,0,0,0,0,0,0,0,0,100},
	},
	{
		"Rate 0.5",
		0.5,
        []int{0,0,1,0,1,5,4,12,32,45},
	},
	{
		"Rate 0.2",
		0.2,
        []int{0,0,0,0,0,0,0,0,19,81},
	},
	{
		"Rate 0.98",
		0.98,
        []int{12,9,9,6,8,13,11,10,13,9},
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
	for i, handCard := range *s.Hands[m.Seat] {
        if handCard == c {
		    return float64(i) / 10.0
        }
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

func scoreFunc(s float64) func(c card.Card, t card.Suit) float64 {
    return func(c card.Card, t card.Suit) float64 {
        return s
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
		&tree{score: scoreFunc(0.5) },
		logic.Logic{},
		card.Card{},
		0.5,
	},
	{
		"Left Tree",
		&tree{
			left:  &tree{score: scoreFunc(0.1)},
			right: &tree{score: scoreFunc(0.9)},
			goLeft: func(i logic.Logic) bool {
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
				left:  &tree{score: scoreFunc(0.1)},
				right: &tree{score: scoreFunc(0.2)},
				goLeft: func(l logic.Logic) bool {
					return false
				},
			},
			right: &tree{
				left:  &tree{score: scoreFunc(0.8)},
				right: &tree{score: scoreFunc(0.9)},
				goLeft: func(l logic.Logic) bool {
					return true
				},
			},
			goLeft: func(l logic.Logic) bool {
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
		got, s := test.tree.evaluate(test.logic, test.card)
		if got != test.want {
			t.Errorf("%s: got %f, want %f", test.name, got, test.want)
			t.Errorf("Tree Trace: %s", s)
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
		[]int{0,0,0,0,0,0,0,0,0,100},
	},
	{
		"Rate 0.01",
		0.01,
		[]int{0,0,0,0,0,0,0,0,0,100},
	},
	{
		"Rate 0.5",
		0.5,
		[]int{0,0,0,0,0,2,7,12,25,54},
	},
	{
		"Rate 0.2",
		0.2,
		[]int{0,0,0,0,0,0,0,6,24,70},
	},
	{
		"Rate 0.98",
		0.98,
		[]int{2,1,1,5,6,5,11,23,16,30},
	},
}

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

func makeTrick(first seat.Seat, cards []card.Card) trick.Trick {
	t := map[seat.Seat]card.Card{}
	st := first
	for _, c := range cards {
		t[st] = c
		st = st.Next()
	}
	return trick.Trick{t, first}
}

func makeHands(st seat.Seat, cards []card.Card) map[seat.Seat]*hand.Hand {
	h := hand.Hand(card.Set(cards))
	return map[seat.Seat]*hand.Hand{
		st: &h,
	}
}

var initialTreeTests = []struct {
	// Always N, always Clubs trump.
	name	string
	// Converts to tricks, last trick can have < 4.
	// Assigns arbitrarily to seats.
	out	string
	// Assumes cards are leading up to N's play
	trick	string
	// TODO(drw): Support offsuits in Shorthand
	hand	string
	best	int
}{
	{
		"Lead the Ace",
		"",
		"",
		"A852",
		0,
	},
	{
		"Five on partner's Ace",
		"",
		"A7",
		"K852",
		2,
	},
	{
		"Deuce on partner's Ace",
		"",
		"A7",
		"KJ832",
		4,
	},
	{
		"Not the five on opponent's Ace",
		"",
		"A",
		"J5",
		0,
	},
	{
		"Not a point on opponent's Ace",
		"",
		"A",
		"KJ5432",
		4,
	},
	{
		"Five on winning hand",
		"",
		"897",
		"AT6532",
		3,
	},
	{
		"Partner high card out",
		"AKQJ",
		"j8",
		"T96532",
		3,
	},
	{
		"Offload junk",
		"",
		"A",
		"KJ8532",
		4,
	},
	{
		"Always try to take the five",
		"",
		"f",
		"AJ832",
		0,
	},	
	{
		"Be wasteful", // TODO(drw): Fix to not be!
		"",
		"f48",
		"K932",
		0,
	},
}

func handFromShorthand(s string) *hand.Hand {
	cards := card.CardsFromShorthand(card.Clubs, s)
	h := hand.Hand(card.Set(cards))
	return &h
}

func outFromShorthand(s string) []trick.Trick {
	tricks := []trick.Trick{}
	cards := card.CardsFromShorthand(card.Clubs, s)
	for i := 0; i < len(cards); i += trick.Size {
		tricks = append(tricks, trick.New(cards[i:i+trick.Size]...))
	}
	if trick.Size * len(tricks) < len(cards) {
		remainder := cards[len(tricks)*trick.Size:]
		tricks = append(tricks, trick.New(remainder...))
	}
	return tricks
}

func currentFromShorthand(s string) trick.Trick {
	t := trick.Trick{Cards: map[seat.Seat]card.Card{}}
	cards := card.CardsFromShorthand(card.Clubs, s)
	st := seat.West
	i := len(cards)-1 
	for i >= 0 && st != seat.North {
		t.Cards[st] = cards[i]
		st = st.Previous()
		i -= 1
	}
	return t
}

func TestBasicTree(t *testing.T) {
	for _, test := range initialTreeTests {
		h := handFromShorthand(test.hand)
		out := outFromShorthand(test.out)
		current := currentFromShorthand(test.trick)
		state := state.State{
			Trump: card.Clubs,
			Hands: map[seat.Seat]*hand.Hand{seat.North: h},
			Played: append(out,current),
		}
		maxCardIndex := -1
		maxCardValue := -1.0
		maxCardTrace := ""
		for i, c := range *h {
			v, s := initialTree.evaluate(logic.Logic{state, seat.North}, c)
			if v > maxCardValue {
				maxCardValue = v
				maxCardIndex = i
				maxCardTrace = s
			}
		}
		if maxCardIndex != test.best {
			t.Errorf("%s: got %v, want %v", test.name, h.Get(maxCardIndex), h.Get(test.best))
			t.Errorf("Tree Trace: %s", maxCardTrace)
		}
	}
}

package action

import (
	"reflect"
	"testing"

	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/hand"
)

var discardsTests = []struct {
	name  string
	trump card.Suit
	hand  *hand.Hand
	want  []int
}{
	{
		"Empty Case",
		card.NoSuit,
		&hand.Hand{},
		nil,
	},
	{
		"Few Trump",
		card.Spades,
		&hand.Hand{
			card.Card{card.Seven, card.Clubs},
			card.Card{card.Eight, card.Diamonds},
			card.Card{card.Nine, card.Spades},
			card.Card{card.Ten, card.Hearts},
		},
		[]int{0, 1, 3},
	},
	{
		"Too Many Trump",
		card.Hearts,
		&hand.Hand{
			card.Card{card.Deuce, card.Hearts},
			card.Card{card.Three, card.Hearts},
			card.Card{card.OffFive, card.Hearts},
			card.Card{card.Seven, card.Hearts},
			card.Card{card.Joker, card.Hearts},
			card.Card{card.King, card.Hearts},
			card.Card{card.Ace, card.Hearts},
			card.Card{card.Ace, card.Spades},
		},
		[]int{1, 3, 5, 7},
	},
}

func TestDiscards(t *testing.T) {
    for _, test := range discardsTests {
        if got := discards(test.trump, test.hand); !reflect.DeepEqual(got, test.want) {
            t.Errorf("%s: got %v, want %v", test.name, got, test.want)
        }
    }
}

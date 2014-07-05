package hand

import (
	"reflect"
	"testing"

	"dr2w.com/hf/model/card"
)

var (
	c1 = card.Card{1, 1}
	c2 = card.Card{2, 2}
)

var addTests = []struct {
	name  string
	cards []card.Card
}{
	{
		"Empty Set",
		[]card.Card{},
	},
	{
		"Single Card",
		[]card.Card{c1},
	},
	{
		"Double Cards",
		[]card.Card{c1, c2},
	},
}

func TestAdd(t *testing.T) {
	for _, test := range addTests {
		h := Hand{}
		for _, card := range test.cards {
			h.Add(card)
		}
		if !reflect.DeepEqual(h, Hand(test.cards)) {
			t.Errorf("Want %v, got %v", Hand(test.cards), h)
		}
	}
}

var removeTests = []struct {
	name   string
	hand   *Hand
	remove int
	card   card.Card
	want   *Hand
	err    bool
}{
	{
		"Single Card",
		&Hand{c1},
		0,
		c1,
		&Hand{},
		false,
	},
	{
		"Remove second",
		&Hand{c1, c2},
		1,
		c2,
		&Hand{c1},
		false,
	},
	{
		"Remove First",
		&Hand{c1, c2},
		0,
		c1,
		&Hand{c2},
		false,
	},
	{
		"Removed nonexistant",
		&Hand{c1, c2},
		2,
		card.Card{},
		&Hand{c1, c2},
		true,
	},
}

func TestRemove(t *testing.T) {
	for _, test := range removeTests {
		card, err := test.hand.Remove(test.remove)
		if err != nil && !test.err {
			t.Errorf("%s: unexpected error (%s)", test.name, err)
		}
		if err == nil && test.err {
			t.Errorf("%s: expected error but got none.", test.name)
		}
		if card != test.card {
			t.Errorf("%s: returned wrong card. Want %v, got %v", test.name, test.card, card)
		}
		if !reflect.DeepEqual(test.hand, test.want) {
			t.Errorf("%s: want %v, got %v", test.name, test.want, test.hand)
		}
	}
}

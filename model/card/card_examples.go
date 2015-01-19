package card

import "strings"

// This file contains functions and exported variables
// that provide example Card Sets for use elsewhere in the
// code.

var (
	BestHand  = hand("AKQJjT982,,,")
	WorstHand = hand("643,43,43,43")

	BestSixHand  = hand("jQ5,Q5,Q5,Q5")
	SolidSixHand = hand("j652,T5,J9,Q")
	WeakSixHand  = hand("T42,9853,,Q4")

	GreatHand1 = hand("86,AKJjT,9,5")
	GoodHand1  = hand("j4,AQ92,5,83")
	GoodHand2  = hand("T3,,AKT87,A7")
	GoodHand3  = hand("AJT72,K2,8,T")
	GoodHand4  = hand("6,KQ,97,AQ72")

	OkHand1 = hand("Q2,J4,875,T4")

	BadHand1 = hand("642,T4,J,Q87")
)

var valueMap = map[rune]Value{
	'A': Ace,
	'K': King,
	'Q': Queen,
	'J': Jack,
	'j': Joker,
	'T': Ten,
	'9': Nine,
	'8': Eight,
	'7': Seven,
	'6': Six,
	'5': Five,
	'4': Four,
	'3': Three,
	'2': Deuce,
}

var suitOrder = []Suit{Spades, Hearts, Diamonds, Clubs}

// hand converts the given string into a card.Set.
// It looks for comma-separated character strings, where
// each character represents a card and the suits are
// separated by commas in the following order:
// Spades, Hearts, Diamonds, Clubs.
// ex.
// AKQ,T984,,j2
// represents nine cards with no diamonds and a very nice
// spades suit. The joker can be included as part of any
// suit and will be parsed correctly.
func hand(s string) Set {
	cards := Set{}
	for i, values := range strings.Split(s, ",") {
		for _, v := range values {
			suit := suitOrder[i]
			if v == 'j' {
				suit = NoSuit
			}
			cards = append(cards, Card{valueMap[v], suit})
		}
	}
	return cards
}

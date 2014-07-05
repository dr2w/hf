// Package bid models players' bids
package bid

type Bid int

func (b Bid) String() string {
    return Names[b]
}

// Score returns the score gained by taking the given points during a hand.
func (b Bid) Score(points int) int {
    if points < Points[b] {
        return -Scores[b]
    }
    if points > Scores[b] {
        return points
    }
    return Scores[b]
}

const (
	Pass Bid = iota
	B6
	B7
	B8
	B9
	B10
	B11
	B12
	B13
	B14
	B1428
	B15
	B1530
)

// Values enumerates all possible Bid Values in increasing rank order.
var Values = []Bid{Pass, B6, B7, B8, B9, B10, B11, B12, B13, B14, B1428, B15, B1530}

var Names = map[Bid]string {
	Pass: "Pass",
	B6: "6",
	B7: "7",
	B8: "8",
	B9: "9",
	B10: "10",
	B11: "11",
	B12: "12",
	B13: "13",
	B14: "14",
	B1428: "14/28",
	B15: "15",
	B1530: "15/30",
}

var Scores = map[Bid]int {
	B6: 6,
	B7: 7,
	B8: 8,
	B9: 9,
	B10: 10,
	B11: 11,
	B12: 12,
	B13: 13,
	B14: 14,
	B1428: 28,
	B15: 15,
	B1530: 30,
}

var Points = map[Bid]int {
	B6: 6,
	B7: 7,
	B8: 8,
	B9: 9,
	B10: 10,
	B11: 11,
	B12: 12,
	B13: 13,
	B14: 14,
	B1428: 14,
	B15: 15,
	B1530: 15,
}

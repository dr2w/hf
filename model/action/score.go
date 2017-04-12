package action

import (
    "dr2w.com/hf/model/state"
    "dr2w.com/hf/model/bid"
    "dr2w.com/hf/model/trick"
    "dr2w.com/hf/model/seat"
    "dr2w.com/hf/model/card"
)

func points(tricks []trick.Trick, trump card.Suit) map[seat.Seat]int {
    p := make(map[seat.Seat]int)
    for _, t := range tricks {
        s, _ := t.Winner(trump)
        points := t.Points(trump)
        p[s] += points
        p[s.Partner()] += points
    }
    return p
}

func resolve(s seat.Seat, b bid.Bid, points map[seat.Seat]int) map[seat.Seat]int {
    scores := make(map[seat.Seat]int)
    nonPts := points[s.Next()]
    bidScore := b.Score(points[s])
    scores[s], scores[s.Partner()] = bidScore, bidScore
    scores[s.Next()], scores[s.Partner().Next()] = nonPts, nonPts
    return scores
}

// score increments the scores appropriately and resets the rest of the
// state for the next hand.
func score(s state.State, m Message) (state.State, Message, error) {
    bs, b := s.WinningBid()
    p := points(s.Played, s.Trump)
    //fmt.Printf("Points => %v\n", p)
    scores := resolve(bs, b, p)
    //fmt.Printf("Scores => %v\n", scores)
    next := s.NextRound(scores)
    return next,
	   Message{Deal, next.Dealer, []int{0}, 1}, nil
}

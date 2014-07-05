package action

import (
	"fmt"

	"dr2w.com/hf/model/bid"
	"dr2w.com/hf/model/card"
	"dr2w.com/hf/model/seat"
	"dr2w.com/hf/model/state"
)

// bids takes a State and a selected bid and returns the next State and Action.
func bids(s state.State, m Message) (state.State, Message, error) {
	if s.Bids == nil {
		s.Bids = make(map[seat.Seat]bid.Bid)
	}

	sel, err := m.Selection()
	if err != nil {
		return state.State{}, Message{}, err
	}
	if sel >= len(bid.Values) {
		return state.State{}, Message{}, fmt.Errorf("unable to process invalid bid Selection (%d) for %s", sel, s)
	}
	s.Bids[m.Seat] = bid.Values[sel]

	if len(s.Bids) == len(seat.Order) {
		return s, reqForChooseSuit(s), nil
	}
	req, err := reqForNextBid(s)
	return s, req, err
}

// bidWinner takes a map of Seat to Bid and returns the Seat and Bid which won the bidding.
func bidWinner(m map[seat.Seat]bid.Bid) (s seat.Seat, b bid.Bid) {
	for seat, bid := range m {
		if bid >= b {
			b = bid
			s = seat
		}
	}
	return s, b
}

// reqForNextBid takes a state which is ready for the next bid and returns the
// corresponding Message.
func reqForNextBid(s state.State) (Message, error) {
	st, b := bidWinner(s.Bids)
	if st == seat.None {
		return Message{}, fmt.Errorf("reqForBid called on State with no Bids (%s).", s)
	}
	if len(s.Bids) >= len(seat.Order) {
		return Message{}, fmt.Errorf("reqForBid called with all bids made (%s).", s)
	}
	for {
		st = st.Next()
		if _, ok := s.Bids[st]; !ok {
			break
		}
	}
	return Message{
		Type:    Bid,
		Seat:    st,
		Options: append([]int{int(bid.Pass)}, SelectionRange(int(b)+1, len(bid.Values))...),
	}, nil
}

// reqForChooseSuit takes a state which has all bids completed and
// returns the corresponding playem Message.
func reqForChooseSuit(s state.State) Message {
	seat, _ := bidWinner(s.Bids)
	return Message{
		Type:    Trump,
		Seat:    seat,
		Options: SelectionRange(0, len(card.Suits)), // 4 Suits
	}
}

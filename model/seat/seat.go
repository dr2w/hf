// Package seat describes the position of a Hand relative to that of the other
// Hands.
package seat

// Seat represents a side of the table where a player would sit.
type Seat int

// Next returns the next Seat to take a turn in clockwise order.
func (s Seat) Next() Seat {
	for i, seat := range Order {
		if seat == s {
			return Order[(i+1)%len(Order)]
		}
	}
	return None
}

// Partner returns the partner of this Seat.
func (s Seat) Partner() Seat {
    return s.Next().Next()
}

// String returns a human readable representation of the seat
func (s Seat) String() string {
	return Names[s]
}

const (
	None Seat = iota
	North
	South
	East
	West
)

var Names = map[Seat]string{
	None:  "X",
	North: "North",
	South: "South",
	East:  "East",
	West:  "West",
}

// Order defines the set of valid Seats in play order.
var Order = []Seat{North, East, South, West}

package seat

import (
	"testing"
)

var seatNextTests = []struct {
	seat Seat
	want Seat
}{
	{
		East,
		South,
	},
	{
		South,
		West,
	},
	{
		West,
		North,
	},
	{
		North,
		East,
	},
}

func TestSeatNext(t *testing.T) {
	for _, test := range seatNextTests {
		if got := test.seat.Next(); got != test.want {
			t.Errorf("Got %s for %s, want %s", got, test.seat, test.want)
		}
	}
}

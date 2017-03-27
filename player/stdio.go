package player

import (
    "fmt"
    "strings"
    "strconv"
    "time"

    "dr2w.com/hf/model/action"
    "dr2w.com/hf/model/bid"
    "dr2w.com/hf/model/card"
    "dr2w.com/hf/model/seat"
    "dr2w.com/hf/model/state"
)

var clearLines = 40

type Stdio struct {
    Seat seat.Seat
}

// Play prints the relevant State and Message Options to stdout and pulls the selection
// from stdin.
func (p Stdio) Play(s state.State, m action.Message) []int {
    clearScreen()
    p.Seat = m.Seat
    displayState(s, m.Type, p.Seat)
    return solicitChoice(m, s)
}

// Update prints the new relevant State to stdout.
func (p Stdio) Update(s state.State, t action.Type) {
    clearScreen()
    displayState(s, t, p.Seat)
    time.Sleep(2*time.Second)
}

// clearScreen scrolls the output to make way for a new update.
func clearScreen() {
    for i := 0; i < clearLines; i++ {
        fmt.Printf("\n")
    }
}

// displayState prints the current game state to stdout in an informative format. 
func displayState(s state.State, t action.Type, st seat.Seat) {
    if t == action.Bid || t == action.Trump {
        displayBidding(s, st)
    } else {
        displayPlays(s, st)
    }
}

// displayBidding shows the current scores as well as the bids placed so far.
func displayBidding(s state.State, st seat.Seat) {
    displayScores(s)
    displayBids(s)
    displayHand(s, st, []int{})
}

// displayScores prints out a single line with the current scores.
func displayScores(s state.State) {
    fmt.Printf("Score:\nEast/West: %d\tNorth/South: %d\n", s.Score[seat.East], s.Score[seat.North])
}

func displayBids(s state.State) {
    fmt.Printf("\nN\tE\tS\tW\n")
    fmt.Printf("%s\t%s\t%s\t%s\n",
               s.Bids[seat.North].String(),
               s.Bids[seat.East].String(),
               s.Bids[seat.South].String(),
               s.Bids[seat.West].String())
}

// displayPlays shows the trump and bid as well as the current hand and current play.
func displayPlays(s state.State, st seat.Seat) {
    displayWinningBid(s)
    displayPlay(s)
}

// displayWinningBid shows the suit and value of the winning bid.
func displayWinningBid(s state.State) {
    st, bid := s.WinningBid()
    if st == seat.None {
        return
    }
    fmt.Printf("Bid is %s %s by %s\n", bid, s.Trump, st)
}

// displayPlay shows the state of the current face-up play.
func displayPlay(s state.State) {
    if len(s.Played) == 0 {
        return
    }
    lastTrick := s.Played[len(s.Played)-1]
    fmt.Printf("\nPlay:\n")
    cardString := func(st seat.Seat) string {
	c := lastTrick.Cards[st]
	if c.Value == card.NoValue {
		return ""
	}
	return c.String()
    }
    fmt.Printf("            N: %s\n", cardString(seat.North))
    fmt.Printf("W: %s                      E: %s\n",
               cardString(seat.West), cardString(seat.East))
    fmt.Printf("            S: %s\n", cardString(seat.South))
}

// displayHand shows the given seat's hand.
// options is an array of ints indicating which indices in the hand are selectable.
func displayHand(s state.State, st seat.Seat, options []int) {
    if s.Hands[st] == nil {
        return
    }
    fmt.Printf("\nHand:\n")
    var suit card.Suit
    cards := []string{}
    for i, c := range *s.Hands[st] {
         if c.Suit != suit {
             if len(cards) > 0 {
                 fmt.Printf("%s: %s\n", suit.String(), strings.Join(cards, ","))
		 cards = []string{}
             }
             suit = c.Suit
         }
	 cardString := c.Value.String()
	 for _, o := range options {
		if o == i {
			cardString = cardString + "[" + strconv.Itoa(i) + "]"
		}
	 }
	 cards = append(cards, cardString)
    }
    if len(cards) > 0 {
        fmt.Printf("%s: %s", suit.String(), strings.Join(cards, ","))
    }
}

// displayBid prints out a selection of available bids along with the
// numerical selection values.
func displayBid(m action.Message) {
	names := make([]string, len(m.Options))
	values := make([]string, len(m.Options))
	for i := range m.Options {
	    names[i] = bid.Bid(i).String()
	    values[i] = "[" + strconv.Itoa(i) + "]"
        }
	fmt.Printf("\n\n%s\n%s\n", strings.Join(names, "\t"), strings.Join(values, "\t"))
	fmt.Printf("Please select a bid: ")
}

// displayTrumpOptions prints out the player's hand along with numerical
// suit selection values.
func displayTrumpOptions(m action.Message, s state.State) {
	var suitOptions []string
	for i := range m.Options {
		suitName := card.Suits[i].String()
		suitOptions = append(suitOptions, suitName + "[" + strconv.Itoa(i) + "]")
	}
	fmt.Printf("\nSuit: " + strings.Join(suitOptions, ","))
	fmt.Printf("\nPlease select a suit: ")
}

// displayChoice shows the options the current player can choose from.
func displayChoice(m action.Message, s state.State) {
    switch m.Type {
	case action.Bid:
		displayBid(m)
	case action.Play:
		displayHand(s, m.Seat, m.Options)
		fmt.Printf("\n\nPlease select a card to play:")
	case action.Discard:
		displayHand(s, m.Seat, m.Options)
		toDiscard := len(card.Set(*s.Hands[m.Seat])) - 6
		fmt.Printf("\n\nPlease select %d cards to discard (use commas):", toDiscard)
	case action.Trump:
		displayTrumpOptions(m, s)
	default:
    		fmt.Printf("\n\n%s: ", m)
    }
}

// solicitChoice prompts the user to select one or more of a set of options and returns
// the selections.
func solicitChoice(m action.Message, s state.State) []int {
    displayChoice(m, s)
    text := ""
    _, err := fmt.Scanln(&text)
    if err != nil {
            fmt.Println("error when reading from stdin, please try again.")
            return solicitChoice(m, s)
    }
    selections := strings.Split(text, ",")
    result := make([]int, len(selections))
    for i, sel := range selections {
        result[i], err = strconv.Atoi(sel)
        if err != nil {
            fmt.Printf("can't interpret %q as a number:\n%s\n", result[i], err)
            return solicitChoice(m, s)
        }
    }
    fmt.Printf("Selected: %v", result)
    return result
}

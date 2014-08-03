package player

import (
    "fmt"
    "strings"
    "strconv"
    "time"

    "dr2w.com/hf/model/state"
    "dr2w.com/hf/model/seat"
    "dr2w.com/hf/model/action"
)

type Stdio struct {
    Seat seat.Seat
}

// Play prints the relevant State and Message Options to stdout and pulls the selection
// from stdin.
func (p Stdio) Play(s state.State, m action.Message) []int {
    clearScreen()
    p.Seat = m.Seat
    displayState(s, m.Type, p.Seat)
    return solicitChoice(m)
}

// Update prints the new relevant State to stdout.
func (p Stdio) Update(s state.State, t action.Type) {
    clearScreen()
    displayState(s, t, p.Seat)
    time.Sleep(2*time.Second)
}

// clearScreen scrolls the output to make way for a new update.
func clearScreen() {
    for i := 0; i < 10; i++ {
        fmt.Printf("\n")
    }
}

// displayState prints the current game state to stdout in an informative format. 
func displayState(s state.State, t action.Type, st seat.Seat) {
    if t == action.Bid || t == action.Trump {
        displayBidding(s)
    }
    displayPlays(s, st)
}

// displayBidding shows the current scores as well as the bids placed so far.
func displayBidding(s state.State) {
    displayScores(s)
    displayBids(s)
}

// displayScores prints out a single line with the current scores.
func displayScores(s state.State) {
    fmt.Printf("East/West: %d\tNorth/South: %d\n", s.Score[seat.East], s.Score[seat.North])
}

func displayBids(s state.State) {
    fmt.Printf("\nN\tE\tS\tW\n")
    fmt.Printf("%d\t%d\t%d\t%d\n", s.Bids[seat.North],s.Bids[seat.East],s.Bids[seat.South],s.Bids[seat.West])
}

// displayPlays shows the trump and bid as well as the current hand and current play.
func displayPlays(s state.State, st seat.Seat) {
    displayWinningBid(s)
    displayPlay(s)
    displayHand(s, st)
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
    fmt.Printf("\nPlay: %s\n", lastTrick)
}

// displayHand shows the controlling (not active) player's hand.
func displayHand(s state.State, st seat.Seat) {
    if s.Hands[st] == nil {
        return
    }
    fmt.Printf("\nHand: %s\n", s.Hands[st])
}

// solicitChoice prompts the user to select one or more of a set of options and returns
// the selections.
func solicitChoice(m action.Message) []int {
    fmt.Printf("\n%s: ", m)
    text := ""
    _, err := fmt.Scanln(&text)
    if err != nil {
            fmt.Println("error when reading from stdin, please try again.")
            return solicitChoice(m)
    }
    selections := strings.Split(text, ",")
    result := make([]int, len(selections))
    for i, s := range selections {
        result[i], err = strconv.Atoi(s)
        if err != nil {
            fmt.Printf("can't interpret %q as a number:\n%s\n", result[i], err)
            return solicitChoice(m)
        }
    }
    return result
}

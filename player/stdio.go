package player

import (
    "fmt"
    "strings"
    "strconv"

    "dr2w.com/hf/model/state"
    "dr2w.com/hf/model/action"
)

type Stdio struct {
    Nothing string
}

// Play prints the relevant State and Message Options to stdout and pulls the selection
// from stdin.
func (p Stdio) Play(s state.State, m action.Message) []int {
    displayState(s, m.Type)
    return solicitChoice(m)
}

// displayState prints the current game state to stdout in an informative format. 
func displayState(s state.State, t action.Type) {
    if t == action.Bid || t == action.Trump {
        displayBidding(s)
    }
    displayPlays(s)
}

// displayBidding shows the current scores as well as the bids placed so far.
func displayBidding(s state.State) {

}

// displayPlays shows the trump and bid as well as the current hand and current play.
func displayPlays(s state.State) {

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

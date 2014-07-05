package player

import (
    "os"
    "io/ioutil"
    "fmt"
    "log"
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
    displayState(s)
    return solicitChoice(m)
}

// displayStatus prints the current game state to stdout in an informative format. 
func displayState(s state.State) {
    fmt.Println(s)
}

// solicitChoice prompts the user to select one or more of a set of options and returns
// the selections.
func solicitChoice(m action.Message) []int {
    fmt.Print("%s %v: ", m.Type, m.Options)
    bytes, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        log.Fatalf("unable to read from stdin: %s", err)
    }
    response := string(bytes)
    selections := strings.Split(response, ",")
    result := make([]int, len(selections))
    for i, s := range selections {
        result[i], err = strconv.Atoi(s)
        if err != nil {
            fmt.Println("can't interpret %q as a number:\n%s", result[i], err)
            return solicitChoice(m)
        }
    }
    return result
}

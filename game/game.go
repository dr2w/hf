//
package game

import (
    "fmt"
   // "log"
    "bytes"
    "strings"

    "dr2w.com/hf/model/seat"
    "dr2w.com/hf/player"
    "dr2w.com/hf/model/state"
    "dr2w.com/hf/model/action"
)

const (
    winningScore = 52
    terribleScore = -104
)

type Game struct {
    Players map[seat.Seat]player.Player
    State state.State
    Message action.Message
}

func (g *Game) String() string {
    var(
        b bytes.Buffer
        names []string
    )
    for seat, player := range g.Players {
        names = append(names, fmt.Sprintf("%s: %s", seat, player))
    }
    b.WriteString(strings.Join(names, " | ") + "\n")
    b.WriteString(fmt.Sprintf("State:\n%s", g.State))
    b.WriteString(fmt.Sprintf("Message:\n%s", g.Message))
    return b.String()
}

// Over returns true iff the Game's state is a terminal one.
func (g *Game) Over() bool {
    for _, score := range g.State.Score {
        if score > winningScore || score < terribleScore {
            return true
        }
    }
    return false
}

// Resolve executes the game to an end state.
func (g *Game) Resolve() error {
    //log.Printf("Starting Game:\n%s", g)
    for !g.Over() {
        if err := g.Advance(); err != nil {
            return err
        }
        for _, p := range g.Players {
            p.Update(g.State, g.Message.Type)
        }
        //log.Printf("Game Advanced to:\n%s", g)
    }
    return nil
}

// Advance advances the Game one step.
func (g *Game) Advance() error {
    if g.Message.Seat != seat.None {
        p := g.Players[g.Message.Seat]
        response := p.Play(g.State, g.Message)
        //log.Printf("Player %s chose %v", p, response)
        g.Message.Options = response
    }
    s, m, err := action.NextState(g.State, g.Message)
    if err != nil {
        return err
    }
    g.State, g.Message = s, m
    return nil
}

// initialMessage returns the initial deal message given the first player.
func initialMessage(first seat.Seat) action.Message {
    return action.Message{
            Type: action.Deal,
            Seat: first,
            Options: []int{0},
        }
}

// New returns a Game initialized to the starting state for the given first
// seat and players (ordered by seat.Order).
func New(first seat.Seat, players ...player.Player) (*Game, error) {
    if len(players) != len(seat.Order) {
        return nil, fmt.Errorf("invalid number of players (%d) supplied to game.New.", len(players))
    }
    p := make(map[seat.Seat]player.Player)
    for i, s := range seat.Order {
        p[s] = players[i]
    }
    return &Game{
        Players: p,
        State: state.Initial(first),
        Message: initialMessage(first),
    }, nil
}

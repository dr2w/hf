package main

import (
    "log"

    "dr2w.com/hf/game"
    "dr2w.com/hf/ai"
    "dr2w.com/hf/player"
    "dr2w.com/hf/model/seat"
)

func main() {
    for i := 0; i < 20000; i++ {
        g, _ := game.New(
            seat.East,
            ai.Dumb,
            ai.Dumb,
            ai.Dumb,
            player.Stdio{},
        )
        err := g.Resolve()
        if err != nil {
            log.Fatalf("Error in Resolving: %s\n%s", err, g)
        }
        log.Printf("Score: %v (%d rounds)", g.State.Score, g.State.Rounds)
    }
}

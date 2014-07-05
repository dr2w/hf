package ai

import (
    "sort"
    "math/rand"

    "dr2w.com/hf/model/action"
    "dr2w.com/hf/model/state"
)

// first implements Decider and returns the first available option.
func first(s state.State, m action.Message) []int {
    return []int{m.Options[0]}
}

// second implements Decider and returns the second available option.
func second(s state.State, m action.Message) []int {
    if len(m.Options) > 1 {
        return []int{m.Options[1]}
    }
    return first(s,m)
}

// last implements Decider and returns the last available option.
func last(s state.State, m action.Message) []int {
    return []int{m.Options[len(m.Options)-1]}
}

// rand1 implements Decider and returns a random available option.
func rand1(s state.State, m action.Message) []int {
    return []int{m.Options[rand.Intn(len(m.Options))]}
}

// randN generates a Decider that chooses N random elements from
// the Options.
func randN(n int) func(s state.State, m action.Message) []int {
    return func(s state.State, m action.Message) []int {
        perm := rand.Perm(len(m.Options))
        sort.Ints(perm)
        return perm[:n]
    }
}

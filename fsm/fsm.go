// Package fsm implements a simple finite-state machine. There is no input field. If you have several permission groups, consider creating one FSM for each group.
package fsm

type Transition[State comparable] struct {
	From State
	To   State
}

type FSM[State comparable] []Transition[State] // easier than a map for a small number of transactions

func (fsm *FSM[State]) Can(from, to State) bool {
	for _, t := range *fsm {
		if t.From == from && t.To == to {
			return true
		}
	}
	return false
}

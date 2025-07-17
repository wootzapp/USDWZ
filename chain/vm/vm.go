package vm

import "strconv"

// Opcode defines a VM instruction.
type Opcode int

const (
	OP_SET Opcode = iota
	OP_ANY
	OP_QUORUM
	OP_ALL
	OP_THEN
)

// Instruction represents a single VM instruction.
type Instruction struct {
	Op  Opcode
	Arg string
}

// VM executes validator approval scripts.
type VM struct{}

// New returns a new VM instance.
func New() *VM { return &VM{} }

// Validate executes the script against the provided votes and validator sets.
// Votes map validator IDs to their boolean vote. Sets map set IDs to validator
// IDs.
func (vm *VM) Validate(script []Instruction, votes map[string]bool, sets map[string][]string) bool {
	i := 0
	for i < len(script) {
		inst := script[i]
		if inst.Op != OP_SET {
			return false
		}
		current := sets[inst.Arg]
		i++
		if i >= len(script) {
			return false
		}
		// evaluate condition for current set
		cond := script[i]
		ok := false
		switch cond.Op {
		case OP_ANY:
			ok = hasAny(current, votes)
		case OP_QUORUM:
			n, err := strconv.Atoi(cond.Arg)
			if err != nil {
				return false
			}
			ok = countYes(current, votes) >= n
		case OP_ALL:
			ok = countYes(current, votes) == len(current)
		default:
			return false
		}
		if !ok {
			return false
		}
		i++
		if i < len(script) && script[i].Op == OP_THEN {
			i++
			continue
		}
	}
	return true
}

func hasAny(set []string, votes map[string]bool) bool {
	for _, v := range set {
		if votes[v] {
			return true
		}
	}
	return false
}

func countYes(set []string, votes map[string]bool) int {
	c := 0
	for _, v := range set {
		if votes[v] {
			c++
		}
	}
	return c
}

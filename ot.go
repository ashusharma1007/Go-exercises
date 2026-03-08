package main

import "fmt"

// Operation types
const (
	OpInsert = "insert"
	OpDelete = "delete"
	OpRetain = "retain"
)

// Operation represents a single OT operation
type Operation struct {
	Type   string `json:"type"`
	Pos    int    `json:"pos"`
	Char   string `json:"char,omitempty"`
	Length int    `json:"length,omitempty"`
}

// Document represents the shared document state
type Document struct {
	Content string
	Version int
}

// Apply applies an operation to the document
func (d *Document) Apply(op Operation) error {
	switch op.Type {
	case OpInsert:
		if op.Pos < 0 || op.Pos > len(d.Content) {
			return fmt.Errorf("invalid insert position: %d", op.Pos)
		}
		d.Content = d.Content[:op.Pos] + op.Char + d.Content[op.Pos:]
		d.Version++
	case OpDelete:
		if op.Pos < 0 || op.Pos >= len(d.Content) {
			return fmt.Errorf("invalid delete position: %d", op.Pos)
		}
		d.Content = d.Content[:op.Pos] + d.Content[op.Pos+1:]
		d.Version++
	default:
		return fmt.Errorf("unknown operation type: %s", op.Type)
	}
	return nil
}

// Transform transforms operation op1 against op2
// This is the core of Operational Transformation
// Returns the transformed version of op1
func Transform(op1, op2 Operation) Operation {
	transformed := op1

	// If both operations are at the same position
	if op1.Pos == op2.Pos {
		if op1.Type == OpInsert && op2.Type == OpInsert {
			// Both inserting at same position - second one shifts right
			transformed.Pos = op1.Pos + 1
		}
		// Delete vs Insert or Delete vs Delete - no transformation needed
		return transformed
	}

	// op2 happened before op1's position
	if op2.Pos < op1.Pos {
		if op2.Type == OpInsert {
			// Insert before op1 - shift op1 right
			transformed.Pos = op1.Pos + 1
		} else if op2.Type == OpDelete {
			// Delete before op1 - shift op1 left
			transformed.Pos = op1.Pos - 1
			if transformed.Pos < 0 {
				transformed.Pos = 0
			}
		}
	}

	// op2 happened after op1's position - no transformation needed
	return transformed
}

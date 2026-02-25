// Package clause provides data structures and operations for working with propositional logic clauses.
package clause

import (
	"fmt"
	"strings"

	"github.com/thxrsxm/res/internal/utils"
)

// Literal represents a propositional variable or its negation.
type Literal int

// ErrorLiteral represents an invalid or unknown literal.
const ErrorLiteral Literal = 0

// lit2Str maps literal values to their string representations.
var lit2Str = []byte{
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

// Lit2Str converts a Literal to its string representation.
// For positive literals, returns the letter (e.g., 1 → "A").
// For negative literals, returns the letter with a minus prefix (e.g., -1 → "-A").
// For ErrorLiteral (0), returns "?".
func Lit2Str(l Literal) string {
	if l == ErrorLiteral {
		return "?"
	}
	s := ""
	if l < 0 {
		s = "-"
		l = -l
	}
	if l < 1 || int(l) > len(lit2Str) {
		return "?"
	}
	return s + string(lit2Str[l-1])
}

// Str2Lit converts a string to a Literal.
// Accepts strings like "A", "a", "-A", "-a".
// Returns ErrorLiteral for invalid input.
func Str2Lit(s string) Literal {
	if len(s) == 0 || len(s) > 2 {
		return ErrorLiteral
	}
	signed := false
	if len(s) == 2 {
		if s[0] != '-' {
			return ErrorLiteral
		}
		signed = true
		s = s[1:]
	}
	if len(s) != 1 {
		return ErrorLiteral
	}
	l := rune(strings.ToUpper(s)[0])
	if l < 'A' || l > 'Z' {
		return ErrorLiteral
	}
	result := Literal(l - 'A' + 1)
	if signed {
		return -result
	}
	return result
}

// Clause represents a disjunction of literals (A ∨ B ∨ ¬C, etc.).
// Internally, it's implemented as a set of literals using a map for efficient operations.
type Clause struct {
	literals map[Literal]struct{}
}

// New creates and returns a new empty Clause.
func New() *Clause {
	return &Clause{literals: make(map[Literal]struct{})}
}

// Contains checks if the clause contains the given literal.
func (c *Clause) Contains(l Literal) bool {
	_, ok := c.literals[l]
	return ok
}

// IsEmpty checks if the clause is empty (contains no literals).
// An empty clause represents a contradiction and is always false.
func (c *Clause) IsEmpty() bool {
	return len(c.literals) == 0
}

// Insert adds a literal to the clause.
// Returns true if the literal was added, false if it was already present or
// if adding it created a contradiction (literal and its negation).
// If a contradiction is created, the conflicting literals are removed.
func (c *Clause) Insert(l Literal) bool {
	if c.Contains(-l) {
		delete(c.literals, -l)
		return false
	}
	if !c.Contains(l) {
		c.literals[l] = struct{}{}
		return true
	}
	return false
}

// Size returns the number of literals in the clause.
func (c *Clause) Size() int {
	return len(c.literals)
}

// Resolve applies the resolution rule between this clause and another.
// Resolution rule: If two clauses contain complementary literals (A and ¬A),
// those literals can be removed and the remaining literals form a new clause (the resolvent).
//
// Returns:
//   - The resolvent clause (union of both clauses minus the complementary literals)
//   - A boolean indicating whether resolution actually occurred (complementary literals were found)
//
// If multiple complementary pairs exist, resolution fails and returns (nil, false).
func (c *Clause) Resolve(other Clause) (*Clause, bool) {
	found := false
	temp := other.Copy()
	for key := range c.literals {
		size := temp.Size()
		temp.Insert(key)
		if temp.Size() < size {
			if found {
				return nil, false
			}
			found = true
		}
	}
	return temp, found
}

// Copy creates and returns a deep copy of the clause.
func (c *Clause) Copy() *Clause {
	temp := New()
	for key := range c.literals {
		temp.Insert(key)
	}
	return temp
}

// Equals checks if this clause is equal to another clause.
// Two clauses are equal if they contain exactly the same literals.
func (c *Clause) Equals(other Clause) bool {
	if c.Size() != other.Size() {
		return false
	}
	for l := range c.literals {
		if !other.Contains(l) {
			return false
		}
	}
	return true
}

// String returns a string representation of the clause in set notation.
// Literals are sorted by their absolute value for consistent output.
// Example: A clause containing literals B, -A, and C would be represented as "{-A, B, C}".
func (c *Clause) String() string {
	values := make([]int, 0, len(c.literals))
	for l := range c.literals {
		values = append(values, int(l))
	}
	utils.UnsignedSort(values)
	s := "{"
	for i, l := range values {
		s += Lit2Str(Literal(l))
		if i < len(values)-1 {
			s += ", "
		}
	}
	return s + "}"
}

// Parse parses a string representation of a clause into a Clause object.
// The input string should be comma-separated literals (e.g., "A,B,-C").
// Returns an error if the input format is invalid.
func Parse(s string) (*Clause, error) {
	if len(s) == 0 {
		return nil, fmt.Errorf("input is empty")
	}
	c := New()
	for _, v := range strings.Split(s, ",") {
		// Trim whitespace from each component
		v = strings.TrimSpace(v)
		// Check length: 1 character (letter) or 2 characters with leading '-'
		if len(v) == 0 || len(v) > 2 || (len(v) == 2 && v[0] != '-') {
			return nil, fmt.Errorf("wrong clause format: %q", v)
		}
		// Check for negative sign
		signed := false
		if len(v) == 2 && v[0] == '-' {
			signed = true
			v = v[1:] // Keep only the letter part
		}
		// Validate we have exactly one character left
		if len(v) != 1 {
			return nil, fmt.Errorf("wrong clause format: %q", v)
		}
		// Convert to literal
		l := Str2Lit(v)
		if l == ErrorLiteral {
			return nil, fmt.Errorf("unknown symbol: %q", v)
		}
		if signed {
			l = -l
		}
		c.Insert(l)
	}
	return c, nil
}

// Res checks if a set of clauses is unsatisfiable using the resolution method.
// It returns true if the clause set is unsatisfiable (contains a contradiction),
// false if it is satisfiable (no contradiction found).
//
// The algorithm works by:
// 1. Starting with the initial set of clauses
// 2. Applying the resolution rule to derive new clauses
// 3. Checking if the empty clause is derived (indicating unsatisfiability)
// 4. Repeating until either the empty clause is found or no new clauses can be derived
//
// Parameters:
//   - set: The set of clauses to check
//   - index: The starting index for resolution (used internally for recursion)
func Res(set []Clause, index int) bool {
	size := len(set)
	for i := len(set) - 1; i >= 0; i-- {
		if set[i].IsEmpty() {
			return true
		}
		for k := len(set) - 1; k >= index; k-- {
			if i == k {
				continue
			}
			c, resolved := set[i].Resolve(set[k])
			if resolved && c != nil {
				// Check if the resolvent is already in the set
				exists := false
				for j := range set {
					if c.Equals(set[j]) {
						exists = true
						break
					}
				}
				if !exists {
					set = append(set, *c)
					// If we found an empty clause, return immediately
					if c.IsEmpty() {
						return true
					}
				}
			}
		}
	}
	for i := size; i < len(set); i++ {
		if set[i].IsEmpty() {
			return true
		}
	}
	if size == len(set) {
		return false
	}
	return Res(set, size)
}

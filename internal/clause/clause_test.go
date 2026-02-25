package clause

import (
	"testing"
)

func TestLit2Str(t *testing.T) {
	tests := []struct {
		name     string
		input    Literal
		expected string
	}{
		{"A", 1, "A"},
		{"B", 2, "B"},
		{"Z", 26, "Z"},
		{"negative A", -1, "-A"},
		{"negative B", -2, "-B"},
		{"negative Z", -26, "-Z"},
		{"zero", 0, "?"},
		{"too large positive", 27, "?"},
		{"too large negative", -27, "?"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Lit2Str(tt.input)
			if result != tt.expected {
				t.Errorf("Lit2Str(%d) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStr2Lit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Literal
	}{
		{"A", "A", 1},
		{"B", "B", 2},
		{"Z", "Z", 26},
		{"lowercase a", "a", 1},
		{"lowercase z", "z", 26},
		{"negative A", "-A", -1},
		{"negative Z", "-Z", -26},
		{"negative lowercase a", "-a", -1},
		{"empty string", "", ErrorLiteral},
		{"too long", "ABC", ErrorLiteral},
		{"invalid char", "1", ErrorLiteral},
		{"invalid char with minus", "-1", ErrorLiteral},
		{"space", " ", ErrorLiteral},
		{"special char", "@", ErrorLiteral},
		{"multiple minuses", "--A", ErrorLiteral},
		{"minus at end", "A-", ErrorLiteral},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Str2Lit(tt.input)
			if result != tt.expected {
				t.Errorf("Str2Lit(%q) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestClauseString(t *testing.T) {
	tests := []struct {
		name     string
		clause   *Clause
		expected string
	}{
		{
			name:     "empty clause",
			clause:   New(),
			expected: "{}",
		},
		{
			name: "single positive literal",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			expected: "{A}",
		},
		{
			name: "single negative literal",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			expected: "{-A}",
		},
		{
			name: "multiple positive literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(3)
				c.Insert(2)
				return c
			}(),
			expected: "{A, B, C}",
		},
		{
			name: "multiple negative literals",
			clause: func() *Clause {
				c := New()
				c.Insert(-3)
				c.Insert(-1)
				c.Insert(-2)
				return c
			}(),
			expected: "{-A, -B, -C}",
		},
		{
			name: "mixed positive and negative literals",
			clause: func() *Clause {
				c := New()
				c.Insert(2)
				c.Insert(-1)
				c.Insert(3)
				c.Insert(-4)
				return c
			}(),
			expected: "{-A, B, C, -D}",
		},
		{
			name: "literals with gaps",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(26)
				c.Insert(13)
				return c
			}(),
			expected: "{A, M, Z}",
		},
		{
			name: "duplicate literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(1)
				c.Insert(2)
				c.Insert(2)
				return c
			}(),
			expected: "{A, B}",
		},
		{
			name: "contradictory literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-1)
				return c
			}(),
			expected: "{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.clause.String()
			if result != tt.expected {
				t.Errorf("String() = %q; want %q", result, tt.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *Clause
		expectError bool
	}{
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:  "single positive literal",
			input: "A",
			expected: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
		},
		{
			name:  "single negative literal",
			input: "-A",
			expected: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
		},
		{
			name:  "multiple positive literals",
			input: "A,B,C",
			expected: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
		},
		{
			name:  "multiple negative literals",
			input: "-A,-B,-C",
			expected: func() *Clause {
				c := New()
				c.Insert(-1)
				c.Insert(-2)
				c.Insert(-3)
				return c
			}(),
		},
		{
			name:  "mixed positive and negative literals",
			input: "A,-B,C,-D",
			expected: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-2)
				c.Insert(3)
				c.Insert(-4)
				return c
			}(),
		},
		{
			name:  "literals with spaces",
			input: " A , B , -C ",
			expected: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(-3)
				return c
			}(),
		},
		{
			name:  "lowercase literals",
			input: "a,b,c",
			expected: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
		},
		{
			name:  "mixed case literals",
			input: "A,b,-C,d",
			expected: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(-3)
				c.Insert(4)
				return c
			}(),
		},
		{
			name:        "invalid format - number",
			input:       "1",
			expectError: true,
		},
		{
			name:        "invalid format - special char",
			input:       "@",
			expectError: true,
		},
		{
			name:        "invalid format - space",
			input:       " ",
			expectError: true,
		},
		{
			name:        "invalid format - multiple minuses",
			input:       "--A",
			expectError: true,
		},
		{
			name:        "invalid format - minus at end",
			input:       "A-",
			expectError: true,
		},
		{
			name:        "invalid format - empty component",
			input:       "A,,B",
			expectError: true,
		},
		{
			name:        "invalid format - too long component",
			input:       "ABC",
			expectError: true,
		},
		{
			name:  "contradictory literals",
			input: "A,-A,B",
			expected: func() *Clause {
				c := New()
				c.Insert(2)
				return c
			}(),
		},
		{
			name:  "duplicate literals",
			input: "A,A,B,B",
			expected: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("Parse(%q) expected error, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
				return
			}
			if !result.Equals(*tt.expected) {
				t.Errorf("Parse(%q) = %q; want %q", tt.input, result.String(), tt.expected.String())
			}
		})
	}
}

func TestClauseEquals(t *testing.T) {
	tests := []struct {
		name     string
		clause1  *Clause
		clause2  *Clause
		expected bool
	}{
		{
			name:     "both empty clauses",
			clause1:  New(),
			clause2:  New(),
			expected: true,
		},
		{
			name: "same single positive literal",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			expected: true,
		},
		{
			name: "same single negative literal",
			clause1: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			expected: true,
		},
		{
			name: "same multiple literals",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(-3)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(2)
				c.Insert(-3)
				c.Insert(1)
				return c
			}(),
			expected: true,
		},
		{
			name: "different single literals",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(2)
				return c
			}(),
			expected: false,
		},
		{
			name: "different sizes",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
			expected: false,
		},
		{
			name: "same literals different signs",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			expected: false,
		},
		{
			name: "one clause has extra literal",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			expected: false,
		},
		{
			name: "clauses with same literals but different insertion order",
			clause1: func() *Clause {
				c := New()
				c.Insert(3)
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(2)
				c.Insert(3)
				c.Insert(1)
				return c
			}(),
			expected: true,
		},
		{
			name: "clauses with contradictory literals",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-1)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				return c
			}(),
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.clause1.Equals(*tt.clause2)
			if result != tt.expected {
				t.Errorf("Equals() = %v; want %v. Clause1: %s, Clause2: %s",
					result, tt.expected, tt.clause1.String(), tt.clause2.String())
			}
		})
	}
}

func TestClauseCopy(t *testing.T) {
	tests := []struct {
		name     string
		original *Clause
		expected *Clause
	}{
		{
			name:     "empty clause",
			original: New(),
			expected: New(),
		},
		{
			name: "single positive literal",
			original: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			expected: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
		},
		{
			name: "single negative literal",
			original: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			expected: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
		},
		{
			name: "multiple positive literals",
			original: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(3)
				c.Insert(2)
				return c
			}(),
			expected: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
		},
		{
			name: "multiple negative literals",
			original: func() *Clause {
				c := New()
				c.Insert(-3)
				c.Insert(-1)
				c.Insert(-2)
				return c
			}(),
			expected: func() *Clause {
				c := New()
				c.Insert(-1)
				c.Insert(-2)
				c.Insert(-3)
				return c
			}(),
		},
		{
			name: "mixed positive and negative literals",
			original: func() *Clause {
				c := New()
				c.Insert(2)
				c.Insert(-1)
				c.Insert(3)
				c.Insert(-4)
				return c
			}(),
			expected: func() *Clause {
				c := New()
				c.Insert(-1)
				c.Insert(2)
				c.Insert(3)
				c.Insert(-4)
				return c
			}(),
		},
		{
			name: "literals with gaps",
			original: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(26)
				c.Insert(13)
				return c
			}(),
			expected: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(13)
				c.Insert(26)
				return c
			}(),
		},
		{
			name: "duplicate literals",
			original: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(1)
				c.Insert(2)
				c.Insert(2)
				return c
			}(),
			expected: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
		},
		{
			name: "contradictory literals",
			original: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-1)
				return c
			}(),
			expected: New(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that copy is equal to original
			copy := tt.original.Copy()
			if !copy.Equals(*tt.expected) {
				t.Errorf("Copy() = %q; want %q", copy.String(), tt.expected.String())
			}
			// Test that copy is independent from original
			if !copy.Equals(*tt.original) {
				t.Errorf("Copy() should equal original: %q != %q", copy.String(), tt.original.String())
			}
			// Test that modifying copy doesn't affect original
			originalSize := tt.original.Size()
			copy.Insert(100) // Insert a literal that doesn't exist in original
			if tt.original.Size() != originalSize {
				t.Errorf("Copy() should be independent from original. Original size changed from %d to %d",
					originalSize, tt.original.Size())
			}
		})
	}
}

func TestClauseSize(t *testing.T) {
	tests := []struct {
		name     string
		clause   *Clause
		expected int
	}{
		{
			name:     "empty clause",
			clause:   New(),
			expected: 0,
		},
		{
			name: "single positive literal",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			expected: 1,
		},
		{
			name: "single negative literal",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			expected: 1,
		},
		{
			name: "multiple positive literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
			expected: 3,
		},
		{
			name: "multiple negative literals",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				c.Insert(-2)
				c.Insert(-3)
				return c
			}(),
			expected: 3,
		},
		{
			name: "mixed positive and negative literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-2)
				c.Insert(3)
				c.Insert(-4)
				return c
			}(),
			expected: 4,
		},
		{
			name: "duplicate literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(1)
				c.Insert(2)
				c.Insert(2)
				return c
			}(),
			expected: 2,
		},
		{
			name: "contradictory literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-1)
				return c
			}(),
			expected: 0,
		},
		{
			name: "literals with gaps",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(10)
				c.Insert(20)
				return c
			}(),
			expected: 3,
		},
		{
			name: "large number of literals",
			clause: func() *Clause {
				c := New()
				for i := 1; i <= 26; i++ {
					c.Insert(Literal(i))
				}
				return c
			}(),
			expected: 26,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.clause.Size()
			if result != tt.expected {
				t.Errorf("Size() = %d; want %d. Clause: %s",
					result, tt.expected, tt.clause.String())
			}
		})
	}
}

func TestClauseIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		clause   *Clause
		expected bool
	}{
		{
			name:     "empty clause",
			clause:   New(),
			expected: true,
		},
		{
			name: "single positive literal",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			expected: false,
		},
		{
			name: "single negative literal",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			expected: false,
		},
		{
			name: "multiple positive literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
			expected: false,
		},
		{
			name: "multiple negative literals",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				c.Insert(-2)
				c.Insert(-3)
				return c
			}(),
			expected: false,
		},
		{
			name: "mixed positive and negative literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-2)
				c.Insert(3)
				c.Insert(-4)
				return c
			}(),
			expected: false,
		},
		{
			name: "duplicate literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(1)
				c.Insert(2)
				c.Insert(2)
				return c
			}(),
			expected: false,
		},
		{
			name: "contradictory literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-1)
				return c
			}(),
			expected: true,
		},
		{
			name: "clause that becomes empty after adding contradictory literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(-1)
				c.Insert(-2)
				return c
			}(),
			expected: true,
		},
		{
			name: "clause with many literals that becomes empty",
			clause: func() *Clause {
				c := New()
				for i := 1; i <= 10; i++ {
					c.Insert(Literal(i))
					c.Insert(Literal(-i))
				}
				return c
			}(),
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.clause.IsEmpty()
			if result != tt.expected {
				t.Errorf("IsEmpty() = %v; want %v. Clause: %s",
					result, tt.expected, tt.clause.String())
			}
		})
	}
}

func TestClauseContains(t *testing.T) {
	tests := []struct {
		name     string
		clause   *Clause
		literal  Literal
		expected bool
	}{
		{
			name:     "empty clause",
			clause:   New(),
			literal:  1,
			expected: false,
		},
		{
			name: "single positive literal - contains",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			literal:  1,
			expected: true,
		},
		{
			name: "single positive literal - does not contain",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			literal:  2,
			expected: false,
		},
		{
			name: "single negative literal - contains",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			literal:  -1,
			expected: true,
		},
		{
			name: "single negative literal - does not contain",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			literal:  -2,
			expected: false,
		},
		{
			name: "multiple positive literals - contains first",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
			literal:  1,
			expected: true,
		},
		{
			name: "multiple positive literals - contains middle",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
			literal:  2,
			expected: true,
		},
		{
			name: "multiple positive literals - contains last",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
			literal:  3,
			expected: true,
		},
		{
			name: "multiple positive literals - does not contain",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
			literal:  4,
			expected: false,
		},
		{
			name: "multiple negative literals - contains",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				c.Insert(-2)
				c.Insert(-3)
				return c
			}(),
			literal:  -2,
			expected: true,
		},
		{
			name: "multiple negative literals - does not contain",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				c.Insert(-2)
				c.Insert(-3)
				return c
			}(),
			literal:  -4,
			expected: false,
		},
		{
			name: "mixed positive and negative literals - contains positive",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-2)
				c.Insert(3)
				c.Insert(-4)
				return c
			}(),
			literal:  3,
			expected: true,
		},
		{
			name: "mixed positive and negative literals - contains negative",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-2)
				c.Insert(3)
				c.Insert(-4)
				return c
			}(),
			literal:  -4,
			expected: true,
		},
		{
			name: "mixed positive and negative literals - does not contain",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-2)
				c.Insert(3)
				c.Insert(-4)
				return c
			}(),
			literal:  5,
			expected: false,
		},
		{
			name: "contradictory literals - contains positive",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-1)
				return c
			}(),
			literal:  1,
			expected: false, // Clause becomes empty when contradictory literals are added
		},
		{
			name: "contradictory literals - contains negative",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-1)
				return c
			}(),
			literal:  -1,
			expected: false, // Clause becomes empty when contradictory literals are added
		},
		{
			name: "literals with gaps - contains",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(10)
				c.Insert(20)
				return c
			}(),
			literal:  10,
			expected: true,
		},
		{
			name: "literals with gaps - does not contain",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(10)
				c.Insert(20)
				return c
			}(),
			literal:  15,
			expected: false,
		},
		{
			name: "check for opposite literal - positive",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			literal:  -1,
			expected: false,
		},
		{
			name: "check for opposite literal - negative",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			literal:  1,
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.clause.Contains(tt.literal)
			if result != tt.expected {
				t.Errorf("Contains(%d) = %v; want %v. Clause: %s",
					tt.literal, result, tt.expected, tt.clause.String())
			}
		})
	}
}

func TestClauseInsert(t *testing.T) {
	tests := []struct {
		name           string
		clause         *Clause
		literal        Literal
		expected       bool
		expectedClause *Clause
	}{
		{
			name:     "insert into empty clause",
			clause:   New(),
			literal:  1,
			expected: true,
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
		},
		{
			name: "insert new positive literal",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			literal:  2,
			expected: true,
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
		},
		{
			name: "insert duplicate positive literal",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			literal:  1,
			expected: false,
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
		},
		{
			name: "insert new negative literal",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			literal:  -2,
			expected: true,
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-2)
				return c
			}(),
		},
		{
			name: "insert duplicate negative literal",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			literal:  -1,
			expected: false,
			expectedClause: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
		},
		{
			name: "insert negation of existing positive literal",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			literal:        -1,
			expected:       false,
			expectedClause: New(), // Clause becomes empty
		},
		{
			name: "insert negation of existing negative literal",
			clause: func() *Clause {
				c := New()
				c.Insert(-1)
				return c
			}(),
			literal:        1,
			expected:       false,
			expectedClause: New(), // Clause becomes empty
		},
		{
			name: "insert into clause with multiple literals",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(-3)
				return c
			}(),
			literal:  4,
			expected: true,
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(-3)
				c.Insert(4)
				return c
			}(),
		},
		{
			name: "insert literal that makes clause empty",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
			literal:  -1,
			expected: false,
			expectedClause: func() *Clause {
				c := New()
				c.Insert(2)
				return c
			}(),
		},
		{
			name: "insert literal with gaps",
			clause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(10)
				return c
			}(),
			literal:  5,
			expected: true,
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(5)
				c.Insert(10)
				return c
			}(),
		},
		{
			name:     "insert large positive literal",
			clause:   New(),
			literal:  26, // Z
			expected: true,
			expectedClause: func() *Clause {
				c := New()
				c.Insert(26)
				return c
			}(),
		},
		{
			name:     "insert large negative literal",
			clause:   New(),
			literal:  -26, // -Z
			expected: true,
			expectedClause: func() *Clause {
				c := New()
				c.Insert(-26)
				return c
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the return value
			result := tt.clause.Insert(tt.literal)
			if result != tt.expected {
				t.Errorf("Insert(%d) return value = %v; want %v. Clause: %s",
					tt.literal, result, tt.expected, tt.clause.String())
			}
			// Test the clause state after insertion
			if !tt.clause.Equals(*tt.expectedClause) {
				t.Errorf("Insert(%d) clause state = %s; want %s",
					tt.literal, tt.clause.String(), tt.expectedClause.String())
			}
		})
	}
}

func TestClauseResolve(t *testing.T) {
	tests := []struct {
		name           string
		clause1        *Clause
		clause2        *Clause
		expectedClause *Clause
		expectedFound  bool
	}{
		{
			name:           "resolve with empty clause",
			clause1:        New(),
			clause2:        New(),
			expectedClause: New(),
			expectedFound:  false,
		},
		{
			name: "resolve identical single literals",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			expectedFound: false,
		},
		{
			name: "resolve different single literals",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(2)
				return c
			}(),
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
			expectedFound: false,
		},
		{
			name: "resolve with one common literal",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
			expectedFound: false,
		},
		{
			name: "resolve with multiple common literals",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(2)
				c.Insert(3)
				c.Insert(4)
				return c
			}(),
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				c.Insert(4)
				return c
			}(),
			expectedFound: false,
		},
		{
			name: "resolve with no common literals",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(3)
				c.Insert(4)
				return c
			}(),
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				c.Insert(3)
				c.Insert(4)
				return c
			}(),
			expectedFound: false,
		},
		{
			name: "resolve with contradictory literals",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(-1)
				c.Insert(3)
				return c
			}(),
			expectedClause: func() *Clause {
				c := New()
				c.Insert(2)
				c.Insert(3)
				return c
			}(),
			expectedFound: true,
		},
		{
			name: "resolve with negative literals",
			clause1: func() *Clause {
				c := New()
				c.Insert(-1)
				c.Insert(-2)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(-2)
				c.Insert(-3)
				return c
			}(),
			expectedClause: func() *Clause {
				c := New()
				c.Insert(-1)
				c.Insert(-2)
				c.Insert(-3)
				return c
			}(),
			expectedFound: false,
		},
		{
			name: "resolve mixed positive and negative literals",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-2)
				c.Insert(3)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(-2)
				c.Insert(4)
				c.Insert(-5)
				return c
			}(),
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(-2)
				c.Insert(3)
				c.Insert(4)
				c.Insert(-5)
				return c
			}(),
			expectedFound: false,
		},
		{
			name: "resolve where one clause is subset of another",
			clause1: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(1)
				return c
			}(),
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(2)
				return c
			}(),
			expectedFound: false,
		},
		{
			name: "resolve with large literals",
			clause1: func() *Clause {
				c := New()
				c.Insert(25)
				c.Insert(26)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(26)
				c.Insert(1)
				return c
			}(),
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(25)
				c.Insert(26)
				return c
			}(),
			expectedFound: false,
		},
		{
			name: "resolve with success",
			clause1: func() *Clause {
				c := New()
				c.Insert(25)
				c.Insert(-26)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(26)
				c.Insert(1)
				return c
			}(),
			expectedClause: func() *Clause {
				c := New()
				c.Insert(1)
				c.Insert(25)
				return c
			}(),
			expectedFound: true,
		},
		{
			name: "resolve with no success",
			clause1: func() *Clause {
				c := New()
				c.Insert(-25)
				c.Insert(-26)
				return c
			}(),
			clause2: func() *Clause {
				c := New()
				c.Insert(26)
				c.Insert(25)
				return c
			}(),
			expectedClause: nil,
			expectedFound:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found := tt.clause1.Resolve(*tt.clause2)
			// Check if result is nil
			if tt.expectedClause == nil {
				if result != nil {
					t.Errorf("Resolve() result = %v; want nil", result)
				}
				return
			}
			// Check if result is not nil when expected
			if result == nil {
				t.Errorf("Resolve() result = nil; want %v", tt.expectedClause)
				return
			}
			// Check the found flag
			if found != tt.expectedFound {
				t.Errorf("Resolve() found = %v; want %v", found, tt.expectedFound)
			}
			// Check the resulting clause
			if !result.Equals(*tt.expectedClause) {
				t.Errorf("Resolve() result = %s; want %s",
					result.String(), tt.expectedClause.String())
			}
		})
	}
}

func TestRes(t *testing.T) {
	tests := []struct {
		name     string
		clauses  []Clause
		expected bool
	}{
		{
			name:     "empty clause set",
			clauses:  []Clause{},
			expected: false,
		},
		{
			name: "single clause",
			clauses: []Clause{
				*func() *Clause {
					c := New()
					c.Insert(1)
					c.Insert(2)
					return c
				}(),
			},
			expected: false,
		},
		{
			name: "unsatisfiable - simple contradiction",
			clauses: []Clause{
				*func() *Clause {
					c := New()
					c.Insert(1)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(-1)
					return c
				}(),
			},
			expected: true,
		},
		{
			name: "unsatisfiable - requires resolution",
			clauses: []Clause{
				*func() *Clause {
					c := New()
					c.Insert(1)
					c.Insert(2)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(-1)
					c.Insert(2)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(1)
					c.Insert(-2)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(-1)
					c.Insert(-2)
					return c
				}(),
			},
			expected: true,
		},
		{
			name: "satisfiable - no contradiction",
			clauses: []Clause{
				*func() *Clause {
					c := New()
					c.Insert(1)
					c.Insert(2)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(-1)
					c.Insert(3)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(2)
					c.Insert(-3)
					return c
				}(),
			},
			expected: false,
		},
		{
			name: "already contains empty clause",
			clauses: []Clause{
				*New(),
				*func() *Clause {
					c := New()
					c.Insert(1)
					return c
				}(),
			},
			expected: true,
		},
		{
			name: "unit propagation leads to empty clause",
			clauses: []Clause{
				*func() *Clause {
					c := New()
					c.Insert(1)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(-1)
					c.Insert(2)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(-1)
					c.Insert(-2)
					return c
				}(),
			},
			expected: true,
		},
		{
			name: "complex satisfiable case",
			clauses: []Clause{
				*func() *Clause {
					c := New()
					c.Insert(1)
					c.Insert(2)
					c.Insert(3)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(-1)
					c.Insert(-2)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(2)
					c.Insert(-3)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(-2)
					c.Insert(3)
					return c
				}(),
			},
			expected: false,
		},
		{
			name: "all clauses already resolved",
			clauses: []Clause{
				*func() *Clause {
					c := New()
					c.Insert(1)
					return c
				}(),
				*func() *Clause {
					c := New()
					c.Insert(2)
					return c
				}(),
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Res(tt.clauses, 0)
			if result != tt.expected {
				t.Errorf("Res() = %v; want %v for clauses %s",
					result, tt.expected, formatClauses(tt.clauses))
			}
		})
	}
}

// Helper function to format clauses for better error messages
func formatClauses(clauses []Clause) string {
	result := "["
	for i, c := range clauses {
		if i > 0 {
			result += ", "
		}
		result += c.String()
	}
	return result + "]"
}

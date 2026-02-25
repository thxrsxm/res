# Resolution Theorem Prover

A command-line tool for checking the satisfiability of propositional logic clauses using the resolution method.

## Overview

This tool implements a resolution-based theorem prover for propositional logic. It takes a set of clauses as input and determines whether they are satisfiable (can all be true simultaneously) or unsatisfiable (contain a contradiction).

## Installation

### Prerequisites

- Go 1.16 or later

### Building from Source

```bash
git clone https://github.com/thxrsxm/res.git
cd res
make build
```
The compiled binary will be placed in the `./bin` directory.

**OR**

```bash
git clone https://github.com/thxrsxm/res.git
cd res
go install
```

## Usage

```
res [options] <clause1> <clause2> ...
```

### Arguments

- `<clause>`: A clause in the format `A,B,-C` (comma-separated literals)
  - Each literal is a single letter (A-Z)
  - Optionally prefixed with `-` for negation
  - Example: `A,-B,C` means "A AND NOT B AND C"

### Output

- `[ ]`: The clause set is unsatisfiable (contradiction found)
- `[x]`: The clause set is satisfiable (no contradiction found)

### Examples

1. Simple contradiction:
   ```bash
   res A,-A
   ```
   Output: `[ ]`

2. Multiple clauses:
   ```bash
   res "A,B" "-A,C" "-B,C" "-C"
   ```
   Output: `[ ]`

3. Satisfiable set:
   ```bash
   res A,B,-C -A,B,C -B,C
   ```
   Output: `[x]`

4. Using `--` to handle clauses starting with `-`:
   ```bash
   res -- -A,B -B,C A
   ```

## How It Works

The tool implements the resolution method from propositional logic:

1. Parses input clauses into an internal representation
2. Applies the resolution rule to derive new clauses
3. Checks for the empty clause (which indicates unsatisfiability)
4. Continues until either:
   - The empty clause is derived (unsatisfiable)
   - No new clauses can be derived (satisfiable)

## Clause Format

- Literals: Single uppercase or lowercase letters (A-Z and a-z)
- Negative literals: Prefix with `-` (e.g., `-A` for "NOT A")
- Clauses: Comma-separated literals (e.g., `A,-B,C`)
- Multiple clauses: Space-separated

## Building

The Makefile provides several targets:

- `make build`: Build the application
- `make export`: Build for multiple platforms (Linux, macOS, Windows)
- `make run`: Run the application
- `make clean`: Clean up build artifacts

## Testing

Run the test suite with:

```bash
go test ./...
```

## License

MIT License

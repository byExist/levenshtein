

# levenshtein [![GoDoc](https://pkg.go.dev/badge/github.com/byExist/levenshtein.svg)](https://pkg.go.dev/github.com/byExist/levenshtein) [![Go Report Card](https://goreportcard.com/badge/github.com/byExist/levenshtein)](https://goreportcard.com/report/github.com/byExist/levenshtein)

## What is levenshtein?

levenshtein is a Go package implementing the Levenshtein distance algorithm with customizable cost functions for insertion, deletion, and substitution operations. It’s lightweight and flexible, suitable for practical string distance computations.

## Features

- Custom cost functions for insert, delete, and replace operations
- Combine multiple cost functions using min, max, or average strategies
- Interface-based design for extensibility
- Supports Unicode, asymmetric lengths, and empty strings

## Installation

```bash
go get github.com/byExist/levenshtein
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/byExist/levenshtein"
	"unicode"
)

func main() {
	// Basic Levenshtein distance
	l := levenshtein.New()
	fmt.Println("Distance between 'kitten' and 'sitting':", l.Distance("kitten", "sitting"))

	// Custom cost functions
	opts := []levenshtein.Option{
		levenshtein.WithInsertCost(func(r rune) float64 {
			if unicode.IsSpace(r) {
				return 0.5
			}
			return 1
		}),
		levenshtein.WithDeleteCost(func(r rune) float64 {
			if unicode.IsPunct(r) {
				return 0.5
			}
			return 1
		}),
		levenshtein.WithReplaceCost(func(a, b rune) float64 {
			if a == b {
				return 0
			}
			if unicode.ToLower(a) == unicode.ToLower(b) {
				return 0.5
			}
			return 1
		}),
	}
	lev := levenshtein.New(opts...)
	fmt.Println("Custom distance:", lev.Distance("Hello, World!", "hello, world!"))
}
```

## Usage

The levenshtein package lets you flexibly calculate edit distances between strings. You can define your own cost logic, combine multiple functions strategically, and extend behavior through the exposed interface. Ideal for fuzzy string matching, text correction, and more.

## API Overview

### Constructors

- `New(options ...Option) Levenshtein`  

### Methods

- `Distance(a, b string) float64`

### Options

- `WithInsertCost(InsertCost) Option`  
- `WithDeleteCost(DeleteCost) Option`  
- `WithReplaceCost(ReplaceCost) Option`  

### Composition

- `ComposeInsertCost(strategy, funcs...)`  
- `ComposeDeleteCost(strategy, funcs...)`  
- `ComposeReplaceCost(strategy, funcs...)`  
- `ComposeWeightedInsertCost([]WeightedInsert) (InsertCost, error)`  
- `ComposeWeightedDeleteCost([]WeightedDelete) (DeleteCost, error)`  
- `ComposeWeightedReplaceCost([]WeightedReplace) (ReplaceCost, error)`  

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
// Package levenshtein provides an implementation of the Levenshtein distance algorithm with customizable cost functions.
package levenshtein

import (
	"fmt"
	"math"
)

// InsertCost defines the cost function for inserting a rune.
// DeleteCost defines the cost function for deleting a rune.
// ReplaceCost defines the cost function for replacing one rune with another.
type (
	InsertCost  func(rune) float64
	DeleteCost  func(rune) float64
	ReplaceCost func(rune, rune) float64
)

// Levenshtein defines the interface for computing edit distance between two strings.
type Levenshtein interface {
	Distance(a, b string) float64
}

// levenshtein calculates the edit distance between two strings using configurable cost functions.
type levenshtein struct {
	insertCost  InsertCost
	deleteCost  DeleteCost
	replaceCost ReplaceCost
}

// Option represents a configuration function for customizing a Levenshtein instance.
type Option func(*levenshtein)

// WithInsertCost sets the insert cost function.
func WithInsertCost(c InsertCost) Option {
	return func(lev *levenshtein) {
		lev.insertCost = c
	}
}

// WithDeleteCost sets the delete cost function.
func WithDeleteCost(c DeleteCost) Option {
	return func(lev *levenshtein) {
		lev.deleteCost = c
	}
}

// WithReplaceCost sets the replace cost function.
func WithReplaceCost(c ReplaceCost) Option {
	return func(lev *levenshtein) {
		lev.replaceCost = c
	}
}

// ComposeStrategy determines how to combine multiple cost functions into one.
type ComposeStrategy int

// StrategyMin uses the minimum value among provided cost functions.
// StrategyMax uses the maximum value among provided cost functions.
// StrategyAvg uses the average value among provided cost functions.
const (
	StrategyMin ComposeStrategy = iota
	StrategyMax
	StrategyAvg
)

// ComposeInsertCost combines multiple InsertCost functions using the specified strategy.
func ComposeInsertCost(strategy ComposeStrategy, funcs ...InsertCost) (InsertCost, error) {
	if len(funcs) == 0 {
		return nil, fmt.Errorf("no insert cost function provided")
	}
	if len(funcs) == 1 {
		return funcs[0], nil
	}
	return func(r rune) float64 {
		var result float64
		switch strategy {
		case StrategyMin:
			result = math.MaxFloat64
			for _, f := range funcs {
				if cost := f(r); cost < result {
					result = cost
				}
			}
		case StrategyMax:
			result = 0
			for _, f := range funcs {
				if cost := f(r); cost > result {
					result = cost
				}
			}
		case StrategyAvg:
			sum := float64(0)
			for _, f := range funcs {
				sum += f(r)
			}
			result = sum / float64(len(funcs))
		default:
			panic("unknown strategy")
		}
		return result
	}, nil
}

// ComposeDeleteCost combines multiple DeleteCost functions using the specified strategy.
func ComposeDeleteCost(strategy ComposeStrategy, funcs ...DeleteCost) (DeleteCost, error) {
	if len(funcs) == 0 {
		return nil, fmt.Errorf("no delete cost function provided")
	}
	if len(funcs) == 1 {
		return funcs[0], nil
	}
	return func(r rune) float64 {
		var result float64
		switch strategy {
		case StrategyMin:
			result = math.MaxFloat64
			for _, f := range funcs {
				if cost := f(r); cost < result {
					result = cost
				}
			}
		case StrategyMax:
			result = 0
			for _, f := range funcs {
				if cost := f(r); cost > result {
					result = cost
				}
			}
		case StrategyAvg:
			sum := float64(0)
			for _, f := range funcs {
				sum += f(r)
			}
			result = sum / float64(len(funcs))
		default:
			panic("unknown strategy")
		}
		return result
	}, nil
}

// ComposeReplaceCost combines multiple ReplaceCost functions using the specified strategy.
func ComposeReplaceCost(strategy ComposeStrategy, funcs ...ReplaceCost) (ReplaceCost, error) {
	if len(funcs) == 0 {
		return nil, fmt.Errorf("no replace cost function provided")
	}
	if len(funcs) == 1 {
		return funcs[0], nil
	}
	return func(a, b rune) float64 {
		var result float64
		switch strategy {
		case StrategyMin:
			result = math.MaxFloat64
			for _, f := range funcs {
				if cost := f(a, b); cost < result {
					result = cost
				}
			}
		case StrategyMax:
			result = 0
			for _, f := range funcs {
				if cost := f(a, b); cost > result {
					result = cost
				}
			}
		case StrategyAvg:
			sum := float64(0)
			for _, f := range funcs {
				sum += f(a, b)
			}
			result = sum / float64(len(funcs))
		default:
			panic("unknown strategy")
		}
		return result
	}, nil
}

// DefaultInsertCost returns a constant insert cost of 1.
func DefaultInsertCost(r rune) float64 { return 1 }

// DefaultDeleteCost returns a constant delete cost of 1.
func DefaultDeleteCost(r rune) float64 { return 1 }

// DefaultReplaceCost returns 0 if runes are equal, otherwise 1.
func DefaultReplaceCost(a, b rune) float64 {
	if a == b {
		return 0
	}
	return 1
}

// New creates a new Levenshtein instance with the specified options.
func New(options ...Option) Levenshtein {
	lev := &levenshtein{
		insertCost:  DefaultInsertCost,
		deleteCost:  DefaultDeleteCost,
		replaceCost: DefaultReplaceCost,
	}
	for _, opt := range options {
		opt(lev)
	}
	return lev
}

// Distance calculates the Levenshtein distance between two strings using the configured cost functions.
func (lev *levenshtein) Distance(a, b string) float64 {
	runesA := []rune(a)
	runesB := []rune(b)
	lenA := len(runesA)
	lenB := len(runesB)

	if lenA == 0 {
		sum := 0.0
		for _, r := range runesB {
			sum += lev.insertCost(r)
		}
		return sum
	}
	if lenB == 0 {
		sum := 0.0
		for _, r := range runesA {
			sum += lev.deleteCost(r)
		}
		return sum
	}

	dpRow := make([]float64, lenB+1)
	dpRow[0] = 0
	for j, r := range runesB {
		dpRow[j+1] = dpRow[j] + lev.insertCost(r)
	}

	for i := 1; i <= lenA; i++ {
		prevCost := dpRow[0]
		dpRow[0] += lev.deleteCost(runesA[i-1])
		for j := 1; j <= lenB; j++ {
			temp := dpRow[j]
			insert := dpRow[j-1] + lev.insertCost(runesB[j-1])
			delete := dpRow[j] + lev.deleteCost(runesA[i-1])
			replace := prevCost + lev.replaceCost(runesA[i-1], runesB[j-1])
			dpRow[j] = min(insert, delete, replace)
			prevCost = temp
		}
	}

	return dpRow[lenB]
}

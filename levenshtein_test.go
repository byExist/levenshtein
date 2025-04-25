package levenshtein_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/byExist/levenshtein"
	"github.com/stretchr/testify/assert"
)

func TestDistanceWithCustomCosts(t *testing.T) {
	insertCost := func(r rune) float64 { return 1.1 }
	deleteCost := func(r rune) float64 { return 1.2 }
	replaceCost := func(a, b rune) float64 {
		if a == b {
			return 0
		}
		return 2.2
	}
	l := levenshtein.New(
		levenshtein.WithInsertCost(insertCost),
		levenshtein.WithDeleteCost(deleteCost),
		levenshtein.WithReplaceCost(replaceCost),
	)
	d := l.Distance("abc", "xbc")
	expected := 2.2 // 'a' -> 'x'
	assert.Equal(t, expected, d)
}

func TestComposeDeleteCostVariousStrategies(t *testing.T) {
	f1 := func(r rune) float64 { return 4 }
	f2 := func(r rune) float64 { return 2 }
	tests := []struct {
		name     string
		strategy levenshtein.ComposeStrategy
		expected float64
	}{
		{"Min", levenshtein.StrategyMin, 2},
		{"Max", levenshtein.StrategyMax, 4},
		{"Avg", levenshtein.StrategyAvg, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost, err := levenshtein.ComposeDeleteCost(tt.strategy, f1, f2)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, cost('a'))
		})
	}
}

func TestComposeInsertCostSingleFunction(t *testing.T) {
	single := func(r rune) float64 { return 7 }
	cost, err := levenshtein.ComposeInsertCost(levenshtein.StrategyMin, single)
	assert.NoError(t, err)
	assert.Equal(t, 7.0, cost('a'))
}

func TestComposeInsertCostVariousStrategies(t *testing.T) {
	f1 := func(r rune) float64 { return 3 }
	f2 := func(r rune) float64 { return 1 }
	tests := []struct {
		name     string
		strategy levenshtein.ComposeStrategy
		expected float64
	}{
		{"Min", levenshtein.StrategyMin, 1},
		{"Max", levenshtein.StrategyMax, 3},
		{"Avg", levenshtein.StrategyAvg, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost, err := levenshtein.ComposeInsertCost(tt.strategy, f1, f2)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, cost('a'))
		})
	}
}

func TestComposeCostFunctionsNoInput(t *testing.T) {
	_, err := levenshtein.ComposeInsertCost(levenshtein.StrategyMin)
	assert.Error(t, err)

	_, err = levenshtein.ComposeDeleteCost(levenshtein.StrategyMin)
	assert.Error(t, err)

	_, err = levenshtein.ComposeReplaceCost(levenshtein.StrategyMin)
	assert.Error(t, err)
}

func TestComposeReplaceCostVariousStrategies(t *testing.T) {
	f1 := func(a, b rune) float64 { return 4 }
	f2 := func(a, b rune) float64 { return 2 }
	tests := []struct {
		name     string
		strategy levenshtein.ComposeStrategy
		expected float64
	}{
		{"Min", levenshtein.StrategyMin, 2},
		{"Max", levenshtein.StrategyMax, 4},
		{"Avg", levenshtein.StrategyAvg, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost, err := levenshtein.ComposeReplaceCost(tt.strategy, f1, f2)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, cost('a', 'b'))
		})
	}
}

func TestDistanceWithCustomDeleteCost(t *testing.T) {
	deleteCost := func(r rune) float64 {
		return 2
	}
	l := levenshtein.New(levenshtein.WithDeleteCost(deleteCost))
	if d := l.Distance("abc", ""); d != 6 {
		t.Errorf("expected 6, got %v", d)
	}
}

func TestDistanceWithCustomInsertCost(t *testing.T) {
	insertCost := func(r rune) float64 {
		return 2
	}
	l := levenshtein.New(levenshtein.WithInsertCost(insertCost))
	if d := l.Distance("", "abc"); d != 6 {
		t.Errorf("expected 6, got %v", d)
	}
}

func TestDistanceWithCustomReplaceCost(t *testing.T) {
	replaceCost := func(a, b rune) float64 {
		return 2
	}
	l := levenshtein.New(levenshtein.WithReplaceCost(replaceCost))
	if d := l.Distance("abc", "xyz"); d != 3*2 {
		t.Errorf("expected 6, got %v", d)
	}
}

func TestDistanceBasicCases(t *testing.T) {
	l := levenshtein.New()
	tests := []struct {
		name     string
		a, b     string
		expected float64
	}{
		{"EmptyStrings", "", "", 0},
		{"InsertOnly", "", "abc", 3},
		{"DeleteOnly", "abc", "", 3},
		{"SameStrings", "abc", "abc", 0},
		{"AllReplace", "aaa", "bbb", 3},
		{"TypicalExample", "kitten", "sitting", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if d := l.Distance(tt.a, tt.b); d != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, d)
			}
		})
	}
}

func TestDistanceWithZeroReplaceCost(t *testing.T) {
	replaceCost := func(a, b rune) float64 {
		if a == b {
			return 0
		}
		return 5
	}
	l := levenshtein.New(levenshtein.WithReplaceCost(replaceCost))
	if d := l.Distance("abc", "abc"); d != 0 {
		t.Errorf("expected 0, got %v", d)
	}
}

func TestDistanceWithUnicodeCharacters(t *testing.T) {
	l := levenshtein.New()
	d := l.Distance("안녕하세요", "안녕하세여")
	assert.Equal(t, 1.0, d)
}

func TestDistanceWithVeryLongStrings(t *testing.T) {
	a := strings.Repeat("a", 10000)
	b := strings.Repeat("a", 9999) + "b"
	l := levenshtein.New()
	d := l.Distance(a, b)
	assert.Equal(t, 1.0, d)
}

func BenchmarkDistanceAsymmetricLengths(b *testing.B) {
	a := "abc"
	bb := strings.Repeat("a", 1000)
	l := levenshtein.New()
	for i := 0; i < b.N; i++ {
		l.Distance(a, bb)
	}
}

func BenchmarkDistanceWithComposedCosts(b *testing.B) {
	f1 := func(r rune) float64 { return 1 }
	f2 := func(r rune) float64 { return 2 }
	rc1 := func(a, b rune) float64 { return 1 }
	rc2 := func(a, b rune) float64 { return 2 }

	insert, _ := levenshtein.ComposeInsertCost(levenshtein.StrategyAvg, f1, f2)
	delete, _ := levenshtein.ComposeDeleteCost(levenshtein.StrategyAvg, f1, f2)
	replace, _ := levenshtein.ComposeReplaceCost(levenshtein.StrategyAvg, rc1, rc2)

	l := levenshtein.New(
		levenshtein.WithInsertCost(insert),
		levenshtein.WithDeleteCost(delete),
		levenshtein.WithReplaceCost(replace),
	)
	a := strings.Repeat("a", 100)
	bb := strings.Repeat("b", 100)
	for i := 0; i < b.N; i++ {
		l.Distance(a, bb)
	}
}

func BenchmarkDistanceWithCustomInsertDeleteCosts(b *testing.B) {
	insert := func(r rune) float64 { return 1.5 }
	delete := func(r rune) float64 { return 1.2 }
	l := levenshtein.New(
		levenshtein.WithInsertCost(insert),
		levenshtein.WithDeleteCost(delete),
	)
	for i := 0; i < b.N; i++ {
		l.Distance("abcdefgh", "abcxefgh")
	}
}

func BenchmarkDistanceWithSimpleReplaceCost(b *testing.B) {
	replace := func(a, b rune) float64 {
		if a == b {
			return 0
		}
		return 2
	}
	l := levenshtein.New(levenshtein.WithReplaceCost(replace))
	for i := 0; i < b.N; i++ {
		l.Distance("abcdefgh", "abcdxfgh")
	}
}

func BenchmarkDistanceWithEmptyStrings(b *testing.B) {
	l := levenshtein.New()
	for i := 0; i < b.N; i++ {
		l.Distance("", "")
	}
}

func BenchmarkDistanceWithLongStrings(b *testing.B) {
	a := strings.Repeat("a", 1000)
	bb := strings.Repeat("b", 1000)
	l := levenshtein.New()
	for i := 0; i < b.N; i++ {
		l.Distance(a, bb)
	}
}

func BenchmarkDistanceWithMediumStrings(b *testing.B) {
	a := strings.Repeat("a", 200)
	bb := strings.Repeat("b", 200)
	l := levenshtein.New()
	for i := 0; i < b.N; i++ {
		l.Distance(a, bb)
	}
}

func BenchmarkDistanceWithShortDifferentStrings(b *testing.B) {
	l := levenshtein.New()
	for i := 0; i < b.N; i++ {
		l.Distance("abc", "xyz")
	}
}

func BenchmarkDistanceWithUnicodeRunes(b *testing.B) {
	a := "안녕하세요😊"
	bb := "안녕하세요😢"
	l := levenshtein.New()
	for i := 0; i < b.N; i++ {
		l.Distance(a, bb)
	}
}

// ExampleComposeDeleteCost_avg demonstrates combining multiple DeleteCost functions using the average strategy.
func ExampleComposeDeleteCost_avg() {
	f1 := func(r rune) float64 { return 2 }
	f2 := func(r rune) float64 { return 4 }
	cost, err := levenshtein.ComposeDeleteCost(levenshtein.StrategyAvg, f1, f2)
	if err != nil {
		panic(err)
	}
	fmt.Println(cost('x'))
	// Output: 3
}

// ExampleComposeInsertCost_min demonstrates combining multiple InsertCost functions using the minimum strategy.
func ExampleComposeInsertCost_min() {
	f1 := func(r rune) float64 { return 3 }
	f2 := func(r rune) float64 { return 1 }
	cost, err := levenshtein.ComposeInsertCost(levenshtein.StrategyMin, f1, f2)
	if err != nil {
		panic(err)
	}
	fmt.Println(cost('x'))
	// Output: 1
}

// ExampleComposeReplaceCost_max demonstrates combining multiple ReplaceCost functions using the maximum strategy.
func ExampleComposeReplaceCost_max() {
	f1 := func(a, b rune) float64 { return 2 }
	f2 := func(a, b rune) float64 { return 5 }
	cost, err := levenshtein.ComposeReplaceCost(levenshtein.StrategyMax, f1, f2)
	if err != nil {
		panic(err)
	}
	fmt.Println(cost('a', 'b'))
	// Output: 5
}

// ExampleNew_basic demonstrates basic usage of the Levenshtein distance with default costs.
func ExampleNew_basic() {
	l := levenshtein.New()
	fmt.Println(l.Distance("kitten", "sitting"))
	// Output: 3
}

// ExampleNew_withAllOptions demonstrates creating a Levenshtein instance
// with custom insert, delete, replace costs and a cutoff distance.
func ExampleNew_withAllOptions() {
	insertCost := func(r rune) float64 {
		if r == '!' {
			return 5
		}
		return 1
	}
	deleteCost := func(r rune) float64 {
		if r == '?' {
			return 5
		}
		return 1
	}
	replaceCost := func(a, b rune) float64 {
		if a == b {
			return 0
		}
		return 2
	}

	l := levenshtein.New(
		levenshtein.WithInsertCost(insertCost),
		levenshtein.WithDeleteCost(deleteCost),
		levenshtein.WithReplaceCost(replaceCost),
	)
	fmt.Println(l.Distance("hello!", "he?lo"))
	// Output: 3
}

// ExampleNew_withDeleteOnly demonstrates creating a Levenshtein instance with custom delete cost only.
func ExampleNew_withDeleteOnly() {
	deleteCost := func(r rune) float64 {
		if r == 'x' {
			return 5
		}
		return 1
	}
	l := levenshtein.New(levenshtein.WithDeleteCost(deleteCost))
	fmt.Println(l.Distance("ax", ""))
	// Output: 6
}

// ExampleNew_withInsertOnly demonstrates creating a Levenshtein instance with custom insert cost only.
func ExampleNew_withInsertOnly() {
	insertCost := func(r rune) float64 {
		if r == 'z' {
			return 10
		}
		return 1
	}
	l := levenshtein.New(levenshtein.WithInsertCost(insertCost))
	fmt.Println(l.Distance("", "az"))
	// Output: 11
}

// ExampleNew_withReplaceOnly demonstrates creating a Levenshtein instance with custom replace cost only.
func ExampleNew_withReplaceOnly() {
	replaceCost := func(a, b rune) float64 {
		if a == 'a' && b == 'b' {
			return 0.5
		}
		return 2
	}
	l := levenshtein.New(levenshtein.WithReplaceCost(replaceCost))
	fmt.Println(l.Distance("a", "b"))
	// Output: 0.5
}

// ExampleReplaceCost_customLogic demonstrates using a replace cost function
// that assigns lower cost for visually similar characters like 'O' and '0'.
func ExampleReplaceCost_customLogic() {
	replaceCost := func(a, b rune) float64 {
		if a == b {
			return 0
		}
		similar := map[rune]rune{'O': '0', '0': 'O', 'l': '1', '1': 'l'}
		if similar[a] == b {
			return 0.3 // visually similar characters have reduced cost
		}
		return 1
	}
	l := levenshtein.New(levenshtein.WithReplaceCost(replaceCost))
	fmt.Println(l.Distance("O1l", "01l"))
	// Output: 0.3
}

// ExampleWithDeleteCost demonstrates customizing the delete cost based on rune type.
func ExampleWithDeleteCost() {
	deleteCost := func(r rune) float64 {
		if r >= 'a' && r <= 'z' {
			return 1
		}
		return 3
	}
	l := levenshtein.New(levenshtein.WithDeleteCost(deleteCost))
	fmt.Println(l.Distance("abcXYZ", ""))
	// Output: 12
}

// ExampleWithInsertCost demonstrates customizing the insert cost based on rune type.
func ExampleWithInsertCost() {
	insertCost := func(r rune) float64 {
		if r >= '0' && r <= '9' {
			return 2 // 숫자는 삽입 비용이 2
		}
		return 1 // 그 외 문자는 기본 1
	}
	l := levenshtein.New(levenshtein.WithInsertCost(insertCost))
	fmt.Println(l.Distance("", "a1b2"))
	// Output: 6
}

// ExampleWithReplaceCost demonstrates customizing the replace cost based on rune similarity.
func ExampleWithReplaceCost() {
	replaceCost := func(a, b rune) float64 {
		if a == b {
			return 0
		}
		if a == '0' && b == 'O' || a == 'O' && b == '0' {
			return 0.5 // 유사한 문자는 낮은 비용
		}
		return 1
	}
	l := levenshtein.New(levenshtein.WithReplaceCost(replaceCost))
	fmt.Println(l.Distance("O0", "00"))
	// Output: 0.5
}

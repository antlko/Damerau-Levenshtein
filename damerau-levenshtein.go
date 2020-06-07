// Package tdl implements the true Damerau–Levenshtein distance.
//
// Reference:
// https://en.wikipedia.org/wiki/Damerau%E2%80%93Levenshtein_distance#Distance_with_adjacent_transpositions
package tdl

// Return the smalles int from a list
func minimum(is ...int) int {
	min := is[0]
	for _, i := range is {
		if min > i {
			min = i
		}
	}
	return min
}

var tdl = New(100)

// Distance is a shortcut func for doing a quick and dirty calculation,
// without having to set up your own struct and stuff.
// Not thread safe!
func Distance(a, b string) int {
	return tdl.Distance(a, b)
}

////////////////////////////////////////////////////////////////////////////////

// TrueDamerauLevenshtein is a struct that allocates memory only once, which is
// used when running Distance().
// This whole struct and associated functions are not thread safe in any way,
// that will be the callers responsibility! At least for now...
type TrueDamerauLevenshtein struct {
	maxSize int
	matrix  [][]int
	da      map[rune]int
}

// New initializes a new struct which allocates memory only once, to be used by
// Distance().
// maxSize sets an upper limit for both input strings used in Distance().
func New(maxSize int) *TrueDamerauLevenshtein {
	t := &TrueDamerauLevenshtein{
		maxSize: maxSize,
		matrix:  make([][]int, maxSize),
		da:      make(map[rune]int),
	}
	for i := range t.matrix {
		t.matrix[i] = make([]int, maxSize)
	}
	return t
}

// Distance calculates and returns the true Damerau–Levenshtein distance of string A and B.
// It's the caller's responsibility if he wants to trim whitespace or fix lower/upper cases.
// Distance is also free from memory allocs and is pretty quick.
func (t *TrueDamerauLevenshtein) Distance(a, b string) int {
	lenA, lenB := len(a), len(b)
	switch {
	case lenA < 1:
		return lenB
	case lenB < 1:
		return lenA
	case lenA > t.maxSize:
		return -1
	case lenB > t.maxSize:
		return -1
	}

	t.matrix[0][0] = lenA + lenB + 1
	for i := 0; i <= lenA; i++ {
		t.matrix[i+1][1] = i
		t.matrix[i+1][0] = t.matrix[0][0]
	}
	for j := 0; j <= lenB; j++ {
		t.matrix[1][j+1] = j
		t.matrix[0][j+1] = t.matrix[0][0]
	}

	for _, r := range a + b {
		t.da[r] = 0
	}

	for i := 1; i <= lenA; i++ {
		db := 0
		for j := 1; j <= lenB; j++ {
			i1 := t.da[rune(b[j-1])]
			j1 := db
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
				db = j
			}

			// By "conventional wisdom", the costs for the ins/del/trans operations are always +1
			t.matrix[i+1][j+1] = minimum(
				t.matrix[i][j]+cost,                  // substitution
				t.matrix[i+1][j]+1,                   // insertion
				t.matrix[i][j+1]+1,                   // deletion
				t.matrix[i1][j1]+(i-i1-1)+1+(j-j1-1), // transposition
			)
		}
		t.da[rune(a[i-1])] = i
	}
	return t.matrix[lenA+1][lenB+1]
}

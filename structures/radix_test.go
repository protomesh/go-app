package structures_test

import (
	"strconv"
	"testing"

	"github.com/protomesh/go-app/structures"

	"github.com/stretchr/testify/assert"
)

func TestRadixMatch(t *testing.T) {
	r := structures.NewRadixTree[int]()
	r.Insert("foo", 1)
	r.Insert("foobar", 2)
	r.Insert("bar", 3)
	r.Insert("barbaz", 4)
	r.Insert("barfoo", 5)
	r.Insert("barfoobaz", 6)
	r.Insert("barbar", 7)
	r.Insert("barbarbaz", 8)
	r.Insert("barbarfoo", 9)
	r.Insert("barbarfoobaz", 10)
	r.Insert("barbarbar", 11)
	r.Insert("barbarbarbaz", 12)
	r.Insert("barbarbarfoo", 13)
	r.Insert("barbarbarfoobaz", 14)
	r.Insert("barbarbarbar", 15)
	r.Insert("barbarbarbarbaz", 16)
	r.Insert("barbarbarbarfoo", 17)
	r.Insert("barbarbarbarfoobaz", 18)

	assert.Equal(t, 1, r.Match("foo"))
	assert.Equal(t, 2, r.Match("foobar"))
	assert.Equal(t, 3, r.Match("bar"))
	assert.Equal(t, 4, r.Match("barbaz"))
	assert.Equal(t, 5, r.Match("barfoo"))
	assert.Equal(t, 6, r.Match("barfoobaz"))
	assert.Equal(t, 7, r.Match("barbar"))
	assert.Equal(t, 8, r.Match("barbarbaz"))
	assert.Equal(t, 9, r.Match("barbarfoo"))
	assert.Equal(t, 10, r.Match("barbarfoobaz"))
	assert.Equal(t, 11, r.Match("barbarbar"))
	assert.Equal(t, 12, r.Match("barbarbarbaz"))
	assert.Equal(t, 13, r.Match("barbarbarfoo"))
	assert.Equal(t, 14, r.Match("barbarbarfoobaz"))
	assert.Equal(t, 15, r.Match("barbarbarbar"))
	assert.Equal(t, 16, r.Match("barbarbarbarbaz"))
	assert.Equal(t, 17, r.Match("barbarbarbarfoo"))
	assert.Equal(t, 18, r.Match("barbarbarbarfoobaz"))

	assert.Equal(t, 0, r.Match("xaxaxaxaxa"))

}

func TestRadixMatchLongest(t *testing.T) {

	r := structures.NewRadixTree[int]()
	r.Insert("foo", 1)
	r.Insert("foobar", 2)
	r.Insert("bar", 3)
	r.Insert("barbaz", 4)
	r.Insert("barfoo", 5)
	r.Insert("barfoobaz", 6)
	r.Insert("barbar", 7)
	r.Insert("barbarbaz", 8)
	r.Insert("barbarfoo", 9)
	r.Insert("barbarfoobaz", 10)
	r.Insert("barbarbar", 11)
	r.Insert("barbarbarbaz", 12)
	r.Insert("barbarbarfoo", 13)
	r.Insert("barbarbarfoobaz", 14)
	r.Insert("barbarbarbar", 15)
	r.Insert("barbarbarbarbaz", 16)
	r.Insert("barbarbarbarfoo", 17)
	r.Insert("barbarbarbarfoo", 20)
	r.Insert("barbarbarbarfoobaz", 18)

	assert.Equal(t, 1, r.MatchLongest("fooo"))
	assert.Equal(t, 1, r.MatchLongest("foo"))
	assert.Equal(t, 2, r.MatchLongest("foobar"))
	assert.Equal(t, 1, r.MatchLongest("foobaar"))
	assert.Equal(t, 3, r.MatchLongest("bar"))
	assert.Equal(t, 3, r.MatchLongest("bararar"))
	assert.Equal(t, 4, r.MatchLongest("barbaz"))
	assert.Equal(t, 3, r.MatchLongest("barbbbbbaz"))
	assert.Equal(t, 5, r.MatchLongest("barfoo"))
	assert.Equal(t, 5, r.MatchLongest("barfoooo"))
	assert.Equal(t, 6, r.MatchLongest("barfoobaz"))
	assert.Equal(t, 3, r.MatchLongest("barfobaobaz"))
	assert.Equal(t, 7, r.MatchLongest("barbar"))
	assert.Equal(t, 0, r.MatchLongest("brrrrarbar"))
	assert.Equal(t, 0, r.MatchLongest("baaarbarbaz"))
	assert.Equal(t, 7, r.MatchLongest("barbarbbbfoo"))
	assert.Equal(t, 7, r.MatchLongest("barbarfvvoobaz"))
	assert.Equal(t, 3, r.MatchLongest("barbbarbar"))
	assert.Equal(t, 3, r.MatchLongest("barrbarbarbaz"))
	assert.Equal(t, 3, r.MatchLongest("barbaarbarfoox"))
	assert.Equal(t, 3, r.MatchLongest("barbaaaaaarbarfoobaz"))
	assert.Equal(t, 0, r.MatchLongest("bnarbarbarbar"))
	assert.Equal(t, 15, r.MatchLongest("barbarbarbarbaaaaaz"))
	assert.Equal(t, 20, r.MatchLongest("barbarbarbarfoooooo"))
	assert.Equal(t, 15, r.MatchLongest("barbarbarbarfffffoobaz"))

	assert.Equal(t, 0, r.MatchLongest("xaxaxaxaxa"))
	assert.Equal(t, 0, r.MatchLongest("fo"))

}

func TestRadixDelete(t *testing.T) {
	r := structures.NewRadixTree[int]()
	r.Insert("foo", 1)
	r.Insert("foobar", 2)
	r.Insert("bar", 3)
	r.Insert("barbaz", 4)
	r.Insert("barfoo", 5)
	r.Insert("barfoobaz", 6)
	r.Insert("barbar", 7)
	r.Insert("barbarbaz", 8)
	r.Insert("barbarfoo", 9)
	r.Insert("barbarfoobaz", 10)
	r.Insert("barbarbar", 11)
	r.Insert("barbarbarbaz", 12)
	r.Insert("barbarbarfoo", 13)
	r.Insert("barbarbarfoobaz", 14)
	r.Insert("barbarbarbar", 15)
	r.Insert("barbarbarbarbaz", 16)
	r.Insert("barbarbarbarfoo", 17)
	r.Insert("barbarbarbarfoobaz", 18)

	r.Delete("foo")
	r.Delete("foobar")
	r.Delete("bar")
	r.Delete("barbaz")
	r.Delete("barfoo")
	r.Delete("barfoobaz")
	r.Delete("barbar")
	r.Delete("barbarbaz")
	r.Delete("barbarfoo")
	r.Delete("barbarfoobaz")
	r.Delete("barbarbar")
	r.Delete("barbarbarbaz")
	r.Delete("barbarbarfoo")
	r.Delete("barbarbarfoobaz")
	r.Delete("barbarbarbar")
	r.Delete("barbarbarbarbaz")
	r.Delete("barbarbarbarfoo")
	r.Delete("barbarbarbarfoobaz")

	assert.Equal(t, 0, r.Match("foo"))
	assert.Equal(t, 0, r.Match("foobar"))
	assert.Equal(t, 0, r.Match("bar"))
	assert.Equal(t, 0, r.Match("barbaz"))
	assert.Equal(t, 0, r.Match("barfoo"))
	assert.Equal(t, 0, r.Match("barfoobaz"))
	assert.Equal(t, 0, r.Match("barbar"))
	assert.Equal(t, 0, r.Match("barbarbaz"))
	assert.Equal(t, 0, r.Match("barbarfoo"))
	assert.Equal(t, 0, r.Match("barbarfoobaz"))
	assert.Equal(t, 0, r.Match("barbarbar"))
	assert.Equal(t, 0, r.Match("barbarbarbaz"))
	assert.Equal(t, 0, r.Match("barbarbarfoo"))
	assert.Equal(t, 0, r.Match("barbarbarfoobaz"))
	assert.Equal(t, 0, r.Match("barbarbarbar"))
	assert.Equal(t, 0, r.Match("barbarbarbarbaz"))
	assert.Equal(t, 0, r.Match("barbarbarbarfoo"))
	assert.Equal(t, 0, r.Match("barbarbarbarfoobaz"))

	assert.Equal(t, 0, r.MatchLongest("fooo"))
	assert.Equal(t, 0, r.MatchLongest("foo"))
	assert.Equal(t, 0, r.MatchLongest("foobar"))
	assert.Equal(t, 0, r.MatchLongest("foobaar"))
	assert.Equal(t, 0, r.MatchLongest("bar"))
	assert.Equal(t, 0, r.MatchLongest("bararar"))
	assert.Equal(t, 0, r.MatchLongest("barbaz"))
	assert.Equal(t, 0, r.MatchLongest("barbbbbbaz"))
	assert.Equal(t, 0, r.MatchLongest("barfoo"))
	assert.Equal(t, 0, r.MatchLongest("barfoooo"))
	assert.Equal(t, 0, r.MatchLongest("barfoobaz"))
	assert.Equal(t, 0, r.MatchLongest("barfobaobaz"))
	assert.Equal(t, 0, r.MatchLongest("barbar"))
	assert.Equal(t, 0, r.MatchLongest("brrrrarbar"))
	assert.Equal(t, 0, r.MatchLongest("baaarbarbaz"))
	assert.Equal(t, 0, r.MatchLongest("barbarbbbfoo"))
	assert.Equal(t, 0, r.MatchLongest("barbarfvvoobaz"))
	assert.Equal(t, 0, r.MatchLongest("barbbarbar"))
	assert.Equal(t, 0, r.MatchLongest("barrbarbarbaz"))
	assert.Equal(t, 0, r.MatchLongest("barbaarbarfoox"))
	assert.Equal(t, 0, r.MatchLongest("barbaaaaaarbarfoobaz"))
	assert.Equal(t, 0, r.MatchLongest("bnarbarbarbar"))
	assert.Equal(t, 0, r.MatchLongest("barbarbarbarbaaaaaz"))
	assert.Equal(t, 0, r.MatchLongest("barbarbarbarfoooooo"))
	assert.Equal(t, 0, r.MatchLongest("barbarbarbarfffffoobaz"))

	assert.Equal(t, 0, r.MatchLongest("xaxaxaxaxa"))
	assert.Equal(t, 0, r.MatchLongest("fo"))

}

func BenchmarkRadixGenerateString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.FormatInt(int64(i), 10)
	}
}
func BenchmarkMapIndex(b *testing.B) {
	r := make(map[string]int)
	b.Run("insert", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r[strconv.FormatInt(int64(i), 10)] = i
		}
	})
	b.Run("match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = r[strconv.FormatInt(int64(i), 10)]
		}
	})
}

func BenchmarkRadix(b *testing.B) {
	r := structures.NewRadixTree[int]()
	b.Run("insert", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r.Insert(strconv.FormatInt(int64(i), 10), i)
		}
	})
	b.Run("match", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r.Match(strconv.FormatInt(int64(i), 10))
		}
	})
	b.Run("matchlongest", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r.MatchLongest(strconv.FormatInt(int64(i), 10))
		}
	})
}

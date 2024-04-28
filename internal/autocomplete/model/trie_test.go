package model

import (
	"testing"
)

func TestNewTrie(t *testing.T) {

	t.Run("OK", func(t *testing.T) {
		trie := NewTrie()
		t.Logf("actual value: %v", trie)
	})
}

func TestInsert(t *testing.T) {
	type Input struct {
		word      string
		prequency int64
		trie      Trie
	}

	type TestCase struct {
		Name     string
		Input    Input
		Expected bool
	}

	testCases := []TestCase{
		{
			Name: "Success Case",
			Input: Input{
				word:      "test",
				prequency: 1,
				trie:      NewTrie(),
			},
			Expected: true,
		},
	}

	for _, v := range testCases {
		t.Run(v.Name, func(t *testing.T) {
			v.Input.trie.Insert(v.Input.word, v.Input.prequency)
			if exist := v.Input.trie.Contain(v.Input.word, v.Input.prequency); exist {
				t.Logf("actual value: %v; expected value: %v", exist, v.Expected)
			} else {
				t.Errorf("actual value: %v; expected value: %v", exist, v.Expected)
			}
		})
	}
}

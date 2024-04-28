package model

import (
	"encoding/json"
	"strings"
	"time"
)

type Trie struct {
	Root      *TrieNode  `json:"root" bson:"root"`
	CreatedAt *time.Time `json:"createdAt" bson:"createdAt"`
}

func NewTrie() Trie {
	createdAt := time.Now()
	return Trie{
		Root:      NewTrieNode(),
		CreatedAt: &createdAt,
	}
}

// should add limit lenght of word will insert
func (trie *Trie) Insert(word string, prequency int64) {

	word = strings.TrimSpace(word)

	current := trie.Root
	for i := range word {
		k := word[0:i]
		node, exist := current.Children[k]
		if !exist {
			node = NewTrieNode()
			current.Children[k] = node
		}
		current = node
	}

	current.Prequency = prequency
	current.EndOfWord = true
}

func (trie *Trie) Search(word string) (result []string) {
	current := trie.Root
	if current == nil {
		return result
	}
	for i := 0; i < len(word); i++ {
		node, exist := current.Children[word[0:i]]
		if !exist {
			return result
		}
		current = node
	}

	trie.searchChild(current, &result)
	return result
}

// traverse child map of identify node
func (trie *Trie) searchChild(node *TrieNode, result *[]string) {
	if len(node.Children) == 0 {
		return
	}
	for k := range node.Children {
		*result = append(*result, k)
		trie.searchChild(node.Children[k], result)
	}
}

func (trie Trie) MarshalBinary() ([]byte, error) {
	return json.Marshal(trie)
}

func (s *Trie) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (trie *Trie) Contain(word string, prequency int64) bool {
	current := trie.Root
	if current == nil {
		return false
	}
	for i := 0; i < len(word); i++ {
		node, exist := current.Children[word[0:i]]
		if !exist && node.Prequency == prequency && word[0:1] == word {
			return false
		}
		current = node
	}
	return true
}

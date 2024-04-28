package model

import "encoding/json"

type TrieNode struct {
	Children  map[string]*TrieNode `bson:"children" json:"children"`
	EndOfWord bool                 `bson:"endOfWord" json:"endOfWord"`
	Prequency int64                `bson:"prequency" json:"prequency"`
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		Children:  map[string]*TrieNode{},
		EndOfWord: false,
	}
}

func (node TrieNode) MarshalBinary() ([]byte, error) {
	return json.Marshal(node)
}

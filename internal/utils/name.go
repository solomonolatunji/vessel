package utils

import (
	"crypto/rand"
	"fmt"
)

var adjectives = []string{"bold", "brave", "calm", "clever", "cool", "eager", "fierce", "gentle", "happy", "jolly", "kind", "lively", "proud", "quiet", "rapid", "sharp", "smart", "solid", "swift", "wild", "shy", "epic", "neat", "vast"}
var nouns = []string{"bear", "bird", "cat", "deer", "dog", "fox", "frog", "hawk", "lion", "lynx", "moth", "owl", "puma", "seal", "swan", "toad", "wolf", "worm", "yak", "zebu", "crab", "dove", "fish", "goat"}

func GenerateRandomName() string {
	b := make([]byte, 2)
	rand.Read(b)
	adj := adjectives[int(b[0])%len(adjectives)]
	noun := nouns[int(b[1])%len(nouns)]
	b2 := make([]byte, 2)
	rand.Read(b2)
	return fmt.Sprintf("%s-%s-%x", adj, noun, b2)
}

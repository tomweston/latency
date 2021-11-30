package namesgenerator

import (
	"fmt"
	"math/rand"
)

// Generator must be implemented by any type that wants to be a generator.
type Generator interface {
	Generate() string
}

// NameGenerator stores a random name
type NameGenerator struct {
	random *rand.Rand
}

// NewNameGenerator returns a new name generator.
func NewNameGenerator(seed int64) Generator {
	nameGenerator := &NameGenerator{
		random: rand.New(rand.New(rand.NewSource(99))),
	}
	nameGenerator.random.Seed(seed)

	return nameGenerator
}

// Generate returns a random name.
func (rn *NameGenerator) Generate() string {
	randomAdjective := ADJECTIVES[rn.random.Intn(len(ADJECTIVES))]
	randomNoun := NOUNS[rn.random.Intn(len(NOUNS))]
	randomName := fmt.Sprintf("%v_%v", randomAdjective, randomNoun)

	return randomName
}

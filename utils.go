package omegleapi

import "strings"

// serializeTopics serializes topics into a string usable over the Omegle API.
func serializeTopics(topics []string) (finalTopics string) {
	finalTopics = "%5B%22"
	for _, topic := range topics {
		finalTopics += topic + "%22%2C%22"
	}
	return strings.TrimSuffix(finalTopics, "%22%2C%22") + "%22%5D"
}

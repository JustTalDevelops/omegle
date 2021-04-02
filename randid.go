package omegleapi

import "math/rand"

func generateRandID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randIdSelection[rand.Intn(len(randIdSelection))]
	}
	return string(b)
}
package omegleapi

import (
	"errors"
	"math/rand"
)

// FrontManager is a manager that manages all fronts, also known as nodes, used by Omegle.
// The default fronts were last updated on April 1st, 2021.
type FrontManager struct {
	fronts []string
}

// NewFrontManager creates a new front manager with the fronts selected.
func NewFrontManager(fronts []string) *FrontManager {
	return &FrontManager{fronts}
}

// NewFrontManagerFromDefaultFronts creates a new front manager with the default fronts.
func NewFrontManagerFromDefaultFronts() *FrontManager {
	return NewFrontManager([]string{"front1", "front2", "front3", "front4", "front5", "front6", "front7", "front8", "front9", "front10", "front11", "front12", "front13", "front14", "front15", "front16", "front17", "front18", "front19", "front20", "front21", "front22", "front23", "front24", "front25", "front26", "front27", "front28", "front29", "front30", "front31", "front32", "front33", "front34", "front35", "front36", "front37", "front38", "front39", "front40", "front41", "front42", "front43", "front44", "front45", "front46", "front47", "front48"})
}

// FindFront finds a valid front to communicate to.
func (f *FrontManager) FindFront() (string, error) {
	if len(f.fronts) == 0 {
		return "", errors.New("no fronts in front manager found")
	}
	return "https://" + f.fronts[rand.Int() % len(f.fronts)] + ".omegle.com", nil
}
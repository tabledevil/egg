package ui

type GameState int

const (
	StateIntro GameState = iota
	StateQuestion
	StateTransition
	StateSuccess
)

package theme

import (
	"ctf-tool/pkg/game"
)

type Theme interface {
	Name() string
	Description() string
	Render(width, height int, q *game.Question, inputView string, hint string) string
}

var Registry = []Theme{}

func Register(t Theme) {
	Registry = append(Registry, t)
}

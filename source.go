package koth

import "github.com/df-mc/dragonfly/server/player"

type Source interface{}

type SourcePlayer struct {
	p *player.Player
}
type SourceCapture struct {
	winner *player.Player
}

package koth

import (
	"math"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
)

type Handler interface {
	HandleStartCapturing(ctx *event.Context, p *player.Player)
	HandleStopCapturing(ctx *event.Context, p *player.Player)

	HandleCapture(ctx *event.Context, p *player.Player)

	HandleStart(ctx *event.Context, src Source)
	HandleStop(ctx *event.Context, src Source)
}

type NopHandler struct{}

func (NopHandler) HandleStartCapturing(ctx *event.Context, p *player.Player) {}
func (NopHandler) HandleStopCapturing(ctx *event.Context, p *player.Player)  {}
func (NopHandler) HandleCapture(ctx *event.Context, p *player.Player)        {}
func (NopHandler) HandleStart(ctx *event.Context, src Source)                {}
func (NopHandler) HandleStop(ctx *event.Context, src Source)                 {}

type PlayerHandler struct {
	player.NopHandler
	p *player.Player
}

func NewPlayerHandler(p *player.Player) *PlayerHandler {
	return &PlayerHandler{
		p: p,
	}
}

func (*PlayerHandler) Name() string { return "KothHandler" }

func actuallyMoved(old, new mgl64.Vec3) bool {
	return old.X() != new.X() || old.Z() != new.Z()
}
func (k *PlayerHandler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	if actuallyMoved(k.p.Position(), newPos) {
		k.p.SendTip(math.Round(newPos[0]), math.Round(newPos[2]))
		for _, koth := range koths {
			if koth.started {
				if koth.captureArea.Vec3WithinXZ(newPos) {
					koth.StartCapturing(k.p)
				} else {
					koth.StopCapturing(k.p)
				}
			}
		}
	}
}

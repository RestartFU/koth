package main

import (
	"fmt"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/dragonfly-on-steroids/area"
	"github.com/dragonfly-on-steroids/claim"
	"github.com/dragonfly-on-steroids/koth"
	"github.com/dragonfly-on-steroids/moreHandlers"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sirupsen/logrus"
)

func main() {
	c := server.DefaultConfig()
	c.Players.SaveData = false
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel
	s := server.New(&c, log)
	s.Start()
	k := koth.NewKOTH(s.World(), area.NewVec2(mgl64.Vec2{10, 10}, mgl64.Vec2{15, 15}), 5*time.Second)
	k.Handle(&Handler{})
	koth.Register(k)
	for {
		p, err := s.Accept()
		if err != nil {
			return
		}
		k.Start(koth.SourcePlayer{})
		claimH := claim.NewClaimHandler(p, claim.NewWilderness(s.World(), "entered wilderness", "left wilderness"))
		kothH := koth.NewPlayerHandler(p)
		h := moreHandlers.NewPlayerHandler(claimH, kothH)
		p.Handle(h)
	}
}

type Handler struct {
	koth.NopHandler
}

func (*Handler) HandleCapture(ctx *event.Context, p *player.Player) {
	p.Message("gg")
}
func (*Handler) HandleStopCapturing(ctx *event.Context, p *player.Player) {
	p.Message("no longer capturing")
}
func (*Handler) HandleStartCapturing(ctx *event.Context, p *player.Player) {
	p.Message("now capturing")
}
func (*Handler) HandleStart(ctx *event.Context, src koth.Source) {
	fmt.Println("started")
}
func (*Handler) HandleStop(ctx *event.Context, src koth.Source) {
	fmt.Println("stopped")
}

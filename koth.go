package koth

import (
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/dragonfly-on-steroids/area"
)

// I'll add comments when I have time

var koths []*KOTH

func KOTHS() []*KOTH {
	return koths
}

func Register(k *KOTH) {
	koths = append(koths, k)
}

type KOTH struct {
	name            string
	captureArea     area.Vec2
	duration        time.Duration
	hMutex          sync.RWMutex
	h               Handler
	capturing       *player.Player
	shouldCaptureAt time.Time
	started         bool
}

func NewKOTH(name string, captureArea area.Vec2, duration time.Duration) *KOTH {
	return &KOTH{
		name:        name,
		captureArea: captureArea,
		duration:    duration,
		h:           NopHandler{},
	}
}

func (k *KOTH) Handle(h Handler) {
	k.hMutex.Lock()
	defer k.hMutex.Unlock()
	if h == nil {
		h = NopHandler{}
	}
	k.h = h
}

func (k *KOTH) Name() string            { return k.name }
func (k *KOTH) Started() bool           { return k.started }
func (k *KOTH) CaptureArea() area.Vec2  { return k.captureArea }
func (k *KOTH) Duration() time.Duration { return k.duration }
func (k *KOTH) handler() Handler        { return k.h }
func (k *KOTH) Capturing() (*player.Player, bool) {
	return k.capturing, k.capturing != nil
}
func (k *KOTH) Start(src Source) {
	if !k.started {
		ctx := event.C()
		k.handler().HandleStart(ctx, src)
		ctx.Continue(func() {
			k.started = true

		})
	}
}
func (k *KOTH) Stop(src Source) {
	if k.started {
		ctx := event.C()
		k.handler().HandleStop(ctx, src)
		ctx.Continue(func() {
			k.started = false
		})
	}

}
func (k *KOTH) StartCapturing(p *player.Player) {
	if k.started {
		if k.capturing != p {
			ctx := event.C()
			k.handler().HandleStartCapturing(ctx, p)
			ctx.Continue(func() {
				k.capturing = p
				k.shouldCaptureAt = time.Now().Add(k.duration)
				time.AfterFunc(k.duration, k.captureFunc(p))
			})
		}
	}
}
func (k *KOTH) StopCapturing(p *player.Player) {
	if k.started {
		if k.capturing == p {
			ctx := event.C()
			k.handler().HandleStopCapturing(ctx, p)
			ctx.Continue(func() {
				k.capturing = nil
				k.shouldCaptureAt = time.Now().Add(43830 * time.Minute)
			})
		}
	}
}
func (k *KOTH) captureFunc(p *player.Player) func() {
	return func() {
		if k.capturing != nil && k.capturing == p {
			if k.shouldCaptureAt.Before(time.Now()) || k.shouldCaptureAt.Equal(time.Now()) {
				ctx := event.C()
				k.h.HandleCapture(ctx, p)
				ctx.Continue(func() {
					k.Stop(SourceCapture{winner: p})
				})
			}
		}
	}
}

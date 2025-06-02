package painter

import (
	"image"
	"sync"

	"github.com/roman-mazur/architecture-lab-3/ui"
	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver
	queue    []Operation
	mu       sync.Mutex
	done     chan struct{}

	next screen.Texture
	prev screen.Texture

	mq messageQueue
}

var size = image.Pt(400, 400)

func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.done = make(chan struct{})

	go func() {
		for {
			select {
			case <-l.done:
				return
			default:
				l.mu.Lock()
				if len(l.queue) > 0 {
					op := l.queue[0]
					l.queue = l.queue[1:]
					l.mu.Unlock()
					if visualizer, ok := l.Receiver.(*ui.Visualizer); ok {
						op.Execute(visualizer)
					}
				} else {
					l.mu.Unlock()
				}
			}
		}
	}()
}

func (l *Loop) Post(op Operation) {
	l.mu.Lock()
	l.queue = append(l.queue, op)
	l.mu.Unlock()
}

func (l *Loop) StopAndWait() {
	close(l.done)
}

type messageQueue struct{}

func (mq *messageQueue) push(op Operation) {}

func (mq *messageQueue) pull() Operation {
	return nil
}

func (mq *messageQueue) empty() bool {
	return false
}

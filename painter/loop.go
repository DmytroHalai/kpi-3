package painter

import (
	"image"
	"time"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циклі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправлення останнього разу у Receiver

	mq messageQueue

	stop    chan struct{}
	stopped chan struct{}
	stopReq bool
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	// TODO: стартувати цикл подій.
	l.stop = make(chan struct{})
	l.stopped = make(chan struct{})

	go func() {
		defer close(l.stopped)
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-l.stop:
				return
			case <-ticker.C:
				if l.mq.empty() {
					continue
				}
				op := l.mq.pull()
				if op == nil {
					continue
				}
				if op.Do(l.next) {
					l.Receiver.Update(l.next)
					l.next, l.prev = l.prev, l.next
				}
			}
		}
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	if update := op.Do(l.next); update {
		l.Receiver.Update(l.next)
		l.next, l.prev = l.prev, l.next
	}
}

// StopAndWait сигналізує про необхідність завершити цикл та блокується до моменту його повної зупинки.
func (l *Loop) StopAndWait() {
}

// TODO: Реалізувати чергу подій.
type messageQueue struct{}

func (mq *messageQueue) push(op Operation) {}

func (mq *messageQueue) pull() Operation {
	return nil
}

func (mq *messageQueue) empty() bool {
	return false
}

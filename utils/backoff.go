package utils

import (
	"errors"
	"math"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type BackoffMethod uint

const (
	NoAttempt   BackoffMethod = 0
	ConstantBO  BackoffMethod = 1
	GeometricBO BackoffMethod = 2
	ListBO      BackoffMethod = 3
)

type Backoff struct {
	method     BackoffMethod
	mult       int
	delays     []time.Duration
	maxAttempt int
}

func NewBackoff(method BackoffMethod, mult int, delays []time.Duration, maxAttempt int) (b *Backoff, err error) {
	if len(delays) < 1 {
		err = errors.New("One initial dealy is required for cvalculationg back off")
		return
	}
	b = &Backoff{
		method,
		mult,
		delays,
		maxAttempt,
	}
	return
}

func (b *Backoff) AttemptAllowed(attempt int) (ok bool) {
	ok = attempt < b.maxAttempt
	return
}

func (b *Backoff) GetDelay(attempt int) (delay time.Duration) {
	switch b.method {
	case ConstantBO:
		delay = b.getConstanDelay()
	case GeometricBO:
		delay = b.getGeometricDelay(attempt)
	case ListBO:
		delay = b.getListDelay(attempt)
	default:
		delay = 0 * time.Second
	}

	return
}

func (b *Backoff) getConstanDelay() (delay time.Duration) {
	delay = b.delays[0]
	return
}

func (b *Backoff) getGeometricDelay(attempt int) (delay time.Duration) {
	fMult, fAttempt := float64(b.mult), float64(attempt)
	delay = b.delays[0] * time.Duration(math.Pow(fMult, fAttempt))
	return
}

func (b *Backoff) getListDelay(attempt int) (delay time.Duration) {
	if attempt < len(b.delays) {
		delay = b.delays[attempt]
	} else {
		delay = 0 * time.Second
	}
	return
}

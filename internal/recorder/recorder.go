package recorder

import (
	"github.com/bricks-cloud/bricksllm/internal/event"
	"github.com/bricks-cloud/bricksllm/internal/key"
)

type Recorder struct {
	s  Store
	c  Cache
	ce CostEstimator
	es EventsStore
}

type EventsStore interface {
	InsertEvent(e *event.Event) error
}

type Store interface {
	IncrementCounter(keyId string, incr int64) error
}

type Cache interface {
	IncrementCounter(keyId string, rateLimitUnit key.TimeUnit, incr int64) error
}

type CostEstimator interface {
	EstimatePromptCost(model string, tks int) (float64, error)
	EstimateCompletionCost(model string, tks int) (float64, error)
}

func NewRecorder(s Store, c Cache, ce CostEstimator, es EventsStore) *Recorder {
	return &Recorder{
		s:  s,
		c:  c,
		ce: ce,
		es: es,
	}
}

func (r *Recorder) RecordKeySpend(keyId string, model string, micros int64, costLimitUnit key.TimeUnit) error {
	err := r.s.IncrementCounter(keyId, micros)
	if err != nil {
		return err
	}

	if len(costLimitUnit) != 0 {
		err = r.c.IncrementCounter(keyId, costLimitUnit, int64(micros))
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Recorder) RecordEvent(e *event.Event) error {
	return r.es.InsertEvent(e)
}

package core

import "github.com/FriedCoderZ/friedbot/cache"

const (
	TriggerModeNever = iota
	TriggerModeAll
	TriggerModeAny
	TriggerModeAlways
)

// Params 用于Stream中Handler函数间的参数传递
type Params any

type DefaultParams struct {
	Cache        *cache.Cache
	handlerIndex int
}

type handlerFunc func(stream *Stream, params Params) error

type Trigger interface {
	IsValid(cache *cache.Cache) bool
}

type Stream struct {
	triggers    []Trigger
	triggerMode int
	isTriggered bool
	handlers    []handlerFunc
}

func NewStream(triggers []Trigger, triggerMode int, handlers ...handlerFunc) *Stream {
	return &Stream{
		triggers:    triggers,
		triggerMode: triggerMode,
		handlers:    handlers,
		isTriggered: false,
	}
}

func (s *Stream) IsTriggered() bool {
	return s.isTriggered
}

func (s *Stream) Trigger(cache *cache.Cache) bool {
	var target int
	var passCount int
	var rejectCount int
	switch s.triggerMode {
	case TriggerModeNever:
		return false
	case TriggerModeAlways:
		s.isTriggered = false
		return false
	case TriggerModeAll:
		target = len(s.triggers)
	case TriggerModeAny:
		target = 1
	}
	for _, trigger := range s.triggers {
		if trigger.IsValid(cache) {
			passCount++
		} else {
			rejectCount++
		}
		if passCount >= target {
			s.isTriggered = true
			break
		} else if rejectCount >= len(s.triggers)-target {
			s.isTriggered = false
			break
		}
	}
	return s.isTriggered
}

func (s *Stream) Handle(cache *cache.Cache) (done bool, err error) {
	if !s.Trigger(cache) {
		return false, nil
	}
	params := &DefaultParams{
		Cache:        cache,
		handlerIndex: 0,
	}
	err = s.handlers[params.handlerIndex](s, params)
	if err != nil {
		return true, err
	}
	return true, nil
}

func (s *Stream) Next(params Params) (done bool, err error) {
	dp := params.(*DefaultParams)
	dp.handlerIndex++
	err = s.handlers[dp.handlerIndex](s, params)
	if err != nil {
		return false, err
	}
	return true, nil
}

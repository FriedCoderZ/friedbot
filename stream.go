package friedbot

const (
	TriggerModeNever = iota
	TriggerModeAll
	TriggerModeAny
	TriggerModeAlways
)

type triggerFunc func(*Context) bool

type handlerFunc func(*Context)

type Stream struct {
	triggers    []triggerFunc
	triggerMode int
	isTriggered bool
	handlers    []handlerFunc
}

func NewStream(triggerMode int) *Stream {
	return &Stream{
		triggers:    []triggerFunc{},
		handlers:    []handlerFunc{},
		triggerMode: triggerMode,
		isTriggered: false,
	}
}

func (s *Stream) AppendTriggers(triggers ...triggerFunc) {
	s.triggers = append(s.triggers, triggers...)
}

func (s *Stream) AppendHandlers(handlers ...handlerFunc) {
	s.handlers = append(s.handlers, handlers...)
}
func (s *Stream) IsTriggered() bool {
	return s.isTriggered
}

func (s *Stream) Trigger(ctx *Context) bool {
	var target int
	var passCount int
	switch s.triggerMode {
	case TriggerModeNever:
		return false
	case TriggerModeAlways:
		s.isTriggered = true
		return true
	case TriggerModeAll:
		target = len(s.triggers)
	case TriggerModeAny:
		target = 1
	}
	for index, trigger := range s.triggers {
		if trigger(ctx) {
			passCount++
		}
		if passCount >= target {
			s.isTriggered = true
			break
		}
		if passCount+len(s.triggers)-index < target {
			s.isTriggered = false
			break
		}
	}
	return s.isTriggered
}

func (s *Stream) Handle(ctx *Context) bool {
	if !s.Trigger(ctx) {
		return false
	}
	ctx.Next()
	return true
}

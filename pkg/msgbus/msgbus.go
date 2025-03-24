package msgbus

import "fmt"

type Msgbus struct {
	bus chan Message
}

func New() Msgbus {
	msg := make(chan Message)
	bus := Msgbus{bus: msg}
	return bus
}

func (m Msgbus) Debug(msg string) {
	m.bus <- Message{
		msg:   msg,
		level: MessageDebugLevel,
	}
}

func (m Msgbus) Debugf(format string, args ...any) {
	m.bus <- Message{
		msg:   fmt.Sprintf(format, args...),
		level: MessageDebugLevel,
	}
}

func (m Msgbus) Info(msg string) {
	m.bus <- Message{
		msg:   msg,
		level: MessageInfoLevel,
	}
}

func (m Msgbus) Infof(format string, args ...any) {
	m.bus <- Message{
		msg:   fmt.Sprintf(format, args...),
		level: MessageInfoLevel,
	}
}

func (m Msgbus) Warning(msg string) {
	m.bus <- Message{
		msg:   msg,
		level: MessageWarningLevel,
	}
}

func (m Msgbus) Warningf(format string, args ...any) {
	m.bus <- Message{
		msg:   fmt.Sprintf(format, args...),
		level: MessageWarningLevel,
	}
}

func (m Msgbus) Error(err error) {
	m.bus <- Message{
		err:   err,
		level: MessageErrorLevel,
	}
}

func (m Msgbus) Errorf(format string, args ...any) {
	m.bus <- Message{
		msg:   fmt.Sprintf(format, args...),
		level: MessageErrorLevel,
	}
}

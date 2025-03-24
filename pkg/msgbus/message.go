package msgbus

type MessageLevel int

const (
	MessageDebugLevel MessageLevel = iota
	MessageInfoLevel
	MessageWarningLevel
	MessageErrorLevel
)

func (m MessageLevel) String() string {
	switch m {
	case MessageDebugLevel:
		return "debug"
	case MessageInfoLevel:
		return "info"
	case MessageWarningLevel:
		return "warning"
	case MessageErrorLevel:
		return "error"
	default:
		return ""
	}
}

type Message struct {
	msg   string
	err   error
	level MessageLevel
}

func (m Message) String() string {
	if m.err != nil && m.err.Error() != "" {
		return m.err.Error()
	} else {
		return m.msg
	}
}

func (m Message) Level() MessageLevel {
	return m.level
}

func (m Message) Error() error {
	return m.err
}

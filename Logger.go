package connect

// Logger defines the interface understood by the Connect client for logging.
type Logger interface {
	Printf(format string, v ...interface{})
}

type discardLog struct{}

func (*discardLog) Printf(format string, v ...interface{}) {
}

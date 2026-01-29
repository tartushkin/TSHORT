package audit

import "log"

type Dispatcher struct {
	loggers []Logger
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) AddLogger(logger Logger) {
	d.loggers = append(d.loggers, logger)
}

func (d *Dispatcher) Dispatch(event Event) {
	for _, logger := range d.loggers {
		go func(l Logger) {
			if err := l.Log(event); err != nil {
				// Логируем ошибку аудита (например, в stderr)
				log.Printf("audit error: %v", err)
			}
		}(logger)
	}
}

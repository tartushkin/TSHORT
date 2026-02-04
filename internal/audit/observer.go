// Package audit реализует систему аудита событий с поддержкой множественных бэкендов.
//
// Dispatcher позволяет регистрировать несколько логгеров (например, файл, HTTP, Kafka)
// и асинхронно отправлять события во все из них.
//
// Пример:
//
//	dispatcher := audit.NewDispatcher()
//	dispatcher.AddLogger(&FileLogger{Path: "audit.log"})
//	dispatcher.AddLogger(&HTTPLoader{URL: "https://audit.example.com"})
//
//	dispatcher.Dispatch(logrus.StandardLogger(), audit.Event{
//	    Action: "url.create",
//	    UserID: "user123",
//	})
package audit

import (
	"github.com/sirupsen/logrus"
)

// Dispatcher управляет отправкой событий аудита в зарегистрированные логгеры.
//
// Отправка происходит асинхронно: каждое событие рассылается параллельно.
// Ошибки при отправке логируются через переданный *logrus.Logger.
type Dispatcher struct {
	loggers []Logger
}

// NewDispatcher создаёт новый диспетчер аудита.
//
// Возвращает указатель на инициализированный *Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

// AddLogger добавляет логгер в список получателей событий.
//
// Параметры:
//   - logger: реализация интерфейса audit.Logger
//
// После добавления, все события будут отправляться и в этот логгер.
func (d *Dispatcher) AddLogger(logger Logger) {
	d.loggers = append(d.loggers, logger)
}

// Dispatch отправляет событие аудита во все зарегистрированные логгеры.
//
// Отправка выполняется асинхронно в отдельных горутинах.
// Если один из логгеров завершается с ошибкой — ошибка логируется,
// но не прерывает отправку другим логгерам.
//
// Параметры:
//   - lg: *logrus.Logger для записи ошибок диспетчера
//   - event: событие аудита
func (d *Dispatcher) Dispatch(lg *logrus.Logger, event Event) {
	for _, logger := range d.loggers {
		go func(l Logger) {
			if err := l.Log(event); err != nil {
				// Логируем ошибку аудита (например, в stderr)
				lg.Errorf("dispatch.err -  %v", err)
			}
		}(logger)
	}
}

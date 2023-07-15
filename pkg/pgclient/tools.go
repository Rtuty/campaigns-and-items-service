package pgclient

import (
	"time"
)

// doWithTries пытается подключиться к клиенту указанное количество раз (attempts), прежде чем объявить об ошибке, если клиент БД не получен
func doWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err := fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}

		return nil
	}
	return
}

package config

import (
	"log"
	"time"
)

func Retry(attempts int, delay time.Duration, name string, fn func() error) error {
	var err error

	for i := 1; i <= attempts; i++ {
		err = fn()
		if err == nil {
			log.Printf("%s connected successfully", name)
			return nil
		}

		log.Printf("%s connection failed (attempt %d/%d): %v", name, i, attempts, err)
		time.Sleep(delay)
	}

	return err
}

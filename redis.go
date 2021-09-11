package main

import (
	"os"

	"github.com/gomodule/redigo/redis"
)

func setLastPublishDate(timestamp int64) error {
	c, err := dial()

	if err != nil {
		return err
	}

	defer c.Close()

	_, err = c.Do("SET", "noaa-nhc:last-updated", timestamp)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func getLastPublishDate() (int64, error) {
	c, err := dial()

	if err != nil {
		return 0, err
	}

	defer c.Close()

	str, err := redis.Int64(c.Do("GET", "noaa-nhc:last-updated"))

	if err != nil {
		return 0, nil
	}

	return str, nil
}

func dial() (redis.Conn, error) {
	return redis.DialURL(os.Getenv("REDIS_DSN"))
}
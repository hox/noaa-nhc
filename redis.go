package main

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
)

func setLastPublishDate(timestamp int64, walletId *string) error {
	c, err := dial()

	if err != nil {
		return err
	}

	defer c.Close()

	suffix := ""

	if walletId != nil {
		suffix += ":" + *walletId
	}

	_, err = c.Do("SET", fmt.Sprintf("noaa-nhc:last-updated%s", suffix), timestamp)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func getLastPublishDate(walletId *string) (int64, error) {
	c, err := dial()

	if err != nil {
		return 0, err
	}

	defer c.Close()

	suffix := ""

	if walletId != nil {
		suffix += ":" + *walletId
	}

	str, err := redis.Int64(c.Do("GET", fmt.Sprintf("noaa-nhc:last-updated%s", suffix)))

	if err != nil {
		return 0, nil
	}

	return str, nil
}

func dial() (redis.Conn, error) {
	return redis.DialURL(os.Getenv("REDIS_DSN"))
}

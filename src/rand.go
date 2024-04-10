package main

import (
	"math/rand"
	"time"
)

func createRand() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const pool = "qazwsxedcrfvtgbyhnujmikolpQAZWSXEDCRFVTGBYHNUJMIKOLP1234567890"
	bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = pool[r.Intn(len(pool))]
	}
	return string(bytes)
}

package main

import (
	"math/rand"
	"time"
)

func createRand(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const pool = "qazwsxedcrfvtgbyhnujmikolpQAZWSXEDCRFVTGBYHNUJMIKOLP1234567890"
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = pool[r.Intn(len(pool))]
	}
	return string(bytes)
}

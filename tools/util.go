package tools

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
)

// Get a random integer in a specified range
func RandomInRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

//Get unique identifier
func GetGuid(num int64) (int64, error) {
	node, err := snowflake.NewNode(num)
	if err != nil {
		return 0, err
	}

	id := node.Generate()
	return id.Int64(), nil
}

//  ----------------------- Random key -----------------------
const baseStr = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()[]{}+-*/_=."

// Generate random password
func GetRandomcode(length uint8) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63()))
	bytes := make([]byte, length)
	l := len(baseStr)
	for i := uint8(0); i < length; i++ {
		bytes[i] = baseStr[r.Intn(l)]
	}
	return string(bytes)
}

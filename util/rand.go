package util

import (
	"math/rand"
	"time"
)

// RandSeed 随机数种子
var RandSeed *rand.Rand

func init() {
	RandSeed = rand.New(rand.NewSource(time.Now().UnixNano()))
}

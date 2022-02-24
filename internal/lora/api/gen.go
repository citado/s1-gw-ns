package api

import (
	"encoding/hex"
	"math/rand"
	"sync"
)

// nolint: gochecknoglobals
var (
	rndEUI   *rand.Rand
	rndAddr  *rand.Rand
	onceAddr sync.Once
	onceEUI  sync.Once
)

func GenerateDevAddr() string {
	onceAddr.Do(func() {
		// nolint: gosec,gomnd
		rndAddr = rand.New(rand.NewSource(1378))
	})

	b := make([]byte, DevAddrLen)

	if _, err := rndAddr.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

func GenerateDevEUI() string {
	onceEUI.Do(func() {
		// nolint: gosec,gomnd
		rndEUI = rand.New(rand.NewSource(1378))
	})

	b := make([]byte, DevEUILen)

	if _, err := rndEUI.Read(b); err != nil {
		panic(err)
	}

	return hex.EncodeToString(b)
}

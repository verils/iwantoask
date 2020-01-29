package app

import (
	haikunator "github.com/atrox/haikunatorgo/v2"
	"log"
	"testing"
)

func TestHaikunate(t *testing.T) {
	haikunate := haikunator.New()
	haikunate.Delimiter = " "
	haikunate.TokenLength = 0
	haikunated := haikunate.Haikunate()

	log.Printf(haikunated)
}

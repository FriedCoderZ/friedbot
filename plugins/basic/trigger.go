package basic

import (
	"fmt"

	"github.com/FriedCoderZ/friedbot/cache"
)

type TrueTrigger struct {
}

func (t TrueTrigger) IsValid(cache *cache.Cache) bool {
	fmt.Println("true trigger")
	fmt.Println(cache)
	return true
}

type FalseTrigger struct {
}

func (f FalseTrigger) IsValid(cache *cache.Cache) bool {
	fmt.Println("false trigger")
	fmt.Println(cache)
	return false
}

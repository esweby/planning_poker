package poker

import (
	"fmt"
	"math/rand"
)

func Generate8DigitString() RoomID {
	return fmt.Sprintf("%08d", rand.Intn(100000000))
}

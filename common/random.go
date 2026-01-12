package common

import "crypto/rand"

func GenSalt() string {

	return rand.Text()

}

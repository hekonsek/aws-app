package awsom

import (
	"github.com/Pallinder/sillyname-go"
	"strings"
)

const UnixExitCodeGeneralError = 1

func RandomName() string {
	lowerCased := strings.ToLower(sillyname.GenerateStupidName())
	return strings.Replace(lowerCased, " ", "", -1)
}

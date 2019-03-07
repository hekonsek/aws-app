package awsom

import (
	"github.com/Pallinder/sillyname-go"
	"strings"
)

func RandomName() string {
	return strings.Replace(sillyname.GenerateStupidName(), " ", "", -1)
}

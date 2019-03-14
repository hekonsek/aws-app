package randomstrings

import (
	"crypto/sha1"
	"fmt"
	"github.com/Pallinder/sillyname-go"
	"strings"
	"time"
)

// GenerateLowercaseName returns "sillyname" (like "Pumacloud Lightning"), but without space and in lowercase
// i.e. "Pumacloud Lightning" becomes ("pumacloudlightning"). This function is useful for generating human-friendly
// random strings that should be used as semi-unique identifiers (as lowercase alphabetical-only string with
// reasonable size almost always can be used as valid identifier).
func GenerateLowercaseName() string {
	lowerCased := strings.ToLower(sillyname.GenerateStupidName())
	return strings.Replace(lowerCased, " ", "", -1)
}

func GenerateLowercaseNameWithHash() string {
	lowerCased := GenerateLowercaseName()

	hash := sha1.New()
	hash.Write([]byte(time.Now().Format(time.RFC850)))
	hashString := fmt.Sprintf("%x\n", hash.Sum(nil))
	nameWithHash := lowerCased + hashString[:6]

	return nameWithHash
}

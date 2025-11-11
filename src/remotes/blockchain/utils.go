package blockchain

import (
	"fmt"
	"logger/utils"
	"math"
	"strings"
)

// Format PGP bytes to a string signature
func formatSignatureToPGP(data []byte) string {
	size := 64
	str := utils.ToBase64(data)
	strLength := len(str)
	splitedLength := int(math.Ceil(float64(strLength) / float64(size)))
	splited := make([]string, splitedLength)
	var start, stop int
	for i := 0; i < splitedLength; i += 1 {
		start = i * size
		stop = start + size
		if stop > strLength {
			stop = strLength
		}
		splited[i] = str[start:stop]
	}
	return fmt.Sprintf("-----BEGIN PGP SIGNATURE-----\n\n%s\n-----END PGP SIGNATURE-----", strings.Join(splited, "\n"))
}

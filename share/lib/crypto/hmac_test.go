package crypto

import (
    "fmt"
	"strings"
	"testing"
)

func TestCrypto(t *testing.T) {
	channelName := "private-mychannel-11"
	socketId := "h5zRV67-praOqGvAS_6k"
	secret := "d902b719646ea7489928d24214776b30"

	tosign := strings.Join([]string{socketId, channelName}, ":")
	signature := HmacSignature(tosign, secret)

	fmt.Printf("%v\n", tosign)
	fmt.Printf("%v", signature)
}

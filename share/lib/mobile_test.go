package lib

import (
    "fmt"
	"testing"
)

// MobileEncrypt
func Test_MobileEncrypt(t *testing.T) {
	fmt.Println(MobileEncrypt("18602990888"))
}

// MobileDecrypt
func Test_MobileDecrypt(t *testing.T) {
	fmt.Println(MobileDecrypt("186*0888", "gpCLEoL7Te4IENz8fdOiYw=="))
}

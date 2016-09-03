package oregonpi

// #cgo CFLAGS: -Ofast
// #cgo CXXFLAGS: -Ofast
// #cgo LDFLAGS: -lwiringPi
// #include "oregonpi.h"
import "C"
import "fmt"

func Test() {
	fmt.Println("fuck yes, new sensor and shit")
}

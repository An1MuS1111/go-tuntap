package tuntap

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	IFNAMSIZ = 16
	// flags
	IFF_TUN   = 0x0001
	IFF_TAP   = 0x0002
	IFF_NO_PI = 0x1000

	TUNSETIFF = 0x400454ca
)

type ifreq struct {
	Name  [IFNAMSIZ]byte
	Flags uint16
	_     [22]byte // padding to match kernel size
}

/**
 * fd ‒ the fd to turn into TUN or TAP.
 * name ‒ the name to use. If empty, kernel will assign something by itself.
 *   Must be buffer with capacity at least 33.
 * mode ‒ 1 = TUN, 2 = TAP.
 * packet_info ‒ if packet info should be provided, if the given value is 0 it
 * will not prepend packet info.
 */

func tuntapSetup(fd uintptr, name string, mode Mode, packetInfo bool) (string, error) {
	req := &ifreq{}

	copy(req.Name[:], []byte(name))

	switch mode {
	case TUN:
		req.Flags = IFF_TUN
	case TAP:
		req.Flags = IFF_TAP

	default:
		return "", fmt.Errorf("invalid mode, expected TUN|TAP")
	}

	if !packetInfo {
		req.Flags |= IFF_NO_PI
	}

	// Perform ioctl
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		uintptr(TUNSETIFF),
		uintptr(unsafe.Pointer(req)),
	)
	if errno != 0 {
		return "", errno
	}

	// Extract actual interface name assigned by kernel
	actualName := string(req.Name[:])
	for i, b := range actualName {
		if b == 0 {
			actualName = actualName[:i]
			break
		}
	}

	return actualName, nil
}

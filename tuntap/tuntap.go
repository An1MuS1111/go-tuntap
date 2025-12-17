package tuntap

/*
#include <stdlib.h>
#include <string.h>
int tuntap_setup(int fd, unsigned char *name, int mode, int packet_info);
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

type Mode uint8

const (
	TUN Mode = iota + 1
	TAP
)

type Iface struct {
	fd   *os.File
	mode Mode
	name string
}

func NewIface(ifname string, mode Mode) (*Iface, error) {
	return withOption(ifname, mode, true)
}

func WithoutPacketInfo(ifname string, mode Mode) (*Iface, error) {
	return withOption(ifname, mode, false)
}

func withOption(ifname string, mode Mode, packetInfo bool) (*Iface, error) {
	fileName := "/dev/net/tun"
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to OPEN|CREATE %s file", fileName)
	}
	nameBuffer := make([]byte, len(ifname))
	nameBuffer = append(append(nameBuffer, []byte(ifname)...), make([]byte, 33)...)

	pi := 0
	if packetInfo {
		pi = 1
	}

	res, err := C.tuntap_setup(
		C.int(file.Fd()),
		(*C.uchar)(unsafe.Pointer(&nameBuffer[0])),
		C.int(mode),
		C.int(pi),
	)

	if res < 0 {
		file.Close()
		return nil, fmt.Errorf("tuntap_setup failed: %w", err)
	}

	actualName := C.GoString((*C.char)(unsafe.Pointer(&nameBuffer[0])))

	return &Iface{
		fd:   file,
		mode: mode,
		name: actualName,
	}, nil

}

func (i *Iface) Mode() Mode {
	return i.mode
}
func (i *Iface) Name() string {
	return i.name
}

func (i *Iface) Recv(buf []byte) (int, error) {
	return i.fd.Read(buf)
}
func (i *Iface) Send(buf []byte) (int, error) {
	return i.fd.Write(buf)
}

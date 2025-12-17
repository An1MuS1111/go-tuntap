package tuntap

import (
	"fmt"
	"os"
	"syscall"
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

	actualName, err := tuntapSetup(file.Fd(), ifname, mode, packetInfo)
	if err != nil {
		file.Close()
		return nil, err
	}

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

func (i *Iface) SetNonBlocking(nonBlock bool) error {
	return syscall.SetNonblock(int(i.fd.Fd()), nonBlock)
}

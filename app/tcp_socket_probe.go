package app

import (
	"net"
	"strconv"
)

type TcpSocketProbe struct {
	Port int `yaml:"port"`
}

func (p *TcpSocketProbe) Probe(_ *Process, _ *Probe) (bool, error) {
	serverAddress := "127.0.0.1:" + strconv.Itoa(p.Port)
	tcpAddress, err := net.ResolveTCPAddr("tcp", serverAddress)
	if err != nil {
		return false, err
	}

	connection, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		return false, err
	}
	defer connection.Close()
	return true, nil
}

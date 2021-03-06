package simulator

import (
	"context"
	"fmt"
	"math/rand"
	"net"
)

var (
	scanPorts    = []int{21, 22, 23, 25, 80, 88, 111, 135, 139, 143, 389, 443, 445, 1433, 1521, 3306, 3389, 5432, 5900, 6000, 8443}
	scanIPRanges = []*net.IPNet{
		{
			IP:   net.IPv4(10, 0, 0, 0),
			Mask: net.CIDRMask(8, 32),
		},
		{
			IP:   net.IPv4(172, 16, 0, 0),
			Mask: net.CIDRMask(12, 32),
		},
		{
			IP:   net.IPv4(192, 168, 0, 0),
			Mask: net.CIDRMask(16, 32),
		},
	}
)

// PortScan simulator.
type PortScan struct {
	portNo int
}

// NewPortScan creates port scan simulator.
func NewPortScan() *PortScan {
	return &PortScan{
		portNo: 10,
	}
}

// Simulate port scanning for given host.
func (*PortScan) Simulate(ctx context.Context, extIP net.IP, host string) error {
	d := &net.Dialer{
		LocalAddr: &net.TCPAddr{IP: extIP},
	}

	conn, err := d.DialContext(ctx, "tcp", host)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

// Hosts returns host:port generated from RFC 1918 addresses.
func (s *PortScan) Hosts(size int) ([]string, error) {
	var (
		hosts []string
		idx   = rand.Perm(len(scanPorts))
	)

	for i := 0; i < size; i++ {
		ip := scanIPRanges[rand.Intn(len(scanIPRanges))]
		ip.IP[len(ip.IP)-2] = byte(rand.Intn(256))
		ip.IP[len(ip.IP)-1] = byte(rand.Intn(256))

		for j := 0; j < s.portNo; j++ {
			port := scanPorts[idx[i]]
			hosts = append(hosts, fmt.Sprintf("%s:%d", ip.IP, port))
		}
	}

	return hosts, nil
}

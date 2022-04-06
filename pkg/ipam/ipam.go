package ipam

import (
	"errors"
	"fmt"
	"net"

	"mycni/pkg/config"
	"mycni/pkg/store"

	cip "github.com/containernetworking/plugins/pkg/ip"
)

var (
	IPOverflowError = errors.New("ip overflow")
)

type IPAM struct {
	subnet  *net.IPNet
	gateway net.IP
	store   *store.Store
}

func NewIPAM(conf *config.CNIConf, s *store.Store) (*IPAM, error) {
	_, ipnet, err := net.ParseCIDR(conf.Subnet)
	if err != nil {
		return nil, err
	}

	im := &IPAM{
		subnet: ipnet,
		store:  s,
	}

	im.gateway, err = im.NextIP(im.subnet.IP)
	if err != nil {
		return nil, err
	}

	return im, nil
}

func (im *IPAM) Mask() net.IPMask {
	return im.subnet.Mask
}

func (im *IPAM) Gateway() net.IP {
	return im.gateway
}

func (im *IPAM) IPNet(ip net.IP) *net.IPNet {
	return &net.IPNet{IP: ip, Mask: im.Mask()}
}

func (im *IPAM) NextIP(ip net.IP) (net.IP, error) {
	next := cip.NextIP(ip)
	if !im.subnet.Contains(next) {
		return nil, IPOverflowError
	}

	return next, nil
}

func (im *IPAM) AllocateIP(id, ifName string) (net.IP, error) {
	im.store.Lock()
	defer im.store.Unlock()

	if err := im.store.LoadData(); err != nil {
		return nil, err
	}

	ip, _ := im.store.GetIPByID(id)
	if len(ip) > 0 {
		return ip, nil
	}

	last := im.store.Last()
	if len(last) == 0 {
		last = im.gateway
	}

	start := make(net.IP, len(last))
	copy(start, last)
	for {
		next, err := im.NextIP(start)
		if err == IPOverflowError && !last.Equal(im.gateway) {
			start = im.gateway
			continue
		} else if err != nil {
			return nil, err
		}

		if !im.store.Contain(next) {
			err := im.store.Add(next, id, ifName)
			return next, err
		}

		start = next
		if start.Equal(last) {
			break
		}

		fmt.Printf("ip: %s", next)
	}

	return nil, fmt.Errorf("no available ip")
}

func (im *IPAM) ReleaseIP(id string) error {
	im.store.Lock()
	defer im.store.Unlock()

	if err := im.store.LoadData(); err != nil {
		return err
	}

	return im.store.Del(id)
}

func (im *IPAM) CheckIP(id string) (net.IP, error) {
	im.store.RLock()
	defer im.store.RUnlock()

	if err := im.store.LoadData(); err != nil {
		return nil, err
	}

	ip, ok := im.store.GetIPByID(id)
	if !ok {
		return nil, fmt.Errorf("failed to find container %s ip", id)
	}

	return ip, nil
}

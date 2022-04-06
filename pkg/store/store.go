package store

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"

	"github.com/alexflint/go-filemutex"
)

const (
	defaultDataDir = "/var/lib/cni"
)

type containerNetINfo struct {
	ID     string `json:"id"`
	IFName string `json:"if"`
}

type data struct {
	IPs  map[string]containerNetINfo `json:"ips"`
	Last string                      `json:"last"`
}

type Store struct {
	*filemutex.FileMutex
	dir      string
	data     *data
	dataFile string
}

func NewStore(dataDir, network string) (*Store, error) {
	if dataDir == "" {
		dataDir = defaultDataDir
	}
	dir := filepath.Join(dataDir, network)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	fl, err := newFileLock(dir)
	if err != nil {
		return nil, err
	}

	dataFile := filepath.Join(dir, network+".json")
	data := &data{IPs: make(map[string]containerNetINfo)}

	return &Store{fl, dir, data, dataFile}, nil
}

func (s *Store) LoadData() error {
	data := &data{}
	raw, err := os.ReadFile(s.dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(s.dataFile)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = f.Write([]byte("{}"))
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if err := json.Unmarshal(raw, &data); err != nil {
			return err
		}
	}
	if data.IPs == nil {
		data.IPs = make(map[string]containerNetINfo)
	}

	s.data = data
	return nil
}

func (s *Store) Last() net.IP {
	return net.ParseIP(s.data.Last)
}

func (s *Store) GetIPByID(id string) (net.IP, bool) {
	for ip, info := range s.data.IPs {
		if info.ID == id {
			return net.ParseIP(ip), true
		}
	}
	return nil, false
}

func (s *Store) Add(ip net.IP, id, ifname string) error {
	if len(ip) > 0 {
		s.data.IPs[ip.String()] = containerNetINfo{id, ifname}
		s.data.Last = ip.String()
		return s.Store()
	}
	return nil
}

func (s *Store) Del(id string) error {
	for ip, info := range s.data.IPs {
		if info.ID == id {
			delete(s.data.IPs, ip)
			return s.Store()
		}
	}
	return nil
}

func (s *Store) Contain(ip net.IP) bool {
	_, ok := s.data.IPs[ip.String()]
	return ok
}

func (s *Store) Store() error {
	raw, err := json.Marshal(s.data)
	if err != nil {
		return err
	}

	return os.WriteFile(s.dataFile, raw, 0644)
}

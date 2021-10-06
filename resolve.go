package doxyproxy

import (
	"net"
	"time"
)

//Resolve will return the true IP address of the client
func (p *Proxy) Resolve(addr net.Addr) (string, error) {
	_, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return "", err
	}

	p.mutex.RLock()
	entry, ok := p.cache[port]
	p.mutex.RUnlock()
	if ok == true {
		if entry.Expire > time.Now().Unix() {
			return entry.IP, nil
		}
	}

	entry.proxy = p
	entry.ID = port
	entry.proxy.API = p.API
	entry.proxy.App = p.App

	ip, err := entry.Fetch()
	if err != nil {
		return "", err
	}

	p.mutex.Lock()
	p.cache[port] = IPEntry{
		ID:     port,
		proxy:  p,
		Expire: time.Now().Add(p.CacheTTL).Unix(),
		IP:     ip,
	}
	p.mutex.Unlock()

	return ip, nil
}

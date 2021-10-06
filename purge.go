package doxyproxy

import "net"

//Purge will remove the IP from the cache
func (p *Proxy) Purge(addr net.Addr) error {
	_, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return err
	}

	p.mutex.RLock()
	_, ok := p.cache[port]
	p.mutex.RUnlock()
	if ok == false {
		return nil
	}

	p.mutex.Lock()
	delete(p.cache, port)
	p.mutex.Unlock()

	return nil
}

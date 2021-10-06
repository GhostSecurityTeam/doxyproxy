package doxyproxy

import (
	"encoding/json"
	"net"
	"net/http"
)

//Kill will disconnect the user from the proxy
func (p *Proxy) Kill(addr net.Addr) error {
	_, port, err := net.SplitHostPort(addr.String())
	if err != nil {
		return err
	}

	p.mutex.RLock()
	entry, ok := p.cache[port]
	p.mutex.RUnlock()
	if ok == true {
		p.mutex.Lock()
		delete(p.cache, port)
		p.mutex.Unlock()

		return entry.Kill()
	}

	entry = IPEntry{
		ID:    port,
		proxy: p,
	}

	return entry.Kill()
}

//Kill will disconnect the user from the proxy
func (entry IPEntry) Kill() error {

	req, err := http.NewRequest("GET", entry.proxy.API+entry.ID+"/kill", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", entry.proxy.Key)

	resp, err := entry.proxy.Client.Do(req)
	if err != nil {
		return err
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if response.Success == false {
		if resp.StatusCode == 404 {
			return ErrClientNotFound
		}

		return ErrAPIResponseSuccessFalse
	}

	return nil
}

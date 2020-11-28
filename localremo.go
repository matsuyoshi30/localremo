package localremo

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/grandcat/zeroconf"
)

func GetLocalRemoAddr() (*zeroconf.ServiceEntry, net.IP, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, nil, err
	}

	entries := make(chan *zeroconf.ServiceEntry)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = resolver.Browse(ctx, "_remo._tcp", "local.", entries)
	if err != nil {
		return nil, nil, err
	}

	select {
	case <-ctx.Done():
		return nil, nil, nil
	case entry := <-entries:
		return entry, entry.AddrIPv4[0], nil
	}
}

type IRSignal struct {
	Format string `json:"format"`
	Freq   uint   `json:"freq"`
	Data   []uint `json:"data"`
}

type LocalClient struct {
	Client *http.Client
}

func NewClient() *LocalClient {
	return &LocalClient{Client: &http.Client{}}
}

func (lc *LocalClient) Get(ctx context.Context, targetURL net.IP) (*IRSignal, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+targetURL.String()+"/messages", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Requested-With", "local")
	req.Header.Add("accept", "application/json")

	resp, err := lc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Response status is not OK")
	}
	defer resp.Body.Close()

	var ir IRSignal
	if err := json.NewDecoder(resp.Body).Decode(&ir); err != nil {
		return nil, err
	}

	return &ir, nil
}

func (lc *LocalClient) Post(ctx context.Context, targetURL net.IP, body io.Reader) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+targetURL.String()+"/messages", body)
	if err != nil {
		return err
	}
	req.Header.Add("X-Requested-With", "local")
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := lc.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Response status is not OK")
	}

	return nil
}

func ReadJSON(path string) ([]byte, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ir IRSignal
	if err := json.Unmarshal(content, &ir); err != nil {
		return nil, err
	}

	return content, nil
}

package app

import (
	"net/http"
	"strconv"
)

type HttpHeader struct {
	name  string `yaml:"name"`
	value string `yaml:"value"`
}

type HttpGetProbe struct {
	Path        string        `yaml:"path"`
	Port        int           `yaml:"port"`
	HttpHeaders []*HttpHeader `yaml:"httpHeaders"`
}

func (p *HttpGetProbe) Probe(_ *Process, c *Probe) (bool, error) {
	timeout, err := c.GetTimeout()
	if err != nil {
		return false, err
	}
	url := "http://127.0.0.1:" + strconv.Itoa(p.Port) + p.Path
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	for _, header := range p.HttpHeaders {
		request.Header.Set(header.name, header.value)
	}

	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 300 || response.StatusCode < 200 {
		return false, nil
	}
	return true, nil
}

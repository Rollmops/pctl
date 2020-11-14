package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"
)

type RemoteProcessController struct{}

type BaseArgs struct {
	Names   []string `json:"names"`
	Filters Filters  `json:"filters"`
}

func (p *RemoteProcessController) Start(names []string, filters Filters, comment string) error {
	data, err := json.Marshal(struct {
		BaseArgs `json:",inline"`
		Comment  string `json:"comment"`
	}{
		BaseArgs: BaseArgs{
			Names:   names,
			Filters: filters,
		},
		Comment: comment,
	})
	if err != nil {
		return err
	}

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", "/tmp/pctl-agent.sock")
			},
		},
	}
	response, err := httpc.Post("http://unix/pctl", "application/octet-stream", bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(os.Stdout, response.Body)
	return err
}

func (p *RemoteProcessController) Stop(names []string, filters Filters, noWait bool, kill bool) error {
	return nil
}
func (p *RemoteProcessController) Restart(names []string, filters Filters, comment string, kill bool) error {
	return nil
}
func (p *RemoteProcessController) Kill(names []string, filters Filters) error {
	return nil
}
func (p *RemoteProcessController) Info(names []string, format string, filters Filters, columns []string) error {
	return nil
}

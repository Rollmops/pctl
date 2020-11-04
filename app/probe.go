package app

import (
	"time"
)

type Probe struct {
	InitialDelay     string          `yaml:"initialDelay"`
	Period           string          `yaml:"period"`
	Timeout          string          `yaml:"timeout"`
	FailureThreshold int             `yaml:"failureThreshold"`
	SuccessThreshold int             `yaml:"successThreshold"`
	Exec             *ExecProbe      `yaml:"exec"`
	HttpGet          *HttpGetProbe   `yaml:"httpGet"`
	TcpSocket        *TcpSocketProbe `yaml:"tcpSocket"`
}

func (c *Probe) GetTimeout() (time.Duration, error) {
	if c.Timeout == "" {
		return 1 * time.Second, nil
	}
	return time.ParseDuration(c.Timeout)
}

func (c *Probe) GetPeriod(defaultPeriod time.Duration) (time.Duration, error) {
	if c.Period == "" {
		return defaultPeriod, nil
	}
	return time.ParseDuration(c.Period)
}

func (c *Probe) GetInitialDelay() (time.Duration, error) {
	if c.InitialDelay == "" {
		return 0, nil
	}
	return time.ParseDuration(c.InitialDelay)
}

func (c *Probe) Probe(process *Process, defaultPeriod time.Duration) (bool, error) {
	period, err := c.GetPeriod(defaultPeriod)
	if err != nil {
		return false, err
	}
	initialDelay, err := c.GetInitialDelay()
	if err != nil {
		return false, err
	}
	if c.FailureThreshold == 0 {
		c.FailureThreshold = 1
	}
	if c.SuccessThreshold == 0 {
		c.SuccessThreshold = 1
	}
	var probeImplFunc func(*Process, *Probe) (bool, error)
	if c.Exec != nil {
		probeImplFunc = c.Exec.Probe
	} else if c.HttpGet != nil {
		probeImplFunc = c.HttpGet.Probe
	} else if c.TcpSocket != nil {
		probeImplFunc = c.TcpSocket.Probe
	} else {
		nullProbe := &NullProbe{}
		probeImplFunc = nullProbe.Probe
	}
	successCount := 0
	failureCount := 0
	time.Sleep(initialDelay)
	for {
		ok, err := probeImplFunc(process, c)
		if err != nil || !ok {
			failureCount++
			if failureCount >= c.FailureThreshold {
				return false, err
			}
		} else {
			successCount++
			if successCount >= c.SuccessThreshold {
				return true, nil
			}
		}
		time.Sleep(period)
	}
}

type NullProbe struct{}

func (p *NullProbe) Probe(_ *Process, _ *Probe) (bool, error) {
	return true, nil
}

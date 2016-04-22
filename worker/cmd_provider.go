package worker

import (
	"errors"
	"os"

	"github.com/anmitsu/go-shlex"
)

type cmdConfig struct {
	name                        string
	upstreamURL, command        string
	workingDir, logDir, logFile string
	interval                    int
	env                         map[string]string
}

type cmdProvider struct {
	baseProvider
	cmdConfig
	command []string
	cmd     *cmdJob
}

func newCmdProvider(c cmdConfig) (*cmdProvider, error) {
	// TODO: check config options
	provider := &cmdProvider{
		baseProvider: baseProvider{
			name:     c.name,
			ctx:      NewContext(),
			interval: c.interval,
		},
		cmdConfig: c,
	}

	provider.ctx.Set(_WorkingDirKey, c.workingDir)
	provider.ctx.Set(_LogDirKey, c.logDir)
	provider.ctx.Set(_LogFileKey, c.logFile)

	cmd, err := shlex.Split(c.command, true)
	if err != nil {
		return nil, err
	}
	provider.command = cmd

	return provider, nil
}

func (p *cmdProvider) InitRunner() {
}

func (p *cmdProvider) Run() error {
	env := map[string]string{
		"TUNASYNC_MIRROR_NAME":  p.Name(),
		"TUNASYNC_WORKING_DIR":  p.WorkingDir(),
		"TUNASYNC_UPSTREAM_URL": p.upstreamURL,
		"TUNASYNC_LOG_FILE":     p.LogFile(),
	}
	for k, v := range p.env {
		env[k] = v
	}
	p.cmd = newCmdJob(p.command, p.WorkingDir(), env)

	logFile, err := os.OpenFile(p.LogFile(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	// defer logFile.Close()
	p.cmd.SetLogFile(logFile)

	return p.cmd.Start()
}

func (p *cmdProvider) Wait() error {
	return p.cmd.Wait()
}

func (p *cmdProvider) Terminate() error {
	if p.cmd == nil {
		return errors.New("provider command job not initialized")
	}
	err := p.cmd.Terminate()
	return err
}

// TODO: implement this
func (p *cmdProvider) Hooks() {

}

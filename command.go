package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"text/template"

	"golang.org/x/crypto/ssh"

	"github.com/sgreben/sshtunnel"
	sshtunnelExec "github.com/sgreben/sshtunnel/exec"
)

const dockerComposeCmd = "docker-compose"

// CommandParallel executes the given docker-compose command in parallel on all projects
func (c *ConfigV1) CommandParallel(args []string) error {
	var errs Errors
	var errsMu sync.Mutex
	ps := c.ProjectsSorted()
	var wg sync.WaitGroup
	wg.Add(len(ps))
	for _, p := range ps {
		p := p
		go func() {
			defer wg.Done()
			if err := p.Command(args); err != nil {
				errsMu.Lock()
				errs = append(errs, fmt.Errorf("%s: %v", p.Name, err))
				errsMu.Unlock()
			}
		}()
	}
	wg.Wait()
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Command executes the given docker-compose command on all projects
func (c *ConfigV1) Command(args []string) error {
	var errs Errors
	ps := c.ProjectsSorted()
	for _, p := range ps {
		if err := p.Command(args); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// Command executes the given docker-compose command on this project
func (p *ConfigV1Project) Command(args []string) error {
	dir := filepath.Dir(p.ComposeFile)
	dockerHost, err := expandEnv(p.DockerHost)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tunnelErrCh := make(<-chan error)
	if p.DockerHostTunnel != nil {
		listener, errCh, err := p.DockerHostTunnel.Establish(ctx, dockerHost)
		if err != nil {
			return err
		}
		tunnelErrCh = errCh
		dockerHost = fmt.Sprintf("%s://%s", listener.Addr().Network(), listener.Addr().String())
	}
	cmdArgs := []string{
		"-f", filepath.Base(p.ComposeFile),
		"--project-name", p.Name,
	}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command(dockerComposeCmd, cmdArgs...)
	if dockerHost != "" {
		dockerHostEnv := fmt.Sprintf("DOCKER_HOST=%s", dockerHost)
		cmd.Env = append(os.Environ(), dockerHostEnv)
	}
	cmd.Dir, _ = filepath.Abs(dir)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	dirAbs, _ := filepath.Abs(cmd.Dir)
	if flags.Verbose {
		log.Printf("running in %q", dirAbs)
		log.Printf("exec: %q", cmd.Args)
	}
	cmdErrCh := make(chan error, 1)
	go func() {
		cmdErrCh <- cmd.Run()
	}()
	select {
	case cmdErr := <-cmdErrCh:
		return cmdErr
	case tunnelErr := <-tunnelErrCh:
		log.Printf("tunnel: %v", tunnelErr)
		return <-cmdErrCh
	}
}

// Establish establishes a connection to a remote Docker daemon at `addr`, and
// returns a local forwarding listener.
func (c *ConfigV1ProjectTunnel) Establish(ctx context.Context, addr string) (net.Listener, <-chan error, error) {
	switch {
	case c.SSH != nil:
		sshHost, err := expandEnv(c.SSH.Host)
		if err != nil {
			return nil, nil, err
		}
		userName, err := expandEnv(c.SSH.UserName)
		if err != nil {
			return nil, nil, err
		}
		if flags.SSHExternalClientOpenSSH {
			flags.SSHExternalClient = sshtunnelExec.CommandTemplateOpenSSHText
		}
		if flags.SSHExternalClientPuTTY {
			flags.SSHExternalClient = sshtunnelExec.CommandTemplatePuTTYText
		}
		if flags.SSHExternalClient != "" {
			commandTemplate, err := template.New("").Parse(flags.SSHExternalClient)
			if err != nil {
				return nil, nil, fmt.Errorf("parse ssh client command template: %v", err)
			}
			sshHostName, sshPort, err := net.SplitHostPort(sshHost)
			if err != nil {
				sshHostName, sshPort = sshHost, "22"
			}
			tunnelConfig := sshtunnelExec.Config{
				User:             userName,
				SSHHost:          sshHostName,
				SSHPort:          sshPort,
				CommandTemplate:  commandTemplate,
				CommandExtraArgs: flags.SSHExternalClientExtraArgs,
				Backoff:          flags.SSHReconnectBackoff,
			}
			listener, errCh, err := sshtunnelExec.ListenContext(
				ctx,
				&net.TCPAddr{IP: net.ParseIP("127.0.0.1")},
				addr,
				&tunnelConfig,
			)
			return listener, errCh, err
		}
		var clientConfig ssh.ClientConfig
		clientConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		ssh := c.SSH
		clientConfig.User = userName
		authConfig := &sshtunnel.ConfigAuth{}
		if ssh.UseAgent {
			agentAddr, err := expandEnv(flags.SSHAgentAddr)
			if err != nil {
				return nil, nil, err
			}
			authConfig.SSHAgent = &sshtunnel.ConfigSSHAgent{
				Addr: &net.UnixAddr{
					Net:  "unix",
					Name: agentAddr,
				},
			}
		}
		if ssh.Password != nil {
			password, err := expandEnv(*ssh.Password)
			if err != nil {
				return nil, nil, err
			}
			authConfig.Password = &password
		}
		if ssh.KeyFile != nil {
			path, err := expandEnv(*ssh.KeyFile)
			if err != nil {
				return nil, nil, err
			}
			var passphrase *[]byte
			if ssh.KeyPassphrase != nil {
				passphraseString, err := expandEnv(string(*ssh.KeyPassphrase))
				if err != nil {
					return nil, nil, err
				}
				passphraseBytes := []byte(passphraseString)
				passphrase = &passphraseBytes
			}
			authConfig.Keys = []sshtunnel.KeySource{
				{
					Path:       &path,
					Passphrase: passphrase,
				},
			}
		}
		clientConfig.Auth, err = authConfig.Methods()
		if err != nil {
			return nil, nil, err
		}
		tunnelConfig := sshtunnel.Config{
			SSHAddr:   sshHost,
			SSHClient: &clientConfig,
		}
		listener, errCh, err := sshtunnel.ListenContext(
			ctx,
			&net.TCPAddr{IP: net.ParseIP("127.0.0.1")},
			"unix", addr,
			&tunnelConfig,
			flags.SSHReconnectBackoff)
		return listener, errCh, err
	default:
		return nil, nil, fmt.Errorf("error: empty connection block")
	}
}

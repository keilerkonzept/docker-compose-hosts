package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
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
	cmdArgs := []string{
		"-f", filepath.Base(p.ComposeFile),
		"--project-name", p.Name,
	}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.CommandContext(ctx, dockerComposeCmd, cmdArgs...)
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
	if cmdErr := cmd.Run(); cmdErr != nil {
		return cmdErr
	}
	return nil
}

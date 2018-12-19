package main

import "sort"

// ConfigV1 is the configuration for a set of docker-compose projects
type ConfigV1 struct {
	Version  string
	Projects map[string]*ConfigV1Project
}

// ProjectsSorted returns a slice of projects lexicographically sorted by name
func (c *ConfigV1) ProjectsSorted() (out []*ConfigV1Project) {
	for n, p := range c.Projects {
		p.Name = n
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return
}

// ConfigV1Project is the configuration for a single docker-compose project
type ConfigV1Project struct {
	Name             string                 `yaml:"-"`
	DockerHost       string                 `yaml:"docker_host,omitempty"`
	DockerHostTunnel *ConfigV1ProjectTunnel `yaml:"docker_host_tunnel,omitempty"`
	ComposeFile      string                 `yaml:"compose_file"`
}

// ConfigV1ProjectTunnel is the configuration for a single docker-compose project's Docker host connection
type ConfigV1ProjectTunnel struct {
	SSH *ConfigV1ProjectTunnelSSH `yaml:"ssh,omitempty"`
}

// ConfigV1ProjectTunnelSSH is the configuration for an SSH connection to a Docker host
type ConfigV1ProjectTunnelSSH struct {
	Host          string  `yaml:"host"`
	HostKeyFile   *string `yaml:"host_key_file,omitempty"`
	UseAgent      bool    `yaml:"agent,omitempty"`
	UserName      string  `yaml:"user"`
	Password      *string `yaml:"password,omitempty"`
	KeyFile       *string `yaml:"key_file,omitempty"`
	KeyPassphrase *[]byte `yaml:"key_passphrase,omitempty"`
}

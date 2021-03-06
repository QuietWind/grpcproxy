package config

import (
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

func ReadConfig(filename string) (ServerConfig, error) {
	cfg := ServerConfig{}
	return cfg, cfg.Read(filename)
}

type ServerConfig struct {
	Bind []string                `hcl:"bind,omitempty" json:"bind,omitempty"`
	Cert []string                `hcl:"cert,omitempty" json:"cert,omitempty"`
	CA   []string                `hcl:"ca" json:"ca"`
	GRPC bool                    `hcl:"grpc,omitempty" json:"grpc,omitempty"`
	AppM []map[string]*AppConfig `hcl:"app,omitempty" json:"app,omitempty"`
	App  []*AppConfig            `hcl:"-" json:"-"`
}

func (this *ServerConfig) Read(filename string) error {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := hcl.Unmarshal(in, this); err != nil {
		return err
	}

	this.Init()

	return nil
}

func (this *ServerConfig) Init() {
	for _, m := range this.AppM {
		for name, app := range m {
			app.link()
			app.server = this
			app.Name = name
			if app.Host == "" {
				app.Host = name
			}

			this.App = append(this.App, app)
		}
	}
}

type AppConfig struct {
	server *ServerConfig

	Name   string                    `hcl:"-" json:"-"`
	Host   string                    `hcl:"host,omitempty" json:"host,omitempty"`
	GRPC   *bool                     `hcl:"grpc,omitempty" json:"grpc,omitempty"`
	CA     []string                  `hcl:"ca" json:"ca"`
	ProxyM []map[string]*ProxyConfig `hcl:"proxy,omitempty" json:"proxy,omitempty"`
	Proxy  []*ProxyConfig            `hcl:"-" json:"-"`
}

func (this *AppConfig) GetGRPC() bool {
	if this.GRPC == nil {
		return this.server.GRPC
	}

	return *this.GRPC
}

func (this *AppConfig) GetCA() []string {
	if len(this.CA) == 0 {
		return this.server.CA
	}

	return this.CA
}

func (this *AppConfig) link() {
	for _, m := range this.ProxyM {
		for name, proxy := range m {
			proxy.link()
			proxy.app = this
			proxy.Name = name
			if proxy.URI == "" {
				proxy.URI = name
			}

			this.Proxy = append(this.Proxy, proxy)
		}
	}
}

type ProxyConfig struct {
	app *AppConfig

	Name               string   `hcl:"-" json:"-"`
	URI                string   `hcl:"uri,omitempty" json:"uri,omitempty"`
	Host               string   `hcl:"host,omitempty" json:"host,omitempty"`
	GRPC               *bool    `hcl:"grpc,omitempty" json:"grpc,omitempty"`
	Backend            string   `hcl:"backend,omitempty" json:"backend,omitempty"`
	Policy             string   `hcl:"policy,omitempty" json:"policy,omitempty"`
	CA                 []string `hcl:"ca" json:"ca"`
	TLS                bool     `hcl:"tls,omitempty" json:"tls,omitempty"`
	InsecureSkipVerify bool     `hcl:"insecure_skip_verify,omitempty" json:"insecure_skip_verify,omitempty"`
}

func (this *ProxyConfig) GetGRPC() bool {
	if this.GRPC == nil {
		return this.app.GetGRPC()
	}

	return *this.GRPC
}

func (this *ProxyConfig) GetCA() []string {
	if len(this.CA) == 0 {
		return this.app.GetCA()
	}

	return this.CA
}

func (this *ProxyConfig) link() {

}

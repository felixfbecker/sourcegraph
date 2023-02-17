package servegit

import (
	"context"
	"net/url"
	"os"
	"strings"

	"github.com/sourcegraph/sourcegraph/internal/debugserver"
	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/internal/service"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

type Config struct {
	env.BaseConfig

	Addr       string
	ReposRoots []string
}

func (c *Config) Load() {
	defaultReposRoot := ""
	if pwd, err := os.Getwd(); err == nil {
		defaultReposRoot = pwd
	}

	roots := strings.Fields(c.Get("SRC", defaultReposRoot, "Root dir containing repos."))
	c.ReposRoots = roots

	url, err := url.Parse(c.Get("SRC_SERVE_GIT_URL", "http://127.0.0.1:3434", "URL that servegit should listen on."))
	if err != nil {
		c.AddError(errors.Wrapf(err, "failed to parse SRC_SERVE_GIT_URL"))
	} else if url.Scheme != "http" {
		c.AddError(errors.Errorf("only support http scheme for SRC_SERVE_GIT_URL got scheme %q", url.Scheme))
	} else {
		c.Addr = url.Host
	}
}

type svc struct{}

func (s svc) Name() string {
	return "servegit"
}

func (s svc) Configure() (env.Config, []debugserver.Endpoint) {
	c := &Config{}
	c.Load()
	return c, nil
}

func (s svc) Start(ctx context.Context, observationCtx *observation.Context, ready service.ReadyFunc, configI env.Config) (err error) {
	defer func() {
		if err == nil {
			ready()
		}
	}()

	config := configI.(*Config)

	// Start servegit which walks ReposRoot to find repositories and exposes
	// them over HTTP for Sourcegraph's syncer to discover and clone.
	srv := &Serve{
		Addr:   config.Addr,
		Logger: observationCtx.Logger,
	}
	if err := srv.Start(); err != nil {
		return errors.Wrap(err, "failed to start servegit server which discovers local repositories")
	}

	// Now that servegit is running, we can add the external service which
	// connects to it.
	//
	// Note: src.Addr is updated to reflect the actual listening address.
	if err := ensureExtSVC(observationCtx, "http://"+srv.Addr, config.ReposRoots); err != nil {
		return errors.Wrap(err, "failed to create external service which imports local repositories")
	}

	return nil
}

var Service service.Service = svc{}

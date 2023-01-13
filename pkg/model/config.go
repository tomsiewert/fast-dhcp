package model

import (
	"github.com/tomsiewert/fast-dhcp/pkg/dhcp"
	"github.com/tomsiewert/fast-dhcp/pkg/prometheus"
)

type Config struct {
	SentryDSN   string            `json:"sentry_dsn"`
	PProfListen string            `json:"pprof_listen"`
	DHCP        dhcp.Config       `json:"dhcp"`
	Prometheus  prometheus.Config `json:"prometheus"`
}

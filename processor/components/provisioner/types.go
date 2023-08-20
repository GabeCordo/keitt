package provisioner

import (
	"github.com/GabeCordo/mango-go/processor/components/supervisor"
	"github.com/GabeCordo/mango/components/cluster"
	"sync"
)

const (
	DefaultFrameworkModule = "common"
)

type ClusterWrapper struct {
	registry *supervisor.Registry

	Identifier        string          `json:"identifier"`
	Mode              cluster.EtlMode `json:"mode"`
	Mounted           bool            `json:"mounted"`
	MarkedForDeletion bool            `json:"marked-for-deletion"`
	DefaultConfig     cluster.Config  `json:"default-config"`

	mutex sync.RWMutex
}

type ModuleWrapper struct {
	clusters map[string]*ClusterWrapper

	Mounted         bool `json:"mounted"`
	MarkForDeletion bool `json:"mark-for-deletion"`

	Identifier string  `json:"identifier"`
	Version    float64 `json:"version"`

	mutex sync.RWMutex
}

type Provisioner struct {
	modules map[string]*ModuleWrapper
	mutex   sync.RWMutex
}

package backends

import (
	"fmt"

	"github.com/varsitytutors/gostatsd"
	"github.com/varsitytutors/gostatsd/pkg/backends/cloudwatch"
	"github.com/varsitytutors/gostatsd/pkg/backends/datadog"
	"github.com/varsitytutors/gostatsd/pkg/backends/graphite"
	"github.com/varsitytutors/gostatsd/pkg/backends/newrelic"
	"github.com/varsitytutors/gostatsd/pkg/backends/null"
	"github.com/varsitytutors/gostatsd/pkg/backends/statsdaemon"
	"github.com/varsitytutors/gostatsd/pkg/backends/stdout"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// All known backends.
var backends = map[string]gostatsd.BackendFactory{
	datadog.BackendName:     datadog.NewClientFromViper,
	graphite.BackendName:    graphite.NewClientFromViper,
	null.BackendName:        null.NewClientFromViper,
	statsdaemon.BackendName: statsdaemon.NewClientFromViper,
	stdout.BackendName:      stdout.NewClientFromViper,
	cloudwatch.BackendName:  cloudwatch.NewClientFromViper,
	newrelic.BackendName:    newrelic.NewClientFromViper,
}

// GetBackend creates an instance of the named backend, or nil if
// the name is not known. The error return is only used if the named backend
// was known but failed to initialize.
func GetBackend(name string, v *viper.Viper) (gostatsd.Backend, error) {
	f, found := backends[name]
	if !found {
		return nil, nil
	}
	return f(v)
}

// InitBackend creates an instance of the named backend.
func InitBackend(name string, v *viper.Viper) (gostatsd.Backend, error) {
	if name == "" {
		log.Info("No backend specified")
		return nil, nil
	}

	backend, err := GetBackend(name, v)
	if err != nil {
		return nil, fmt.Errorf("could not init backend %q: %v", name, err)
	}
	if backend == nil {
		return nil, fmt.Errorf("unknown backend %q", name)
	}
	log.Infof("Initialised backend %q", name)

	return backend, nil
}

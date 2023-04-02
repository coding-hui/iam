package options

import (
	"flag"

	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"

	"github.com/wecoding/iam/pkg/apiserver/config"
)

// ServerRunOptions contains everything necessary to create and run api apiserver
type ServerRunOptions struct {
	GenericServerRunOptions *config.Config
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters
func NewServerRunOptions() *ServerRunOptions {
	s := &ServerRunOptions{
		GenericServerRunOptions: config.NewConfig(),
	}
	return s
}

// Flags returns the complete NamedFlagSets
func (s *ServerRunOptions) Flags() (fss cliflag.NamedFlagSets) {
	fs := fss.FlagSet("generic")
	s.GenericServerRunOptions.AddFlags(fs, s.GenericServerRunOptions)
	local := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(local)
	fs.AddGoFlagSet(local)
	return fss
}

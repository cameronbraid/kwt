package core

import (
	"fmt"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ConfigFactory interface {
	ConfigurePathResolver(func() (string, error))
	ConfigureContextResolver(func() (string, error))
	RESTConfig() (*rest.Config, error)
	DefaultNamespace() (string, error)
}

type ConfigFactoryImpl struct {
	pathResolverFunc    func() (string, error)
	contextResolverFunc func() (string, error)
}

var _ ConfigFactory = &ConfigFactoryImpl{}

func NewConfigFactoryImpl() *ConfigFactoryImpl {
	return &ConfigFactoryImpl{}
}

func (f *ConfigFactoryImpl) ConfigurePathResolver(resolverFunc func() (string, error)) {
	f.pathResolverFunc = resolverFunc
}

func (f *ConfigFactoryImpl) ConfigureContextResolver(resolverFunc func() (string, error)) {
	f.contextResolverFunc = resolverFunc
}

func (f *ConfigFactoryImpl) RESTConfig() (*rest.Config, error) {
	config, err := f.clientConfig()
	if err != nil {
		return nil, err
	}

	restConfig, err := config.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("Building Kubernetes config: %s", err)
	}

	return restConfig, nil
}

func (f *ConfigFactoryImpl) DefaultNamespace() (string, error) {
	config, err := f.clientConfig()
	if err != nil {
		return "", err
	}

	name, _, err := config.Namespace()
	return name, err
}

func (f *ConfigFactoryImpl) clientConfig() (clientcmd.ClientConfig, error) {
	path, err := f.pathResolverFunc()
	if err != nil {
		return nil, fmt.Errorf("Resolving config path: %s", err)
	}

	context, err := f.contextResolverFunc()
	if err != nil {
		return nil, fmt.Errorf("Resolving config context: %s", err)
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: path},
		&clientcmd.ConfigOverrides{CurrentContext: context},
	), nil
}

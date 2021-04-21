package kubeconfig

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Kubeconfig struct {
	clientConfig clientcmd.ClientConfig
}

func New() *Kubeconfig {
	return &Kubeconfig{
		clientConfig: clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientcmd.NewDefaultClientConfigLoadingRules(), &clientcmd.ConfigOverrides{}),
	}
}

func (k *Kubeconfig) ReadCurrentContext() (*api.Context, error) {
	rawConfig, err := k.clientConfig.RawConfig()
	if err != nil {
		return nil, err
	}

	return k.ReadContext(rawConfig.CurrentContext)
}

func (k *Kubeconfig) ReadContext(name string) (*api.Context, error) {
	rawConfig, err := k.clientConfig.RawConfig()
	if err != nil {
		return nil, err
	}

	obj, ok := rawConfig.Contexts[name]
	if !ok {
		return nil, fmt.Errorf("'%s' context can not read", name)
	}
	if obj == nil {
		return nil, fmt.Errorf("'%s' context was nof found in your kubeconfig", name)
	}

	return obj, nil
}

func (k *Kubeconfig) ReadUser(name string) (*api.AuthInfo, error) {
	rawConfig, err := k.clientConfig.RawConfig()
	if err != nil {
		return nil, err
	}

	obj, ok := rawConfig.AuthInfos[name]
	if !ok {
		return nil, fmt.Errorf("'%s' user was nof found in your kubeconfig", name)
	}

	return obj, nil
}

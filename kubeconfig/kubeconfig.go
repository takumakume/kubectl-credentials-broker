package kubeconfig

import (
	"fmt"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Kubeconfig struct {
	clientConfig   clientcmd.ClientConfig
	configFilePath string
}

type Credential struct {
	ClientCertificate string
	ClientKey         string
	Token             string
}

func New() *Kubeconfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})

	return &Kubeconfig{
		clientConfig:   clientConfig,
		configFilePath: rules.GetLoadingPrecedence()[0],
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
		return nil, fmt.Errorf("'%s' context was not found in your kubeconfig", name)
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
		return nil, fmt.Errorf("'%s' user was not found in your kubeconfig", name)
	}

	return obj, nil
}

func (k *Kubeconfig) ReadCluster(name string) (*api.Cluster, error) {
	rawConfig, err := k.clientConfig.RawConfig()
	if err != nil {
		return nil, err
	}

	obj, ok := rawConfig.Clusters[name]
	if !ok {
		return nil, fmt.Errorf("'%s' cluster was not found in your kubeconfig", name)
	}

	return obj, nil
}

func (k *Kubeconfig) ReadCurrentUserExecVersion() (string, error) {
	cc, err := k.ReadCurrentContext()
	if err != nil {
		return "", err
	}

	user, err := k.ReadUser(cc.AuthInfo)
	if err != nil {
		return "", err
	}

	if user.Exec == nil {
		return "", fmt.Errorf("exec is not specified for user in current-context, this command expects to be run as a credential plugin for kubeconfig")
	}

	return user.Exec.APIVersion, nil
}

func (k *Kubeconfig) UpdateCurrentUserExecConfig(apiVersion, cmd string, args []string, envs map[string]string) error {
	cc, err := k.ReadCurrentContext()
	if err != nil {
		return err
	}

	rawConfig, err := k.clientConfig.RawConfig()
	if err != nil {
		return err
	}

	if rawConfig.AuthInfos[cc.AuthInfo] == nil {
		return fmt.Errorf("'%s' user was not found in your kubeconfig", cc.AuthInfo)
	}

	envVars := []api.ExecEnvVar{}
	for name, value := range envs {
		envVars = append(envVars, api.ExecEnvVar{
			Name:  name,
			Value: value,
		})
	}

	rawConfig.AuthInfos[cc.AuthInfo].Exec = &api.ExecConfig{
		APIVersion: apiVersion,
		Command:    cmd,
		Args:       args,
		Env:        envVars,
	}

	rawConfig.AuthInfos[cc.AuthInfo].Token = ""

	if err := k.write(rawConfig); err != nil {
		return err
	}

	return nil
}

func (k *Kubeconfig) write(rawConfig api.Config) error {
	if err := clientcmd.Validate(rawConfig); err != nil {
		return err
	}

	if err := clientcmd.WriteToFile(rawConfig, k.configFilePath); err != nil {
		return err
	}

	return nil
}

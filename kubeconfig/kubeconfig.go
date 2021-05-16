package kubeconfig

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sergi/go-diff/diffmatchpatch"
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

func (k *Kubeconfig) updateCurrentUserExecConfig(apiVersion, cmd string, args []string, envs map[string]string) (api.Config, error) {
	cc, err := k.ReadCurrentContext()
	if err != nil {
		return api.Config{}, err
	}

	rawConfig, err := k.clientConfig.RawConfig()
	if err != nil {
		return api.Config{}, err
	}

	if rawConfig.AuthInfos[cc.AuthInfo] == nil {
		return api.Config{}, fmt.Errorf("'%s' user was not found in your kubeconfig", cc.AuthInfo)
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

	return rawConfig, nil
}

func (k *Kubeconfig) UpdateCurrentUserExecConfig(apiVersion, cmd string, args []string, envs map[string]string) error {
	newApiConfig, err := k.updateCurrentUserExecConfig(apiVersion, cmd, args, envs)
	if err != nil {
		return err
	}

	return k.update(newApiConfig)
}

func (k *Kubeconfig) UpdateCurrentUserExecConfigDryRun(apiVersion, cmd string, args []string, envs map[string]string) (string, error) {
	newApiConfig, err := k.updateCurrentUserExecConfig(apiVersion, cmd, args, envs)
	if err != nil {
		return "", err
	}

	return k.diff(newApiConfig)
}

func (k *Kubeconfig) write(path string, rawConfig api.Config) error {
	if err := clientcmd.Validate(rawConfig); err != nil {
		return err
	}

	if err := clientcmd.WriteToFile(rawConfig, path); err != nil {
		return err
	}

	return nil
}

func (k *Kubeconfig) update(rawConfig api.Config) error {
	if err := k.write(k.configFilePath, rawConfig); err != nil {
		return err
	}

	return nil
}

func (k *Kubeconfig) diff(newRawConfig api.Config) (string, error) {
	tmpfile, err := ioutil.TempFile("", "kubectl-credentials-broker-diff-tempfile")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name())

	if err := k.write(tmpfile.Name(), newRawConfig); err != nil {
		return "", err
	}

	oldConfig, err := ioutil.ReadFile(k.configFilePath)
	if err != nil {
		return "", err
	}

	newConfig, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		return "", err
	}

	return diff(string(oldConfig), string(newConfig)), nil
}

func diff(old, new string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(old, new, false)
	for _, diff := range diffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			return dmp.DiffPrettyText(diffs)
		}
	}

	return ""
}

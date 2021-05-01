package kubeconfig

import (
	"fmt"
	"io/ioutil"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Kubeconfig struct {
	clientConfig clientcmd.ClientConfig
}

type CertificateBundle struct {
	Certificate string
	Key         string
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

func (k *Kubeconfig) CurrentCertificateBundle() (*CertificateBundle, error) {
	cc, err := k.ReadCurrentContext()
	if err != nil {
		return nil, err
	}

	user, err := k.ReadUser(cc.AuthInfo)
	if err != nil {
		return nil, err
	}

	cluster, err := k.ReadCluster(cc.Cluster)
	if err != nil {
		return nil, err
	}

	var cert string
	var key string
	var ca string

	if len(user.ClientCertificateData) > 0 && len(user.ClientKeyData) > 0 {
		cert = string(user.ClientCertificateData)
		key = string(user.ClientKeyData)
	} else if len(user.ClientCertificate) > 0 && len(user.ClientKey) > 0 {
		certBuf, err := ioutil.ReadFile(user.ClientCertificate)
		if err != nil {
			return nil, err
		}
		cert = string(certBuf)

		keyBuf, err := ioutil.ReadFile(user.ClientKey)
		if err != nil {
			return nil, err
		}
		key = string(keyBuf)
	}

	if len(cluster.CertificateAuthorityData) > 0 {
		ca = string(cluster.CertificateAuthorityData)
	} else if len(cluster.CertificateAuthority) > 0 {
		caBuf, err := ioutil.ReadFile(cluster.CertificateAuthority)
		if err != nil {
			return nil, err
		}
		ca = string(caBuf)
	}

	if len(ca) > 0 {
		cert = fmt.Sprintf("%s\n%s", cert, ca)
	}

	certificateBundle := &CertificateBundle{
		Certificate: cert,
		Key:         key,
	}

	return certificateBundle, nil
}

func (k *Kubeconfig) CurrentUserToken() (string, error) {
	cc, err := k.ReadCurrentContext()
	if err != nil {
		return "", err
	}

	user, err := k.ReadUser(cc.AuthInfo)
	if err != nil {
		return "", err
	}

	if len(user.Token) > 0 {
		return user.Token, nil
	}

	return "", nil
}

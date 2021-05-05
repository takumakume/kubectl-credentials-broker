package credentials

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientauthenticationv1alpha1 "k8s.io/client-go/pkg/apis/clientauthentication/v1alpha1"
)

type V1Alpha1 struct{}

func (o *V1Alpha1) APIVersionString() string {
	return "client.authentication.k8s.io/v1alpha1"
}

func (o *V1Alpha1) ToJSON(opts *CredentialOption) ([]byte, error) {
	status := &clientauthenticationv1alpha1.ExecCredentialStatus{}
	if len(opts.ClientCertificateData) > 0 {
		status.ClientCertificateData = opts.ClientCertificateData
	}
	if len(opts.ClientKeyData) > 0 {
		status.ClientKeyData = opts.ClientKeyData
	}
	if len(opts.Token) > 0 {
		status.Token = opts.Token
	}

	return json.Marshal(&clientauthenticationv1alpha1.ExecCredential{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "client.authentication.k8s.io/v1alpha1",
			Kind:       "ExecCredential",
		},
		Status: status,
	})
}

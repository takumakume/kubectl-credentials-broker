package credentials

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientauthenticationv1beta1 "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

type V1Beta1 struct{}

func (o *V1Beta1) ToJSON(opts *CredentialOption) ([]byte, error) {
	status := &clientauthenticationv1beta1.ExecCredentialStatus{}
	if len(opts.ClientCertificateData) > 0 {
		status.ClientCertificateData = opts.ClientCertificateData
	}
	if len(opts.ClientKeyData) > 0 {
		status.ClientKeyData = opts.ClientKeyData
	}
	if len(opts.Token) > 0 {
		status.Token = opts.Token
	}

	return json.Marshal(&clientauthenticationv1beta1.ExecCredential{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "client.authentication.k8s.io/v1beta1",
			Kind:       "ExecCredential",
		},
		Status: status,
	})
}

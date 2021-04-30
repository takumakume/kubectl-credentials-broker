package kubeconfig

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"text/template"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestKubeconfig_ReadCurrentContext(t *testing.T) {
	type fields struct {
		kubeconfigString string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *api.Context
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    token: hoge`,
			},
			want: &api.Context{
				Cluster:    "server1",
				AuthInfo:   "user1",
				Namespace:  "kube-system",
				Extensions: map[string]runtime.Object{},
			},
		},
		{
			name: "no current config",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    token: hoge`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(tt.fields.kubeconfigString))
			if err != nil {
				t.Errorf("Kubeconfig.ReadCurrentContext() test data error = %+v", err)
				return
			}
			k := &Kubeconfig{
				clientConfig: clientConfig,
			}
			got, err := k.ReadCurrentContext()
			if (err != nil) != tt.wantErr {
				t.Errorf("Kubeconfig.ReadCurrentContext() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kubeconfig.ReadCurrentContext() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestKubeconfig_ReadContext(t *testing.T) {
	type fields struct {
		kubeconfigString string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *api.Context
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    token: hoge`,
			},
			args: args{
				name: "context1",
			},
			want: &api.Context{
				Cluster:    "server1",
				AuthInfo:   "user1",
				Namespace:  "kube-system",
				Extensions: map[string]runtime.Object{},
			},
		},
		{
			name: "not found",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    token: hoge`,
			},
			args: args{
				name: "notfound",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(tt.fields.kubeconfigString))
			if err != nil {
				t.Errorf("Kubeconfig.ReadContext() test data error = %+v", err)
				return
			}
			k := &Kubeconfig{
				clientConfig: clientConfig,
			}
			got, err := k.ReadContext(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kubeconfig.ReadContext() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kubeconfig.ReadContext() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestKubeconfig_ReadUser(t *testing.T) {
	type fields struct {
		kubeconfigString string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *api.AuthInfo
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    token: hoge`,
			},
			args: args{
				name: "user1",
			},
			want: &api.AuthInfo{
				Token:      "hoge",
				Extensions: map[string]runtime.Object{},
			},
		},
		{
			name: "not found",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    token: hoge`,
			},
			args: args{
				name: "notfound",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(tt.fields.kubeconfigString))
			if err != nil {
				t.Errorf("Kubeconfig.ReadUser() test data error = %+v", err)
				return
			}
			k := &Kubeconfig{
				clientConfig: clientConfig,
			}
			got, err := k.ReadUser(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kubeconfig.ReadUser() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kubeconfig.ReadUser() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestKubeconfig_ReadCluster(t *testing.T) {
	type fields struct {
		kubeconfigString string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *api.Cluster
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    token: hoge`,
			},
			args: args{
				name: "server1",
			},
			want: &api.Cluster{
				Server:     "https://127.0.0.1",
				Extensions: map[string]runtime.Object{},
			},
		},
		{
			name: "not found",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    token: hoge`,
			},
			args: args{
				name: "notfound",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(tt.fields.kubeconfigString))
			if err != nil {
				t.Errorf("Kubeconfig.ReadCluster() test data error = %+v", err)
				return
			}
			k := &Kubeconfig{
				clientConfig: clientConfig,
			}
			got, err := k.ReadCluster(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kubeconfig.ReadCluster() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kubeconfig.ReadCluster() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestKubeconfig_ReadCurrentUserExecVersion(t *testing.T) {
	type fields struct {
		kubeconfigString string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1`,
			},
			want: "client.authentication.k8s.io/v1beta1",
		},
		{
			name: "not found",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(tt.fields.kubeconfigString))
			if err != nil {
				t.Errorf("Kubeconfig.ReadCurrentUserExecVersion() test data error = %+v", err)
				return
			}
			k := &Kubeconfig{
				clientConfig: clientConfig,
			}
			got, err := k.ReadCurrentUserExecVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("Kubeconfig.ReadCurrentUserExecVersion() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Kubeconfig.ReadCurrentUserExecVersion() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestKubeconfig_CurrentCertificateBundle(t *testing.T) {
	testDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("TempDir() error = %v", err)
		return
	}
	defer os.Remove(testDir)

	certFile, err := ioutil.TempFile(testDir, "tls.crt")
	if _, err = certFile.Write([]byte("client-certificate-file-sample-string")); err != nil {
		t.Errorf("ioutil.Write() error = %v", err)
		return
	}

	keyFile, err := ioutil.TempFile(testDir, "tls.key")
	if _, err = keyFile.Write([]byte("client-key-file-sample-string")); err != nil {
		t.Errorf("ioutil.Write() error = %v", err)
		return
	}

	caFile, err := ioutil.TempFile(testDir, "ca.crt")
	if _, err = caFile.Write([]byte("certificate-authority-file-sample-string")); err != nil {
		t.Errorf("ioutil.Write() error = %v", err)
		return
	}

	type fields struct {
		kubeconfigString string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *CertificateBundle
		wantErr bool
	}{
		{
			name: "from files",
			fields: fields{
				kubeconfigString: tmpl(`---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
    certificate-authority: {{ .ca }}
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    client-certificate: {{ .cert }}
    client-key: {{ .key }}`, map[string]interface{}{
					"ca":   caFile.Name(),
					"cert": certFile.Name(),
					"key":  keyFile.Name()}),
			},
			want: &CertificateBundle{
				Certificate: "client-certificate-file-sample-string",
				Key:         "client-key-file-sample-string",
				CA:          "certificate-authority-file-sample-string",
			},
		},
		{
			name: "*-data used",
			fields: fields{
				kubeconfigString: tmpl(`---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
    certificate-authority: {{ .ca }}
    certificate-authority-data: Y2VydGlmaWNhdGUtYXV0aG9yaXR5LWRhdGEtc2FtcGxlLXN0cmluZw== # certificate-authority-data-sample-string
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    client-certificate: {{ .cert }}
    client-certificate-data: Y2xpZW50LWNlcnRpZmljYXRlLWRhdGEtc2FtcGxlLXN0cmluZw== # client-certificate-data-sample-string
    client-key: {{ .key }}
    client-key-data: Y2xpZW50LWtleS1kYXRhLXNhbXBsZS1zdHJpbmc= #client-key-data-sample-string`, map[string]interface{}{
					"ca":   caFile.Name(),
					"cert": certFile.Name(),
					"key":  keyFile.Name()}),
			},
			want: &CertificateBundle{
				Certificate: "client-certificate-data-sample-string",
				Key:         "client-key-data-sample-string",
				CA:          "certificate-authority-data-sample-string",
			},
		},
		{
			name: "*-data only",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
    certificate-authority-data: Y2VydGlmaWNhdGUtYXV0aG9yaXR5LWRhdGEtc2FtcGxlLXN0cmluZw== # certificate-authority-data-sample-string
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    client-certificate-data: Y2xpZW50LWNlcnRpZmljYXRlLWRhdGEtc2FtcGxlLXN0cmluZw== # client-certificate-data-sample-string
    client-key-data: Y2xpZW50LWtleS1kYXRhLXNhbXBsZS1zdHJpbmc= #client-key-data-sample-string`,
			},
			want: &CertificateBundle{
				Certificate: "client-certificate-data-sample-string",
				Key:         "client-key-data-sample-string",
				CA:          "certificate-authority-data-sample-string",
			},
		},
		{
			name: "all empty",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:`,
			},
			want: &CertificateBundle{
				Certificate: "",
				Key:         "",
				CA:          "",
			},
		},
		{
			name: "client-certificate-data only (empty if client-key-data is missing)",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    client-certificate-data: Y2xpZW50LWNlcnRpZmljYXRlLWRhdGEtc2FtcGxlLXN0cmluZw== # client-certificate-data-sample-string`,
			},
			want: &CertificateBundle{
				Certificate: "",
				Key:         "",
				CA:          "",
			},
		},
		{
			name: "client-certificate only (empty if client-key is missing)",
			fields: fields{
				kubeconfigString: tmpl(`---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    client-certificate: {{ .cert }}`, map[string]interface{}{
					"cert": certFile.Name()}),
			},
			want: &CertificateBundle{
				Certificate: "",
				Key:         "",
				CA:          "",
			},
		},
		{
			name: "file not found",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    client-certificate: /path/to/notfound.crt
    client-key: /path/to/notfound.key`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(tt.fields.kubeconfigString))
			if err != nil {
				t.Errorf("Kubeconfig.CurrentCertificateBundle() test data error = %+v", err)
				return
			}
			k := &Kubeconfig{
				clientConfig: clientConfig,
			}
			got, err := k.CurrentCertificateBundle()
			if (err != nil) != tt.wantErr {
				t.Errorf("Kubeconfig.CurrentCertificateBundle() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Kubeconfig.CurrentCertificateBundle() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestKubeconfig_CurrentUserToken(t *testing.T) {
	type fields struct {
		kubeconfigString string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1
  user:
    token: hoge`,
			},
			want: "hoge",
		},
		{
			name: "not found",
			fields: fields{
				kubeconfigString: `---
apiVersion: v1
kind: Config
current-context: context1
clusters:
- cluster:
    server: https://127.0.0.1
  name: server1
contexts:
- context:
    cluster: server1
    namespace: kube-system
    user: user1
  name: context1
users:
- name: user1`,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(tt.fields.kubeconfigString))
			if err != nil {
				t.Errorf("Kubeconfig.CurrentUserToken() test data error = %+v", err)
				return
			}
			k := &Kubeconfig{
				clientConfig: clientConfig,
			}
			got, err := k.CurrentUserToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("Kubeconfig.CurrentUserToken() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Kubeconfig.CurrentUserToken() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func tmpl(tpl string, params map[string]interface{}) string {
	var tmpl = template.Must(template.New("").Parse(tpl))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		panic(err)
	}

	return buf.String()
}

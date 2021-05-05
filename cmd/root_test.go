package cmd

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/takumakume/kubectl-credentials-broker/credentials"
	"github.com/takumakume/kubectl-credentials-broker/kubeconfig"
)

func Test_rootCmdArgs_validate(t *testing.T) {
	type fields struct {
		clientCertificatePath string
		clientKeyPath         string
		tokenPath             string
		beforeExecCommand     string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				clientCertificatePath: "/path/to/tls.crt",
				clientKeyPath:         "/path/to/tls.key",
				tokenPath:             "/path/to/token",
				beforeExecCommand:     "/path/to/script.sh",
			},
		},
		{
			name: "certificate and key only",
			fields: fields{
				clientCertificatePath: "/path/to/tls.crt",
				clientKeyPath:         "/path/to/tls.key",
				tokenPath:             "",
				beforeExecCommand:     "",
			},
		},
		{
			name: "token only",
			fields: fields{
				clientCertificatePath: "",
				clientKeyPath:         "",
				tokenPath:             "/path/to/token",
				beforeExecCommand:     "",
			},
		},
		{
			name: "both client-certificate-path and client-key-path must be provided (only client-certificate-path)",
			fields: fields{
				clientCertificatePath: "/path/to/tls.crt",
				clientKeyPath:         "",
				tokenPath:             "",
				beforeExecCommand:     "",
			},
			wantErr: true,
		},
		{
			name: "both client-certificate-path and client-key-path must be provided (only client-key-path)",
			fields: fields{
				clientCertificatePath: "",
				clientKeyPath:         "/path/to/tls.key",
				tokenPath:             "",
				beforeExecCommand:     "",
			},
			wantErr: true,
		},
		{
			name: "requires either certificate token",
			fields: fields{
				clientCertificatePath: "",
				clientKeyPath:         "",
				tokenPath:             "",
				beforeExecCommand:     "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &rootCmdArgs{
				clientCertificatePath: tt.fields.clientCertificatePath,
				clientKeyPath:         tt.fields.clientKeyPath,
				tokenPath:             tt.fields.tokenPath,
				beforeExecCommand:     tt.fields.beforeExecCommand,
			}
			if err := args.validate(); (err != nil) != tt.wantErr {
				t.Errorf("rootCmdArgs.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRun(t *testing.T) {
	testDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("TempDir() error = %v", err)
		return
	}
	defer os.Remove(testDir)

	tokenFile, err := ioutil.TempFile(testDir, "token")
	if _, err = tokenFile.Write([]byte("token-from-file")); err != nil {
		t.Errorf("ioutil.Write() error = %v", err)
		return
	}

	type args struct {
		rootCmdArgs    rootCmdArgs
		kubeconfigData string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "API version: client.authentication.k8s.io/v1beta1",
			args: args{
				rootCmdArgs: rootCmdArgs{
					tokenPath: tokenFile.Name(),
				},
				kubeconfigData: `---
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
    namespace: default
    user: user1
  name: context1
users:
- name: user1
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1`,
			},
			want: []byte(`{"kind":"ExecCredential","apiVersion":"client.authentication.k8s.io/v1beta1","spec":{},"status":{"token":"token-from-file"}}`),
		},
		{
			name: "API version: client.authentication.k8s.io/v1alpha1",
			args: args{
				rootCmdArgs: rootCmdArgs{
					tokenPath: tokenFile.Name(),
				},
				kubeconfigData: `---
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
    namespace: default
    user: user1
  name: context1
users:
- name: user1
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1`,
			},
			want: []byte(`{"kind":"ExecCredential","apiVersion":"client.authentication.k8s.io/v1alpha1","spec":{},"status":{"token":"token-from-file"}}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kubeconfigFile, err := ioutil.TempFile("", "kube-config-")
			defer os.Remove(kubeconfigFile.Name())

			if _, err = kubeconfigFile.Write([]byte(tt.args.kubeconfigData)); err != nil {
				t.Errorf("kubeconfigFile.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Run error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			os.Setenv("KUBECONFIG", kubeconfigFile.Name())

			runner, err := newRootCmdRunner(&tt.args.rootCmdArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run newRunner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := runner.run()
			if (err != nil) != tt.wantErr {
				t.Errorf("Run newRunner().run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run newRunner().run() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_makeCredentialOptions(t *testing.T) {
	testDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("TempDir() error = %v", err)
		return
	}
	defer os.Remove(testDir)

	certFile, err := ioutil.TempFile(testDir, "tls.crt")
	if _, err = certFile.Write([]byte("client-certificate-from-file")); err != nil {
		t.Errorf("ioutil.Write() error = %v", err)
		return
	}

	keyFile, err := ioutil.TempFile(testDir, "tls.key")
	if _, err = keyFile.Write([]byte("client-key-from-file")); err != nil {
		t.Errorf("ioutil.Write() error = %v", err)
		return
	}

	caFile, err := ioutil.TempFile(testDir, "ca.crt")
	if _, err = caFile.Write([]byte("certificate-authority-from-file")); err != nil {
		t.Errorf("ioutil.Write() error = %v", err)
		return
	}

	tokenFile, err := ioutil.TempFile(testDir, "token")
	if _, err = tokenFile.Write([]byte("token-from-file")); err != nil {
		t.Errorf("ioutil.Write() error = %v", err)
		return
	}

	type args struct {
		args                 *rootCmdArgs
		kubeConfigCredential *kubeconfig.Credential
	}
	tests := []struct {
		name    string
		args    args
		want    *credentials.CredentialOption
		wantErr bool
	}{
		{
			name: "args certificate/key and token overrides kubeconfig value",
			args: args{
				args: &rootCmdArgs{
					clientCertificatePath: certFile.Name(),
					clientKeyPath:         keyFile.Name(),
					tokenPath:             tokenFile.Name(),
				},
				kubeConfigCredential: &kubeconfig.Credential{
					ClientCertificate: "client-certificate-from-kubeconfig\ncertificate-authority-from-kubeconfig",
					ClientKey:         "client-key-from-kubeconfig",
					Token:             "token-from-kubeconfig",
				},
			},
			want: &credentials.CredentialOption{
				ClientCertificateData: "client-certificate-from-file",
				ClientKeyData:         "client-key-from-file",
				Token:                 "token-from-file",
			},
		},
		{
			name: "args certificate/key overrides kubeconfig value",
			args: args{
				args: &rootCmdArgs{
					clientCertificatePath: certFile.Name(),
					clientKeyPath:         keyFile.Name(),
				},
				kubeConfigCredential: &kubeconfig.Credential{
					ClientCertificate: "client-certificate-from-kubeconfig\ncertificate-authority-from-kubeconfig",
					ClientKey:         "client-key-from-kubeconfig",
					Token:             "token-from-kubeconfig",
				},
			},
			want: &credentials.CredentialOption{
				ClientCertificateData: "client-certificate-from-file",
				ClientKeyData:         "client-key-from-file",
				Token:                 "token-from-kubeconfig",
			},
		},
		{
			name: "args token overrides kubeconfig value",
			args: args{
				args: &rootCmdArgs{
					tokenPath: tokenFile.Name(),
				},
				kubeConfigCredential: &kubeconfig.Credential{
					ClientCertificate: "client-certificate-from-kubeconfig\ncertificate-authority-from-kubeconfig",
					ClientKey:         "client-key-from-kubeconfig",
					Token:             "token-from-kubeconfig",
				},
			},
			want: &credentials.CredentialOption{
				ClientCertificateData: "client-certificate-from-kubeconfig\ncertificate-authority-from-kubeconfig",
				ClientKeyData:         "client-key-from-kubeconfig",
				Token:                 "token-from-file",
			},
		},
		{
			name: "kubeconfig.Credential is empty",
			args: args{
				args: &rootCmdArgs{
					clientCertificatePath: certFile.Name(),
					clientKeyPath:         keyFile.Name(),
					tokenPath:             tokenFile.Name(),
				},
				kubeConfigCredential: &kubeconfig.Credential{},
			},
			want: &credentials.CredentialOption{
				ClientCertificateData: "client-certificate-from-file",
				ClientKeyData:         "client-key-from-file",
				Token:                 "token-from-file",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeCredentialOptions(tt.args.args, tt.args.kubeConfigCredential)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeCredentialOptions() error = %+v, wantErr %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeCredentialOptions() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

package cmd

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func Test_arguments_validate(t *testing.T) {
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
			args := &arguments{
				clientCertificatePath: tt.fields.clientCertificatePath,
				clientKeyPath:         tt.fields.clientKeyPath,
				tokenPath:             tt.fields.tokenPath,
				beforeExecCommand:     tt.fields.beforeExecCommand,
			}
			if err := args.validate(); (err != nil) != tt.wantErr {
				t.Errorf("arguments.validate() error = %v, wantErr %v", err, tt.wantErr)
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
	if _, err = tokenFile.Write([]byte("token-string")); err != nil {
		t.Errorf("ioutil.Write() error = %v", err)
		return
	}
	certFile, err := ioutil.TempFile(testDir, "tls.crt")
	if _, err = certFile.Write([]byte("cert-string")); err != nil {
		t.Errorf("ioutil.Write() error = %v", err)
		return
	}
	keyFile, err := ioutil.TempFile(testDir, "tls.key")
	if _, err = keyFile.Write([]byte("key-string")); err != nil {
		t.Errorf("ioutil.Write() error = %v", err)
		return
	}

	type args struct {
		arguments      arguments
		kubeconfigData string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				arguments: arguments{
					tokenPath:             tokenFile.Name(),
					clientCertificatePath: certFile.Name(),
					clientKeyPath:         keyFile.Name(),
				},
				kubeconfigData: `apiVersion: v1
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
      apiVersion: client.authentication.k8s.io/v1beta1
`,
			},
			want: []byte(`{"kind":"ExecCredential","apiVersion":"client.authentication.k8s.io/v1beta1","spec":{},"status":{"token":"token-string","clientCertificateData":"cert-string","clientKeyData":"key-string"}}`),
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

			runner, err := newRunner(&tt.args.arguments)
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

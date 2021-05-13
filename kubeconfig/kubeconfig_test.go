package kubeconfig

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"text/template"

	"github.com/google/go-cmp/cmp"
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

func TestKubeconfig_UpdateCurrentUserExecConfig(t *testing.T) {
	type fields struct {
		kubeconfigString string
	}
	type args struct {
		apiVersion string
		cmd        string
		args       []string
		envs       map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				apiVersion: "client.authentication.k8s.io/v1beta1",
				cmd:        "/cmd",
				args:       []string{"-args1", "-args2"},
				envs:       map[string]string{"ENV1": "val1"},
			},
			fields: fields{
				kubeconfigString: `apiVersion: v1
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
current-context: context1
kind: Config
preferences: {}
users:
- name: user1
  user:
    token: hoge
`,
			},
			want: `apiVersion: v1
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
current-context: context1
kind: Config
preferences: {}
users:
- name: user1
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      args:
      - -args1
      - -args2
      command: /cmd
      env:
      - name: ENV1
        value: val1
      provideClusterInfo: false
`,
		},
	}

	testDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("TempDir() error = %v", err)
		return
	}
	defer os.Remove(testDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kubeConfigFile, err := ioutil.TempFile(testDir, "-kubeconfig")
			if _, err = kubeConfigFile.Write([]byte(tt.fields.kubeconfigString)); err != nil {
				t.Errorf("ioutil.Write() error = %v", err)
				return
			}
			os.Setenv("KUBECONFIG", kubeConfigFile.Name())

			// k := New()
			// if err := k.UpdateCurrentUserExecConfig(tt.args.apiVersion, tt.args.cmd, tt.args.args, tt.args.envs); (err != nil) != tt.wantErr {
			// 	t.Errorf("Kubeconfig.UpdateCurrentUserExecConfig() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }

			buf, err := ioutil.ReadFile(kubeConfigFile.Name())
			if err != nil {
				t.Errorf("ioutil.ReadFile() error = %v", err)
			}
			got := string(buf)
			if d := cmp.Diff(tt.want, got); d != "" {
				t.Errorf("Kubeconfig.UpdateCurrentUserExecConfig() mismatch:\n%s", d)
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

func TestKubeconfig_Diff(t *testing.T) {
	type fields struct {
		currentKubeconfigString string
	}
	type args struct {
		newRawConfig api.Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "same file",
			fields: fields{
				currentKubeconfigString: `apiVersion: v1
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
current-context: context1
kind: Config
preferences: {}
users:
- name: user1
  user: {}
`,
			},
			want: "",
			args: args{
				newRawConfig: api.Config{
					Clusters: map[string]*api.Cluster{
						"server1": {
							Server: "https://127.0.0.1",
						},
					},
					AuthInfos: map[string]*api.AuthInfo{
						"user1": {},
					},
					Contexts: map[string]*api.Context{
						"context1": {
							Cluster:   "server1",
							Namespace: "kube-system",
							AuthInfo:  "user1",
						},
					},
					CurrentContext: "context1",
					Extensions:     map[string]runtime.Object{},
				},
			},
		},
		{
			name: "add user exec",
			fields: fields{
				currentKubeconfigString: `apiVersion: v1
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
current-context: context1
kind: Config
preferences: {}
users:
- name: user1
  user: {}
`,
			},
			want: `  (
  	"""
  	... // 14 identical lines
  	users:
  	- name: user1
- 	  user: {}
+ 	  user:
+ 	    exec:
+ 	      apiVersion: client.authentication.k8s.io/v1beta1
+ 	      args: null
+ 	      command: cmd
+ 	      env: null
+ 	      provideClusterInfo: false
  	"""
  )
`,
			args: args{
				newRawConfig: api.Config{
					Clusters: map[string]*api.Cluster{
						"server1": {
							Server: "https://127.0.0.1",
						},
					},
					AuthInfos: map[string]*api.AuthInfo{
						"user1": {
							Exec: &api.ExecConfig{
								APIVersion: "client.authentication.k8s.io/v1beta1",
								Command:    "cmd",
							},
						},
					},
					Contexts: map[string]*api.Context{
						"context1": {
							Cluster:   "server1",
							Namespace: "kube-system",
							AuthInfo:  "user1",
						},
					},
					CurrentContext: "context1",
					Extensions:     map[string]runtime.Object{},
				},
			},
		},
	}
	testDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Errorf("TempDir() error = %v", err)
		return
	}
	defer os.Remove(testDir)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kubeConfigFile, err := ioutil.TempFile(testDir, "-kubeconfig")
			if _, err = kubeConfigFile.Write([]byte(tt.fields.currentKubeconfigString)); err != nil {
				t.Errorf("ioutil.Write() error = %v", err)
				return
			}
			os.Setenv("KUBECONFIG", kubeConfigFile.Name())

			k := New()

			got, err := k.Diff(tt.args.newRawConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("Kubeconfig.Diff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Kubeconfig.Diff() = \ngot:\n%v, want\n%v", got, tt.want)
			}
		})
	}
}

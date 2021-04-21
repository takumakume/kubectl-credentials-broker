package kubeconfig

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var kubeconfigString = `
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
    token: hoge
`

var kubeconfigStringNoCurrentConfig = `
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
    token: hoge
`

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
				kubeconfigString: kubeconfigString,
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
				kubeconfigString: kubeconfigStringNoCurrentConfig,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientConfig, err := clientcmd.NewClientConfigFromBytes([]byte(tt.fields.kubeconfigString))
			if err != nil {
				t.Errorf("Kubeconfig.ReadCurrentContext() test data error = %+v", err)
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
				kubeconfigString: kubeconfigString,
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
				kubeconfigString: kubeconfigString,
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
				kubeconfigString: kubeconfigString,
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
				kubeconfigString: kubeconfigString,
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

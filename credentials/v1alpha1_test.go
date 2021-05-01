package credentials

import (
	"reflect"
	"testing"
)

func TestV1Alpha1_ToJSON(t *testing.T) {
	type args struct {
		opts *CredentialOption
	}
	tests := []struct {
		name    string
		o       *V1Alpha1
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{opts: &CredentialOption{
				ClientCertificateData: "hoge",
				ClientKeyData:         "foo",
				Token:                 "bar",
			}},
			want: []byte(`{"kind":"ExecCredential","apiVersion":"client.authentication.k8s.io/v1alpha1","spec":{},"status":{"token":"bar","clientCertificateData":"hoge","clientKeyData":"foo"}}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &V1Alpha1{}
			got, err := o.ToJSON(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("V1Alpha1.ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("V1Alpha1.ToJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

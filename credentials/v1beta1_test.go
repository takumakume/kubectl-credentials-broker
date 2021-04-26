package credentials

import (
	"reflect"
	"testing"
)

func TestV1Beta1_ToJSON(t *testing.T) {
	type args struct {
		opts CredentialOptions
	}
	tests := []struct {
		name    string
		o       *V1Beta1
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{opts: CredentialOptions{
				ClientCertificateData: "hoge",
				ClientKeyData:         "foo",
				Token:                 "bar",
			}},
			want: []byte(`{"kind":"ExecCredential","apiVersion":"client.authentication.k8s.io/v1beta1","spec":{},"status":{"token":"bar","clientCertificateData":"hoge","clientKeyData":"foo"}}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &V1Beta1{}
			got, err := o.ToJSON(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("V1Beta1.ToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("V1Beta1.ToJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

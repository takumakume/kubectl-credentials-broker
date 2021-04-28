package cmd

import "testing"

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

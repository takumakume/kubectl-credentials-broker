package credentials

type CredentialOptions struct {
	ClientCertificateData string
	ClientKeyData         string
	Token                 string
}

type Credentials interface {
	ToJSON(opts CredentialOptions) ([]byte, error)
}

package credentials

type Credential interface {
	APIVersionString() string
	ToJSON(opts *CredentialOption) ([]byte, error)
}

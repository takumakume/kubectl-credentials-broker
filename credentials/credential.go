package credentials

type Credential interface {
	ToJSON(opts *CredentialOption) ([]byte, error)
}

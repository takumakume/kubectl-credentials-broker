package credentials

type Credentials interface {
	ToJSON(opts *CredentialOptions) ([]byte, error)
}

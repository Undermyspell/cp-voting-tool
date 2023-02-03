package jwks

func New() KeyfuncProvider {
	return &MockKeyfuncProvider{}
}

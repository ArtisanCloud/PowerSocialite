package contracts

type Params interface {
	Get(string) string
}

type ISession interface {
	GetAuthURL() (string, error)
	Marshal() string
	Authorize(IProvider, Params) (string, error)
}

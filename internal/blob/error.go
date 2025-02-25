package blob

var ErrException = &exception{}

func NewException(err error) *exception {
	return &exception{
		error: err,
	}
}

type exception struct {
	error
}

func (op *exception) mustEmbedded() {}

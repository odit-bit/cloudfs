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

// record in db is existed

type RecordIsExist struct {
	Err error
}

func (r *RecordIsExist) Error() string {
	if r.Err != nil {
		return r.Error()
	}
	return ""
}
func (r *RecordIsExist) mustEmbedded() {}

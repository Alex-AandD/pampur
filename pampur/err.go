package pampur

type Error interface {
	error
}	

type HttpError struct {
	Code 	int
	Msg string
}

func (err HttpError) Status() int {
	return err.Code
}

func (err HttpError) Error() string {
	return err.Msg
}

type RouterError struct {
	Msg string
}

func (err RouterError) Error() string {
	return err.Msg
}

func NewRouterError(m string) RouterError {
	return RouterError{Msg: m}
}

func NewHttpError(s int, m string) HttpError {
	return HttpError{Msg: m, Code: s}
}
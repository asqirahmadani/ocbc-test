package utils

type NotFoundError struct {
	Msg string
}

func (e NotFoundError) Error() string {
	return e.Msg
}

func ErrNotFound(msg string) error {
	return NotFoundError{
		Msg: msg,
	}
}

type BadRequestError struct {
	Msg string
}

func (e BadRequestError) Error() string {
	return e.Msg
}

func ErrBadRequest(msg string) error {
	return BadRequestError{
		Msg: msg,
	}
}

type InternalServerError struct {
	Msg string
}

func (e InternalServerError) Error() string {
	return e.Msg
}

func ErrInternalServer(msg string) error {
	return InternalServerError{
		Msg: msg,
	}
}

type AuthorizationError struct {
	Msg string
}

func (e AuthorizationError) Error() string {
	return e.Msg
}

func ErrUnauthorize(msg string) error {
	return AuthorizationError{
		Msg: msg,
	}
}

type UnprocessableEntityError struct {
	Msg string
}

func (e UnprocessableEntityError) Error() string {
	return e.Msg
}

func ErrUnprocessableEntity(msg string) error {
	return UnprocessableEntityError{
		Msg: msg,
	}
}

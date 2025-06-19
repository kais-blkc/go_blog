package shared

import "errors"

var (
	ErrUnauthorized      = errors.New("требуется авторизация")
	ErrInvalidUserIDType = errors.New("некорректный тип ID пользователя")
)

package errors

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("Адрес электронной почты уже существует")
	ErrInvalidCredentials = errors.New("Неверный адрес электронной почты или пароль")

	ErrUserNotFound     = errors.New("Пользователь не найден")
	ErrUserInactive     = errors.New("Пользователь неактивен")
	ErrAlreadyOrganizer = errors.New("Пользователь уже является организатором")
	ErrNotOrganizer     = errors.New("Пользователь не является организатором")
)

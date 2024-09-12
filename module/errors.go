package module

import "errors"

var (
	ErrInvalidModuleType  = errors.New("invalid plugin type")
	ErrModulePathNotFound = errors.New("Module path not found")
	ErrModuleNotFound     = errors.New("Module not found")
)

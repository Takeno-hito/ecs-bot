package bot

import "errors"

var (
	ErrUnderConstruction = errors.New("this command is under construction")
	ErrUnknownCommand    = errors.New("unknown command")
	ErrNeedRegister      = errors.New("you need to register first")
)

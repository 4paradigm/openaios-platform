package token

import (
	"github.com/pkg/errors"

	"time"
)

//TtyParameter kubectl tty param
type TtyParameter struct {
	Title string
	Arg   []string
}

var (
	InvalidToken = errors.New("ERROR:Invalid Token")
	NoTokenProvided = errors.New("ERROR:No Token Provided")
)

//interface that defines token cache behavior
type Cache interface {
	Get(token string) *TtyParameter
	Delete(token string) error
	Add(token string, param *TtyParameter, d time.Duration) error
}
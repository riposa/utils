package errors

import (
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	log.Print(New(ErrSuccess))
	log.Print(New(ErrUnknownError))
	log.Print(New(ErrUnstableNetwork))
	log.Print(New(ErrPermissionDeny))
	log.Print(New(ErrServiceUnderMaintaining))
	log.Print(New(ErrTooMuchRequest))
	log.Print(New(ErrServiceNotFound))

	log.Print(New(ErrNeedLogin))
	log.Print(New(ErrTokenExpired))
}

func TestNewFormat(t *testing.T) {
	log.Print(NewFormat(30, "test", 10))
}

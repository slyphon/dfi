package pathlib

import (
	"errors"
	"os"
	"syscall"
)

var (
	notADirError = errors.New("not a directory")
)

func IsNotDir(err error) bool {
	return underlyingErrorIs(err, notADirError)
}

// copied out of the os/error.go stdlib implemntation of 1.13.7
func underlyingErrorIs(err, target error) bool {
	// Note that this function is not errors.Is:
	// underlyingError only unwraps the specific error-wrapping types
	// that it historically did, not all errors.Wrapper implementations.
	err = underlyingError(err)
	if err == target {
		return true
	}
	// To preserve prior behavior, only examine syscall errors.
	e, ok := err.(syscall.Errno)
	return ok && e.Is(target)
}

// underlyingError returns the underlying error for known os error types.
func underlyingError(err error) error {
	switch err := err.(type) {
	case *os.PathError:
		return err.Err
	case *os.LinkError:
		return err.Err
	case *os.SyscallError:
		return err.Err
	}
	return err
}

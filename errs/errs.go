package errs

import (
	"errors"
	"github.com/cihub/seelog"
)

func ReturnIfErrNotNil(err error) {
	if err != nil {
		return
	}
}

func LogIfErrNotNil(err error, msg ...string) {
	if err != nil {
		msgs := []string{err.Error()}
		msgs = append(msgs, msg...)
		seelog.Error(msgs)
	}
}

func ErrToString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func StringToErr(msg string) error {
	if msg != "" {
		return errors.New(msg)
	}
	return nil
}

package lib

import (
	"github.com/amenzhinsky/go-memexec"
)

func Wrap(bin []byte) (*memexec.Exec, error) {
	exec, err := memexec.New(bin)

	if err != nil {
		return nil, err
	}

	return exec, nil

}

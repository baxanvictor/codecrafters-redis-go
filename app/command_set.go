package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal/storage"
)

func (s *server) handleCommandSetOnMaster(args []string) ([]byte, error) {
	err := s.handleCommandSet(args)
	if err != nil {
		return nil, err
	}

	return okSimpleString(), nil
}

func (s *server) handleCommandSetOnSlave(args []string) ([]byte, error) {
	return nil, s.handleCommandSet(args)
}

func (s *server) handleCommandSet(args []string) error {
	if len(args) < 2 {
		return errors.New("command set accepts two arguments")
	}

	if len(args)%2 != 0 {
		return errors.New("invalid arguments list, must come in pairs")
	}

	key := string(args[0])

	expVal := storage.NewExpiringValue(args[1])

	if len(args) > 2 {
		extraArg := args[2]
		if !strings.EqualFold(extraArg, "px") {
			return fmt.Errorf("unknown extra argument \"%s\"", extraArg)
		}

		exp, err := strconv.Atoi(args[3])
		if err != nil {
			return err
		}

		expVal.ExpiresIn = exp
	}

	s.dataMu.Lock()
	s.data[key] = expVal
	s.dataMu.Unlock()

	return nil
}

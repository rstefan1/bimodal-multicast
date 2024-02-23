/*
Copyright 2024 Robert Andrei STEFAN

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"strconv"
)

var (
	errEmptyAddress   = errors.New("empty address")
	errPortNotInteger = errors.New("port must be an integer number")
	errPortOutOfRange = errors.New("port must be between 1 and 65535")
)

// addrValidator returns an address validator.
func addrValidator() func(string) error {
	return func(addr string) error {
		if len(addr) == 0 {
			return errEmptyAddress
		}

		return nil
	}
}

// portValidator returns a port validator.
func portValidator() func(int) error {
	return func(port int) error {
		if port < 1 || port > 65535 {
			return errPortOutOfRange
		}

		return nil
	}
}

// portAsStringValidator returns a port (as string) validator.
func portAsStringValidator() func(string) error {
	return func(port string) error {
		iPort, err := strconv.ParseInt(port, 10, 32)
		if err != nil {
			return errPortNotInteger
		}

		return portValidator()(int(iPort))
	}
}

/*
Copyright 2020 Robert Andrei STEFAN

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

package validators

import "errors"

var errEmptyAddress = errors.New("empty address")

// AddrValidator returns an address validator.
func AddrValidator() func(string) error {
	return func(addr string) error {
		if len(addr) == 0 {
			return errEmptyAddress
		}

		return nil
	}
}

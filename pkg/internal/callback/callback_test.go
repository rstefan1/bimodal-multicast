/*
Copyright 2019 Robert Andrei STEFAN

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

package callback

import (
	"errors"
	"log/slog"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Callback interface", func() {
	Describe("NewRegistry func", func() {
		It("creates new registry when given callbacks map is empty", func() {
			cb := map[string]func(any, *slog.Logger) error{}
			r, err := NewRegistry(cb)
			Expect(err).To(Succeed())
			Expect(r.Callbacks).To(BeEquivalentTo(cb))
		})

		It("creates new registry when given callbacks map has more callbacks", func() {
			cb := map[string]func(any, *slog.Logger) error{
				"first-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
				"second-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
			}
			r, err := NewRegistry(cb)
			Expect(err).To(Succeed())
			Expect(r.Callbacks).To(BeEquivalentTo(cb))
		})

		It("returns error if given callbacks map is nil", func() {
			r, err := NewRegistry(nil)
			Expect(err).To(Not(Succeed()))
			Expect(r).To(BeNil())
		})
	})

	Describe("GetCallback func", func() {
		It("returns proper callback func when given callback type exists in registry", func() {
			cbType := "my-callback"
			cbFn := func(_ any, _ *slog.Logger) error {
				return nil
			}

			cb := map[string]func(any, *slog.Logger) error{
				"a-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
				cbType: cbFn,
				"another-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
			}
			r, err := NewRegistry(cb)
			Expect(err).To(Succeed())

			fn := r.GetCallback(cbType)
			Expect(reflect.ValueOf(fn)).To(Equal(reflect.ValueOf(cbFn)))
		})

		It("returns error when given callback type doesn't exist in registry", func() {
			cb := map[string]func(any, *slog.Logger) error{
				"a-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
				"another-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
			}
			r, err := NewRegistry(cb)
			Expect(err).To(Succeed())

			fn := r.GetCallback("not-existent-callback")
			Expect(fn).To(BeNil())
		})
	})

	Describe("ValidateCustomCallbacks func", func() {
		It("returns error when callbacks contain a `add-peer` type", func() {
			cb := map[string]func(any, *slog.Logger) error{
				"a-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
				"add-peer": func(_ any, _ *slog.Logger) error {
					return nil
				},
				"another-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
			}

			Expect(ValidateCustomCallbacks(cb)).To(MatchError(errors.New("callback type is not allowed"))) //nolint: goerr113
		})

		It("returns error when callbacks contain a `remove-peer` type", func() {
			cb := map[string]func(any, *slog.Logger) error{
				"a-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
				"remove-peer": func(_ any, _ *slog.Logger) error {
					return nil
				},
				"another-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
			}

			Expect(ValidateCustomCallbacks(cb)).To(MatchError(errors.New("callback type is not allowed"))) //nolint: goerr113
		})

		It("doesn't return error when all callback are valid", func() {
			cb := map[string]func(any, *slog.Logger) error{
				"a-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
				"valid-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
				"another-valid-callback": func(_ any, _ *slog.Logger) error {
					return nil
				},
			}

			Expect(ValidateCustomCallbacks(cb)).To(Succeed())
		})
	})
})

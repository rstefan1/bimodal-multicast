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
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CustomCallbackRegistry interface", func() {
	It("creates new registry when given callbacks map is empty", func() {
		cb := map[string]func(string) (bool, error){}
		r, err := NewCustomRegistry(cb)
		Expect(err).To(Succeed())
		Expect(r.callbacks).To(Equal(cb))
	})

	It("creates new registry when given callbacks map has more callbacks", func() {
		cb := map[string]func(string) (bool, error){
			"first-callback": func(msg string) (bool, error) {
				return true, nil
			},
			"second-callback": func(msg string) (bool, error) {
				return false, nil
			},
		}
		r, err := NewCustomRegistry(cb)
		Expect(err).To(Succeed())
		Expect(r.callbacks).To(Equal(cb))
	})

	It("returns error if given callbacks map is empty", func() {
		r, err := NewCustomRegistry(nil)
		Expect(err).To(Not(Succeed()))
		Expect(r).To(BeNil())
	})

	It("returns proper callback func when given callback type exists in registry", func() {
		var (
			cbType string
			cbFn   func(string) (bool, error)
		)

		cbType = "my-callback"
		cbFn = func(msg string) (bool, error) {
			return true, nil
		}

		cb := map[string]func(string) (bool, error){
			cbType: cbFn,
		}
		r, err := NewCustomRegistry(cb)
		Expect(err).To(Succeed())

		fn, err := r.GetCustomCallback(cbType)
		Expect(err).To(Succeed())
		Expect(reflect.ValueOf(fn)).To(Equal(reflect.ValueOf(cbFn)))
	})

	It("returns error when given callback type doesn't exist in registry", func() {
		cb := map[string]func(string) (bool, error){}
		r, err := NewCustomRegistry(cb)
		Expect(err).To(Succeed())

		fn, err := r.GetCustomCallback("mu-callback")
		Expect(err).To(Not(Succeed()))
		Expect(fn).To(BeNil())
	})
})

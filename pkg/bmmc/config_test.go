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

package bmmc

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

// newDummyConfig creates new dummy bmmc config.
func newDummyConfig() *Config {
	host, err := peer.NewHTTPPeer("localhost", "19123", &http.Client{})
	Expect(err).ToNot(HaveOccurred())

	return &Config{
		Host:   host,
		Beta:   0.45,
		Logger: log.New(os.Stdout, "", 0),
		Callbacks: map[string]func(interface{}, *log.Logger) error{
			"awesome-callback": func(_ interface{}, _ *log.Logger) error {
				return nil
			},
		},
		RoundDuration:     time.Millisecond * 100,
		BufferSize:        32,
		ReadHeaderTimeout: time.Second * 19,
	}
}

var _ = Describe("BMMC Config", func() {
	var cfg *Config

	BeforeEach(func() {
		cfg = newDummyConfig()
	})

	Describe("validate func", func() {
		It("doesn't return error when full config is given", func() {
			Expect(cfg.validate()).To(Succeed())
		})

		It("returns error when buffer size is invalid", func() {
			cfg.BufferSize = 0
			Expect(cfg.validate()).To(MatchError(errInvalidBufSize))
		})

		It("returns error when callback map contains an invalid callback (a default callback)", func() {
			cfg.Callbacks = map[string]func(interface{}, *log.Logger) error{
				"add-peer": func(_ interface{}, _ *log.Logger) error {
					return nil
				},
			}
			Expect(cfg.validate()).To(MatchError(errors.New("callback type is not allowed"))) //nolint: goerr113
		})
	})

	Describe("fillEmptyFields func", func() {
		It("set default values for all empty and nil optional fields", func() {
			cfg.Beta = 0
			cfg.RoundDuration = 0
			cfg.ReadHeaderTimeout = 0
			cfg.Logger = nil
			cfg.Callbacks = nil

			cfg.fillEmptyFields()

			Expect(cfg.Beta).To(Equal(0.3))
			Expect(cfg.RoundDuration).To(Equal(time.Millisecond * 100)) // defaultRoundDuration
			Expect(cfg.ReadHeaderTimeout).To(Equal(time.Second * 30))   // defaultReadHeaderTimeout
			Expect(cfg.Logger).NotTo(BeNil())
			Expect(cfg.Callbacks).NotTo(BeNil())
		})
	})
})

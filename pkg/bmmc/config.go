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
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/validators"
)

const (
	defaultBeta          = 0.3
	defaultRoundDuration = time.Millisecond * 100
)

const (
	invalidBufSizeErr = "invalid buffer size"
)

// Config is the config for the protocol
type Config struct {
	// Addr is HTTP address for node which runs http servers
	// Required
	Addr string
	// Port is HTTP port for node which runs http servers
	// Required
	Port string
	// Beta is the expected fanout for gossip rounds
	// Optional
	Beta float64
	// Logger
	// Optional
	Logger *log.Logger
	// Callbacks funtions
	// Optional
	Callbacks map[string]func(interface{}, *log.Logger) error
	// Gossip round duration
	// Optional
	RoundDuration time.Duration
	// Buffer size
	// Required
	BufferSize int
}

// validate validates given config
func (cfg *Config) validate() error {
	if err := validators.AddrValidator()(cfg.Addr); err != nil {
		return err
	}

	if err := validators.PortAsStringValidator()(cfg.Port); err != nil {
		return err
	}

	if cfg.BufferSize <= 0 {
		return fmt.Errorf(invalidBufSizeErr)
	}

	return nil
}

// fillEmptyFields set default values for optional empty fields.
func (cfg *Config) fillEmptyFields() {
	if cfg.Beta == 0 {
		cfg.Beta = defaultBeta
	}

	if cfg.Logger == nil {
		cfg.Logger = log.New(os.Stdout, "", 0)
	}

	if cfg.RoundDuration == 0 {
		cfg.RoundDuration = defaultRoundDuration
	}

	if cfg.Callbacks == nil {
		cfg.Callbacks = map[string]func(interface{}, *log.Logger) error{}
	}
}

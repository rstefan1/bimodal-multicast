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

package bmmc

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/internal/callback"
	"github.com/rstefan1/bimodal-multicast/pkg/internal/peer"
)

const (
	defaultBeta          = 0.3
	defaultRoundDuration = time.Millisecond * 100
)

var errInvalidBufSize = errors.New("invalid buffer size")

// Config is the config for the protocol.
type Config struct {
	// Host is the host peer.
	// Required.
	Host peer.Peer
	// Beta is the expected fanout for gossip rounds.
	// Optional
	Beta float64
	// Logger.
	// Optional
	Logger *slog.Logger
	// Callbacks functions.
	// Optional
	Callbacks map[string]func(any, *slog.Logger) error
	// Gossip round duration.
	// Optional
	RoundDuration time.Duration
	// Buffer size.
	// When the buffer is full, the oldest message will be removed.
	// Required
	BufferSize int
}

// validate validates given config.
func (cfg *Config) validate() error {
	if cfg.BufferSize <= 0 {
		return errInvalidBufSize
	}

	return callback.ValidateCustomCallbacks(cfg.Callbacks) //nolint: wrapcheck
}

// fillEmptyFields set default values for optional empty fields.
func (cfg *Config) fillEmptyFields() {
	if cfg.Beta == 0 {
		cfg.Beta = defaultBeta
	}

	if cfg.Logger == nil {
		cfg.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	if cfg.RoundDuration == 0 {
		cfg.RoundDuration = defaultRoundDuration
	}

	if cfg.Callbacks == nil {
		cfg.Callbacks = map[string]func(any, *slog.Logger) error{}
	}
}

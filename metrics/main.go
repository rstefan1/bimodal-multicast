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

package main

import (
	"fmt"
	"time"

	"github.com/rstefan1/bimodal-multicast/pkg/bmmc"
)

func main() {
	noPeers := 50
	noRetries := 20
	minLoss := 0.0
	maxLoss := 0.9
	minBeta := 0.0
	maxBeta := 0.9
	timeout := time.Second * 3

	for loss := minLoss; loss <= maxLoss; loss = loss + 0.1 {
		for beta := minBeta; beta <= maxBeta; beta = beta + 0.1 {
			fmt.Println("Running BMMC for loss =", loss, "and beta =", beta, "...")
			if err := bmmc.RunWithSpec(noRetries, noPeers, loss, beta, timeout); err != nil {
				fmt.Println(err)
			}
		}
	}
}

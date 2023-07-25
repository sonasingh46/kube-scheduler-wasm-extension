/*
   Copyright 2023 The Kubernetes Authors.

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

// Package plugin is ported from the native Go version of the same plugin with
// some changes:
//
//   - The description was rewritten for clarity.
//   - Logic was refactored to be cleaner and more testable.
//   - Doesn't return an error if state has the wrong type, as it is
//     impossible: this panics instead with the default message.
//   - TODO: uses PreFilter instead of PreScore
//   - TODO: logging
//   - TODO: config
//
// See https://github.com/kubernetes-sigs/kube-scheduler-simulator/blob/simulator/v0.1.0/simulator/docs/sample/nodenumber/plugin.go
//
// Note: This is intentionally separate from the main package, for testing.
package plugin

import (
	"sigs.k8s.io/kube-scheduler-wasm-extension/guest/api"
	"sigs.k8s.io/kube-scheduler-wasm-extension/guest/api/proto"
)

type GPU struct {
	reverse bool
}

// EventsToRegister implements api.EnqueueExtensions
func (g *GPU) EventsToRegister() []api.ClusterEvent {
	return []api.ClusterEvent{
		{Resource: api.Node, ActionType: api.Add},
	}
}

func (pl *GPU) Filter(state api.CycleState, pod proto.Pod, nodeInfo api.NodeInfo) *api.Status {
	nodeLabels := nodeInfo.Node().Metadata().Labels
	podLabels := pod.Metadata().Labels
	podIsGPUIntensive := false
	for key := range podLabels {
		if key == "gpu" {
			podIsGPUIntensive = true
		}
	}
	nodeIsGPUIntensive := false
	for key := range nodeLabels {
		if key == "gpu" {
			nodeIsGPUIntensive = true
		}
	}

	if podIsGPUIntensive == nodeIsGPUIntensive {
		return &api.Status{
			Code:   api.StatusCodeSuccess,
			Reason: "[WASM Plugin]: Scheduling GPU intensive pod on GPU node and non GPU on general node",
		}
	}

	return &api.Status{
		Code:   api.StatusCodeUnschedulableAndUnresolvable,
		Reason: "[WASM Plugin]: No matching node for GPU/General pods",
	}

}

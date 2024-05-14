/*
Copyright 2025 The Karmada Authors.

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

package queue

import (
	"k8s.io/client-go/util/workqueue"
)

// rateLimitingSchedulingQueue is a wrapper of TypedRateLimitingInterface.
// Note: This queue may be deprecated in a future release.
type rateLimitingSchedulingQueue struct {
	delegate workqueue.TypedRateLimitingInterface[*QueuedBindingInfo]
}

func (rbq *rateLimitingSchedulingQueue) Push(bindingInfo *QueuedBindingInfo) {
	rbq.delegate.Add(bindingInfo)
}

func (rbq *rateLimitingSchedulingQueue) PushUnschedulableIfNotPresent(bindingInfo *QueuedBindingInfo) {
	// Default behavior: add to the rate limiter.
	rbq.PushBackoffIfNotPresent(bindingInfo)
}

func (rbq *rateLimitingSchedulingQueue) PushBackoffIfNotPresent(bindingInfo *QueuedBindingInfo) {
	rbq.delegate.AddRateLimited(bindingInfo)
}

func (rbq *rateLimitingSchedulingQueue) Pop() (*QueuedBindingInfo, bool) {
	return rbq.delegate.Get()
}

func (rbq *rateLimitingSchedulingQueue) Done(bindingInfo *QueuedBindingInfo) {
	rbq.delegate.Done(bindingInfo)
}

func (rbq *rateLimitingSchedulingQueue) Len() int {
	return rbq.delegate.Len()
}

func (rbq *rateLimitingSchedulingQueue) Forget(bindingInfo *QueuedBindingInfo) {
	rbq.delegate.Forget(bindingInfo)
}

func (rbq *rateLimitingSchedulingQueue) Run() {
	// ignore.
}

func (rbq *rateLimitingSchedulingQueue) Close() {
	rbq.delegate.ShutDown()
}

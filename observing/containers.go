/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package observing

import coreV1 "k8s.io/api/core/v1"

func getRuntimeAndHeadContainer(pod coreV1.PodSpec) (runtime, head coreV1.Container, ok bool) {
	var hasHeadContainer, hasRuntimeContainer = false, false

	for _, container := range pod.Containers {
		if container.Name == "runtime" {
			runtime = container
			hasRuntimeContainer = true
		}
		if container.Name == "head" {
			head = container
			hasHeadContainer = true
		}
	}

	ok = hasRuntimeContainer && hasHeadContainer
	return
}

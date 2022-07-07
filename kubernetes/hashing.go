/*
 * Copyright (c) Dolittle. All rights reserved.
 * Licensed under the MIT license. See LICENSE file in the project root for full license information.
 */

package kubernetes

import (
	"crypto/sha512"
	"fmt"
	"hash"
	coreV1 "k8s.io/api/core/v1"
	"sort"
)

type ConfigHasher struct {
	hasher hash.Hash
}

func NewConfigHasher() ConfigHasher {
	return ConfigHasher{
		hasher: sha512.New(),
	}
}

func (h ConfigHasher) WriteConfigMap(configMap *coreV1.ConfigMap) {
	if len(configMap.Data) > 0 {
		keys := make([]string, 0, len(configMap.Data))
		for key := range configMap.Data {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			h.hasher.Write([]byte(key))
			h.hasher.Write([]byte(configMap.Data[key]))
		}
	}
	if len(configMap.BinaryData) > 0 {
		keys := make([]string, 0, len(configMap.BinaryData))
		for key := range configMap.BinaryData {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			h.hasher.Write([]byte(key))
			h.hasher.Write(configMap.BinaryData[key])
		}
	}
}

func (h ConfigHasher) WriteSecret(secret *coreV1.Secret) {
	if len(secret.Data) > 0 {
		keys := make([]string, 0, len(secret.Data))
		for key := range secret.Data {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			h.hasher.Write([]byte(key))
			h.hasher.Write([]byte(secret.Data[key]))
		}
	}
}

func (h ConfigHasher) GetComputedHash() string {
	hash := h.hasher.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

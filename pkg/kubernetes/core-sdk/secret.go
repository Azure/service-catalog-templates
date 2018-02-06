// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package coresdk

import (
	core "k8s.io/api/core/v1"
)

// GetSecretFromCache retrieves a Secret by name from the informer cache.
func (sdk *SDK) GetSecretFromCache(namespace, name string) (*core.Secret, error) {
	s, err := sdk.SecretCache().Secrets(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return s.DeepCopy(), nil
}

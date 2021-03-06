// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package builder

import (
	"reflect"
	"strings"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
)

const (
	// SecretSuffix is the suffix applied to a secret name to build the service catalog managed secret name.
	SecretSuffix = "-template"
)

func BuildBoundSecret(secret *core.Secret, tbnd *templates.TemplatedBinding) (*core.Secret, error) {
	shadowSecret := &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BoundSecretName(secret.Name),
			Namespace: secret.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(secret, core.SchemeGroupVersion.WithKind("Secret")),
			},
		},
		Type: secret.Type,
		Data: mapSecretKeys(tbnd.Spec.SecretKeys, secret.Data),
	}

	return shadowSecret, nil
}

func RefreshSecret(svcSecret *core.Secret, tbnd *templates.TemplatedBinding, secret *core.Secret) (*core.Secret, bool) {
	// TODO: Sync all fields

	if reflect.DeepEqual(svcSecret.Data, secret.Data) {
		return nil, false
	}

	secret.Data = mapSecretKeys(tbnd.Spec.SecretKeys, svcSecret.Data)

	return secret, true
}

func ShadowSecretName(name string) string {
	return name + SecretSuffix
}

func BoundSecretName(name string) string {
	return strings.TrimSuffix(name, SecretSuffix)
}

func mapSecretKeys(keys map[string]string, data map[string][]byte) map[string][]byte {
	mappedData := make(map[string][]byte, len(data))

	for k, v := range data {
		if mappedKey, ok := keys[k]; ok {
			k = mappedKey
		}

		mappedData[k] = v
	}

	return mappedData
}

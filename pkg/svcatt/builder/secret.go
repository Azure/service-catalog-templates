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

func BuildBoundSecret(secret core.Secret, binding templates.TemplatedBinding) (*core.Secret, error) {
	shadowSecret := &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      BoundSecretName(secret.Name),
			Namespace: secret.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&secret, core.SchemeGroupVersion.WithKind("Secret")),
			},
		},
		Type: secret.Type,
		Data: mapSecretKeys(binding.Spec.SecretKeys, secret.Data),
	}

	return shadowSecret, nil
}

func RefreshSecret(svcSecret core.Secret, binding templates.TemplatedBinding, secret core.Secret) (*core.Secret, bool) {
	// TODO: Sync all fields

	if reflect.DeepEqual(svcSecret.Data, secret.Data) {
		return nil, false
	}

	updatedSecret := secret.DeepCopy()
	updatedSecret.Data = mapSecretKeys(binding.Spec.SecretKeys, svcSecret.Data)

	return updatedSecret, true
}

func ShadowSecretName(name string) string {
	return name + SecretSuffix
}

func BoundSecretName(name string) string {
	return strings.TrimRight(name, SecretSuffix)
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

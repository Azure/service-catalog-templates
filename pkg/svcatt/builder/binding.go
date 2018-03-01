package builder

import (
	"github.com/peterbourgon/mergemap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func BuildServiceBinding(binding templates.TemplatedBinding) *svcat.ServiceBinding {
	return &svcat.ServiceBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      binding.Name,
			Namespace: binding.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&binding, templates.SchemeGroupVersion.WithKind(templates.BindingKind)),
			},
		},
		Spec: svcat.ServiceBindingSpec{
			ServiceInstanceRef: binding.Spec.InstanceRef,
			Parameters:         binding.Spec.Parameters,
			ParametersFrom:     binding.Spec.ParametersFrom,
			SecretName:         ShadowSecretName(binding.Spec.SecretName),
		},
	}
}

func RefreshServiceBinding(bnd *templates.TemplatedBinding, svcBnd *svcat.ServiceBinding) *svcat.ServiceBinding {
	svcBnd = svcBnd.DeepCopy()

	svcBnd.Spec.Parameters = bnd.Spec.Parameters
	svcBnd.Spec.ParametersFrom = bnd.Spec.ParametersFrom

	return svcBnd
}

func ApplyBindingTemplate(binding templates.TemplatedBinding, template templates.BindingTemplate) (*templates.TemplatedBinding, error) {
	finalBinding := binding.DeepCopy()

	// Default the secret name to the instance name, if empty
	if finalBinding.Spec.SecretName == "" {
		finalBinding.Spec.SecretName = finalBinding.Spec.InstanceRef.Name
	}

	var err error
	finalBinding.Spec.Parameters, err = mergeParameters(finalBinding.Spec.Parameters, template.Spec.Parameters)
	if err != nil {
		return nil, err
	}

	finalBinding.Spec.ParametersFrom = selectParametersFromSource(finalBinding.Spec.ParametersFrom, template.Spec.ParametersFrom)

	finalBinding.Spec.SecretKeys = mergeSecretKeys(finalBinding.Spec.SecretKeys, template.Spec.SecretKeys)

	return finalBinding, nil
}

func mergeSecretKeys(bndKeys map[string]string, tmplKeys map[string]string) map[string]string {
	// TODO: Add tests and remove these ifs
	if tmplKeys == nil {
		return bndKeys
	}

	if bndKeys == nil {
		return tmplKeys
	}

	bndMap := make(map[string]interface{}, len(bndKeys))
	for k, v := range bndKeys {
		bndMap[k] = v
	}

	tmplMap := make(map[string]interface{}, len(bndKeys))
	for k, v := range tmplKeys {
		tmplMap[k] = v
	}

	mergedMap := mergemap.Merge(bndMap, tmplMap)

	mergedKeys := make(map[string]string, len(mergedMap))
	for k, v := range mergedMap {
		mergedKeys[k] = v.(string)
	}

	return mergedKeys
}

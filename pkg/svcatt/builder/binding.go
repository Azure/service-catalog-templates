package builder

import (
	"github.com/peterbourgon/mergemap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func BuildServiceBinding(tbnd *templates.TemplatedBinding) *svcat.ServiceBinding {
	return &svcat.ServiceBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tbnd.Name,
			Namespace: tbnd.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(tbnd, templates.SchemeGroupVersion.WithKind(templates.BindingKind)),
			},
		},
		Spec: svcat.ServiceBindingSpec{
			ServiceInstanceRef: tbnd.Spec.InstanceRef,
			Parameters:         tbnd.Spec.Parameters,
			ParametersFrom:     tbnd.Spec.ParametersFrom,
			SecretName:         ShadowSecretName(tbnd.Spec.SecretName),
		},
	}
}

func RefreshServiceBinding(tbnd *templates.TemplatedBinding, svcBnd *svcat.ServiceBinding) *svcat.ServiceBinding {
	svcBnd.Spec.Parameters = tbnd.Spec.Parameters
	svcBnd.Spec.ParametersFrom = tbnd.Spec.ParametersFrom

	return svcBnd
}

func ApplyBindingTemplate(tbnd *templates.TemplatedBinding, template templates.BindingTemplateInterface) (*templates.TemplatedBinding, error) {
	// Default the secret name to the instance name, if empty
	if tbnd.Spec.SecretName == "" {
		tbnd.Spec.SecretName = tbnd.Spec.InstanceRef.Name
	}

	var err error
	tbnd.Spec.Parameters, err = MergeParameters(tbnd.Spec.Parameters, template.GetParameters())
	if err != nil {
		return nil, err
	}

	tbnd.Spec.ParametersFrom = MergeParametersFromSource(tbnd.Spec.ParametersFrom, template.GetParametersFrom())

	tbnd.Spec.SecretKeys = MergeSecretKeys(tbnd.Spec.SecretKeys, template.GetSecretKeys())

	return tbnd, nil
}

func MergeSecretKeys(bndKeys map[string]string, tmplKeys map[string]string) map[string]string {
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

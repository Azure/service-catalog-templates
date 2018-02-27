package svcatt

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/peterbourgon/mergemap"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

const (
	// SecretSuffix is the suffix applied to a secret name to build the service catalog managed secret name.
	SecretSuffix = "-template"
)

func BuildServiceInstance(instance templates.CatalogInstance, template templates.InstanceTemplate) (*svcat.ServiceInstance, error) {
	finalInstance, err := mergeTemplateWithInstance(instance, template)
	if err != nil {
		return nil, err
	}

	// Verify we resolved a plan
	if finalInstance.Spec.ClassExternalName == "" || finalInstance.Spec.PlanExternalName == "" {
		return nil, errors.New("could not resolve a class and plan")
	}

	return &svcat.ServiceInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      finalInstance.Name,
			Namespace: finalInstance.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(finalInstance, templates.SchemeGroupVersion.WithKind(templates.InstanceKind)),
			},
		},
		Spec: svcat.ServiceInstanceSpec{
			PlanReference: svcat.PlanReference{
				ClusterServiceClassExternalName: finalInstance.Spec.ClassExternalName,
				ClusterServicePlanExternalName:  finalInstance.Spec.PlanExternalName,
			},
			Parameters:     finalInstance.Spec.Parameters,
			ParametersFrom: finalInstance.Spec.ParametersFrom,
			ExternalID:     finalInstance.Spec.ExternalID,
			UpdateRequests: finalInstance.Spec.UpdateRequests,
		},
	}, nil
}

func BuildServiceBinding(binding templates.CatalogBinding) *svcat.ServiceBinding {
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
			SecretName:         toSVCSecretName(binding.Spec.SecretName),
		},
	}
}

func BuildShadowSecret(secret *core.Secret, binding templates.CatalogBinding) (*core.Secret, error) {
	shadowSecret := &core.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      toShadowSecretName(secret.Name),
			Namespace: secret.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(secret, core.SchemeGroupVersion.WithKind("Secret")),
			},
		},
		Type: secret.Type,
		Data: mapSecretKeys(binding.Spec.SecretKeys, secret.Data),
	}

	return shadowSecret, nil
}

func RefreshServiceInstance(inst *templates.CatalogInstance, svcInst *svcat.ServiceInstance) *svcat.ServiceInstance {
	svcInst = svcInst.DeepCopy()

	svcInst.Spec.Parameters = inst.Spec.Parameters
	svcInst.Spec.ParametersFrom = inst.Spec.ParametersFrom
	svcInst.Spec.UpdateRequests = inst.Spec.UpdateRequests

	// TODO: Figure out what can be synced, what's immutable

	// TODO: Figure out how to sync resolved values, like plan
	if inst.Spec.ClassExternalName != "" && inst.Spec.PlanExternalName != "" {
		svcInst.Spec.ClusterServiceClassExternalName = inst.Spec.ClassExternalName
		svcInst.Spec.ClusterServicePlanExternalName = inst.Spec.PlanExternalName
	}

	return svcInst
}

func RefreshServiceBinding(bnd *templates.CatalogBinding, svcBnd *svcat.ServiceBinding) *svcat.ServiceBinding {
	svcBnd = svcBnd.DeepCopy()

	svcBnd.Spec.Parameters = bnd.Spec.Parameters
	svcBnd.Spec.ParametersFrom = bnd.Spec.ParametersFrom

	return svcBnd
}

func RefreshSecret(svcSecret core.Secret, secret core.Secret) (*core.Secret, bool) {
	// TODO: Sync all fields

	if reflect.DeepEqual(svcSecret.Data, secret.Data) {
		return nil, false
	}

	updatedSecret := secret.DeepCopy()
	updatedSecret.Data = svcSecret.Data

	return updatedSecret, true
}

func toSVCSecretName(name string) string {
	return name + SecretSuffix
}

func toShadowSecretName(name string) string {
	return strings.TrimRight(name, SecretSuffix)
}

func mergeTemplateWithInstance(instance templates.CatalogInstance, template templates.InstanceTemplate) (*templates.CatalogInstance, error) {
	finalInstance := instance.DeepCopy()

	if finalInstance.Spec.ClassExternalName == "" {
		finalInstance.Spec.ClassExternalName = template.Spec.ClassExternalName
	}
	if finalInstance.Spec.PlanExternalName == "" {
		finalInstance.Spec.PlanExternalName = template.Spec.PlanExternalName
	}

	var err error
	finalInstance.Spec.Parameters, err = mergeParameters(finalInstance.Spec.Parameters, template.Spec.Parameters)
	if err != nil {
		return nil, err
	}

	finalInstance.Spec.ParametersFrom = selectParametersFromSource(finalInstance.Spec.ParametersFrom, template.Spec.ParametersFrom)

	return finalInstance, nil
}

func ApplyBindingTemplate(binding templates.CatalogBinding, template templates.BindingTemplate) (*templates.CatalogBinding, error) {
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

func mergeParameters(instParams *runtime.RawExtension, tmplParams *runtime.RawExtension) (*runtime.RawExtension, error) {
	if tmplParams == nil {
		return instParams, nil
	}

	if instParams == nil {
		return tmplParams, nil
	}

	var instMap, tmplMap map[string]interface{}
	json.Unmarshal(instParams.Raw, &instMap)
	json.Unmarshal(tmplParams.Raw, &tmplMap)

	merged := mergemap.Merge(instMap, tmplMap)

	result, err := json.Marshal(merged)
	if err != nil {
		return nil, fmt.Errorf("could not merge the instance and template parameters: %s", err)
	}

	return &runtime.RawExtension{Raw: result}, nil
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

func selectParametersFromSource(instParams []svcat.ParametersFromSource, tmplParams []svcat.ParametersFromSource) []svcat.ParametersFromSource {
	// TODO: I don't believe that merging is the right thing, so I'm only using the template if the instance didn't define anything
	if len(instParams) == 0 {
		return tmplParams
	}

	return instParams
}

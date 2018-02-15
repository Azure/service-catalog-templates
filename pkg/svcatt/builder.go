package svcatt

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/peterbourgon/mergemap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func BuildServiceInstance(instance templates.Instance, template templates.InstanceTemplate) (*svcat.ServiceInstance, error) {
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
				*metav1.NewControllerRef(finalInstance, schema.GroupVersionKind{
					Group:   templates.SchemeGroupVersion.Group,
					Version: templates.SchemeGroupVersion.Version,
					Kind:    templates.InstanceKind,
				}),
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

func mergeTemplateWithInstance(instance templates.Instance, template templates.InstanceTemplate) (*templates.Instance, error) {
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

func selectParametersFromSource(instParams []svcat.ParametersFromSource, tmplParams []svcat.ParametersFromSource) []svcat.ParametersFromSource {
	// TODO: I don't believe that merging is the right thing, so I'm only using the template if the instance didn't define anything
	if len(instParams) == 0 {
		return tmplParams
	}

	return instParams
}

func RefreshServiceInstance(inst *templates.Instance, svcInst *svcat.ServiceInstance) *svcat.ServiceInstance {
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

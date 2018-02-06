// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package builder

import (
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func BuildServiceInstance(instance *templates.TemplatedInstance) (*svcat.ServiceInstance, error) {
	// Verify we resolved a plan
	if !isPlanReferenceSpecified(instance.Spec.PlanReference) {
		return nil, errors.New("could not resolve a class and plan")
	}

	return &svcat.ServiceInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(instance, templates.SchemeGroupVersion.WithKind(templates.InstanceKind)),
			},
		},
		Spec: svcat.ServiceInstanceSpec{
			PlanReference:  instance.Spec.PlanReference,
			Parameters:     instance.Spec.Parameters,
			ParametersFrom: instance.Spec.ParametersFrom,
			ExternalID:     instance.Spec.ExternalID,
			UpdateRequests: instance.Spec.UpdateRequests,
		},
	}, nil
}

func RefreshServiceInstance(inst *templates.TemplatedInstance, svcInst *svcat.ServiceInstance) *svcat.ServiceInstance {
	// TODO: Figure out what can be synced, what's immutable

	svcInst.Spec.PlanReference = inst.Spec.PlanReference
	svcInst.Spec.Parameters = inst.Spec.Parameters
	svcInst.Spec.ParametersFrom = inst.Spec.ParametersFrom
	svcInst.Spec.UpdateRequests = inst.Spec.UpdateRequests

	return svcInst
}

func ApplyInstanceTemplate(instance *templates.TemplatedInstance, template templates.InstanceTemplateInterface) (*templates.TemplatedInstance, error) {
	if !isPlanReferenceSpecified(instance.Spec.PlanReference) {
		instance.Spec.PlanReference = template.GetPlanReference()
	}

	var err error
	instance.Spec.Parameters, err = MergeParameters(instance.Spec.Parameters, template.GetParameters())
	if err != nil {
		return nil, err
	}

	instance.Spec.ParametersFrom = MergeParametersFromSource(instance.Spec.ParametersFrom, template.GetParametersFrom())

	return instance, nil
}

func MergePlanReference(pr svcat.PlanReference, template svcat.PlanReference) svcat.PlanReference {
	if !isPlanReferenceSpecified(template) {
		return pr
	}

	if !isPlanReferenceSpecified(pr) {
		return template
	}

	if template.ClusterServiceClassExternalName != "" {
		pr.ClusterServiceClassExternalName = template.ClusterServiceClassExternalName
	}
	if template.ClusterServiceClassName != "" {
		pr.ClusterServiceClassName = template.ClusterServiceClassName
	}

	if template.ClusterServicePlanExternalName != "" {
		pr.ClusterServicePlanExternalName = template.ClusterServicePlanExternalName
	}
	if template.ClusterServicePlanName != "" {
		pr.ClusterServicePlanName = template.ClusterServicePlanName
	}

	return pr
}

func isPlanReferenceSpecified(pr svcat.PlanReference) bool {
	return isClassSpecified(pr) && isPlanSpecified(pr)
}

func isClassSpecified(pr svcat.PlanReference) bool {
	return pr.ClusterServiceClassExternalName != "" ||
		pr.ClusterServiceClassName != ""
}

func isPlanSpecified(pr svcat.PlanReference) bool {
	return pr.ClusterServicePlanExternalName != "" ||
		pr.ClusterServicePlanName != ""
}

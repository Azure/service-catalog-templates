package builder

import (
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func BuildServiceInstance(instance templates.TemplatedInstance) (*svcat.ServiceInstance, error) {
	// Verify we resolved a plan
	if instance.Spec.ClassExternalName == "" || instance.Spec.PlanExternalName == "" {
		return nil, errors.New("could not resolve a class and plan")
	}

	return &svcat.ServiceInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(&instance, templates.SchemeGroupVersion.WithKind(templates.InstanceKind)),
			},
		},
		Spec: svcat.ServiceInstanceSpec{
			PlanReference: svcat.PlanReference{
				ClusterServiceClassExternalName: instance.Spec.ClassExternalName,
				ClusterServicePlanExternalName:  instance.Spec.PlanExternalName,
			},
			Parameters:     instance.Spec.Parameters,
			ParametersFrom: instance.Spec.ParametersFrom,
			ExternalID:     instance.Spec.ExternalID,
			UpdateRequests: instance.Spec.UpdateRequests,
		},
	}, nil
}

func RefreshServiceInstance(inst *templates.TemplatedInstance, svcInst *svcat.ServiceInstance) *svcat.ServiceInstance {
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

func ApplyInstanceTemplate(instance templates.TemplatedInstance, template templates.InstanceTemplate) (*templates.TemplatedInstance, error) {
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

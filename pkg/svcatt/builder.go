package svcatt

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func BuildServiceInstance(instance *templates.Instance, template *templates.InstanceTemplate) *svcat.ServiceInstance {
	return &svcat.ServiceInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(instance, schema.GroupVersionKind{
					Group:   templates.SchemeGroupVersion.Group,
					Version: templates.SchemeGroupVersion.Version,
					Kind:    templates.InstanceKind,
				}),
			},
		},
		Spec: svcat.ServiceInstanceSpec{
			PlanReference: svcat.PlanReference{
				ClusterServiceClassExternalName: instance.Spec.ClassExternalName,
				ClusterServicePlanExternalName:  instance.Spec.PlanExternalName,
			},
			Parameters:     instance.Spec.Parameters, // TODO: Figure out if these need deep copies
			ParametersFrom: instance.Spec.ParametersFrom,
			ExternalID:     instance.Spec.ExternalID,
			UpdateRequests: instance.Spec.UpdateRequests,
		},
	}
}

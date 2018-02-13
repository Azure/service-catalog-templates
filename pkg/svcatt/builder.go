package svcatt

import (
	templatesexperimental "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"

	svcatv1beta1 "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func BuildServiceInstance(instance *templatesexperimental.Instance, template *templatesexperimental.InstanceTemplate) *svcatv1beta1.ServiceInstance {
	return &svcatv1beta1.ServiceInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(instance, schema.GroupVersionKind{
					Group:   templatesexperimental.SchemeGroupVersion.Group,
					Version: templatesexperimental.SchemeGroupVersion.Version,
					Kind:    "Instance",
				}),
			},
		},
		Spec: svcatv1beta1.ServiceInstanceSpec{
			PlanReference: svcatv1beta1.PlanReference{
				ClusterServiceClassExternalName: instance.Spec.ClassExternalName,
				ClusterServicePlanExternalName:  instance.Spec.PlanExternalName,
			},
			// TODO: Copy parameters and remaining values
		},
	}
}

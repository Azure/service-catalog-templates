package sdk

import (
	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetTemplatedInstanceFromCache retrieves a TemplatedInstance by name from the informer cache.
func (sdk *SDK) GetInstanceFromCache(namespace, name string) (*templates.TemplatedInstance, error) {
	inst, err := sdk.InstanceCache().TemplatedInstances(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return inst.DeepCopy(), nil
}

func (sdk *SDK) GetInstanceTemplateByServiceType(serviceType, namespace string) (*templates.InstanceTemplate, error) {
	opts := metav1.ListOptions{
		LabelSelector: sdk.filterByServiceTypeLabel(serviceType).String(),
	}
	results, err := sdk.Templates().InstanceTemplates(namespace).List(opts)
	if err != nil {
		return nil, err
	}
	if len(results.Items) == 0 {
		return nil, nil
	}

	return &results.Items[0], nil
}

func (sdk *SDK) GetClusterInstanceTemplateByServiceType(serviceType string) (*templates.ClusterInstanceTemplate, error) {
	opts := metav1.ListOptions{
		LabelSelector: sdk.filterByServiceTypeLabel(serviceType).String(),
	}
	results, err := sdk.Templates().ClusterInstanceTemplates().List(opts)
	if err != nil {
		return nil, err
	}
	if len(results.Items) == 0 {
		return nil, nil
	}

	return &results.Items[0], nil
}

func (sdk *SDK) GetBrokerInstanceTemplateByServiceType(serviceType string) (*templates.BrokerInstanceTemplate, error) {
	opts := metav1.ListOptions{
		LabelSelector: sdk.filterByServiceTypeLabel(serviceType).String(),
	}
	results, err := sdk.Templates().BrokerInstanceTemplates().List(opts)
	if err != nil {
		return nil, err
	}
	if len(results.Items) == 0 {
		return nil, nil
	}

	return &results.Items[0], nil
}

package sdk

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

// GetTemplatedBindingFromCache retrieves a TemplatedBinding by name from the informer cache.
func (sdk *SDK) GetBindingFromCache(namespace, name string) (*templates.TemplatedBinding, error) {
	bnd, err := sdk.BindingCache().TemplatedBindings(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return bnd.DeepCopy(), nil
}

func (sdk *SDK) GetBindingOwner(svcBnd *svcat.ServiceBinding) (*templates.TemplatedBinding, error) {
	ownerBnd := metav1.GetControllerOf(svcBnd)
	if ownerBnd == nil {
		// Ignore unmanaged resources
		return nil, nil
	}
	tbnd, err := sdk.GetBindingFromCache(svcBnd.Namespace, ownerBnd.Name)
	return tbnd, err
}

func (sdk *SDK) GetBindingTemplateByServiceType(serviceType, namespace string) (*templates.BindingTemplate, error) {
	opts := metav1.ListOptions{
		LabelSelector: sdk.filterByServiceTypeLabel(serviceType).String(),
	}
	results, err := sdk.Templates().BindingTemplates(namespace).List(opts)
	if err != nil {
		return nil, err
	}
	if len(results.Items) == 0 {
		return nil, nil
	}

	return &results.Items[0], nil
}

func (sdk *SDK) GetClusterBindingTemplateByServiceType(serviceType string) (*templates.ClusterBindingTemplate, error) {
	opts := metav1.ListOptions{
		LabelSelector: sdk.filterByServiceTypeLabel(serviceType).String(),
	}
	results, err := sdk.Templates().ClusterBindingTemplates().List(opts)
	if err != nil {
		return nil, err
	}
	if len(results.Items) == 0 {
		return nil, nil
	}

	return &results.Items[0], nil
}

func (sdk *SDK) GetBrokerBindingTemplateByServiceType(serviceType string) (*templates.BrokerBindingTemplate, error) {
	opts := metav1.ListOptions{
		LabelSelector: sdk.filterByServiceTypeLabel(serviceType).String(),
	}
	results, err := sdk.Templates().BrokerBindingTemplates().List(opts)
	if err != nil {
		return nil, err
	}
	if len(results.Items) == 0 {
		return nil, nil
	}

	return &results.Items[0], nil
}

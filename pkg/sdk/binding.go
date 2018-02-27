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

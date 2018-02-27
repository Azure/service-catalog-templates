package svcatsdk

import (
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetBindingFromCache retrieves a ServiceBinding by name from the informer cache.
func (sdk *SDK) GetBindingFromCache(namespace, name string) (*svcat.ServiceBinding, error) {
	bnd, err := sdk.BindingCache().ServiceBindings(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return bnd.DeepCopy(), nil
}

func (sdk *SDK) GetSecretOwner(svcSecret *core.Secret) (*svcat.ServiceBinding, error) {
	ownerSvcBnd := metav1.GetControllerOf(svcSecret)
	if ownerSvcBnd == nil {
		// Ignore unmanaged secrets
		return nil, nil
	}
	svcBnd, err := sdk.GetBindingFromCache(svcSecret.Namespace, ownerSvcBnd.Name)
	if err != nil {
		return nil, err
	}
	return svcBnd, err
}

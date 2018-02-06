// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package servicecatalogsdk

import (
	servicecatalog "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetBindingFromCache retrieves a ServiceBinding by name from the informer cache.
func (sdk *SDK) GetBindingFromCache(namespace, name string) (*servicecatalog.ServiceBinding, error) {
	bnd, err := sdk.BindingCache().ServiceBindings(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return bnd.DeepCopy(), nil
}

func (sdk *SDK) GetServiceBinding(namespace, name string) (*servicecatalog.ServiceBinding, error) {
	if sdk.HasCache() {
		return sdk.GetBindingFromCache(namespace, name)
	}
	return sdk.ServiceCatalog().ServiceBindings(namespace).Get(name, meta.GetOptions{})
}

func (sdk *SDK) GetSecretOwner(svcSecret *core.Secret) (*servicecatalog.ServiceBinding, error) {
	ownerSvcBnd := meta.GetControllerOf(svcSecret)
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

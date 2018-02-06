// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package servicecatalogsdk

import (
	servicecatalog "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetInstanceFromCache retrieves a ServiceInstance by name from the informer cache.
func (sdk *SDK) GetInstanceFromCache(namespace, name string) (*servicecatalog.ServiceInstance, error) {
	inst, err := sdk.InstanceCache().ServiceInstances(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return inst.DeepCopy(), nil
}

func (sdk *SDK) GetServiceInstance(namespace, name string) (*servicecatalog.ServiceInstance, error) {
	if sdk.HasCache() {
		return sdk.GetInstanceFromCache(namespace, name)
	}
	return sdk.ServiceCatalog().ServiceInstances(namespace).Get(name, meta.GetOptions{})
}

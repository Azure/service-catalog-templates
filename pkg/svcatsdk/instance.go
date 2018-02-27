package svcatsdk

import (
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

// GetInstanceFromCache retrieves a ServiceInstance by name from the informer cache.
func (sdk *SDK) GetInstanceFromCache(namespace, name string) (*svcat.ServiceInstance, error) {
	inst, err := sdk.InstanceCache().ServiceInstances(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return inst.DeepCopy(), nil
}

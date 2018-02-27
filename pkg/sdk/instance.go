package sdk

import (
	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
)

// GetTemplatedInstanceFromCache retrieves a TemplatedInstance by name from the informer cache.
func (sdk *SDK) GetInstanceFromCache(namespace, name string) (*templates.TemplatedInstance, error) {
	inst, err := sdk.InstanceCache().TemplatedInstances(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return inst.DeepCopy(), nil
}

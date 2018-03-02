package sdk

import (
	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"k8s.io/apimachinery/pkg/labels"
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
	opts := labels.SelectorFromSet(map[string]string{
		templates.FieldServiceTypeName: serviceType,
	})
	results, err := sdk.InstanceTemplateCache().InstanceTemplates(namespace).List(opts)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}

	return results[0].DeepCopy(), nil
}

func (sdk *SDK) GetClusterInstanceTemplateByServiceType(serviceType string) (*templates.ClusterInstanceTemplate, error) {
	opts := labels.SelectorFromSet(map[string]string{
		templates.FieldServiceTypeName: serviceType,
	})
	results, err := sdk.ClusterInstanceTemplateCache().List(opts)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}

	return results[0].DeepCopy(), nil
}

func (sdk *SDK) GetBrokerInstanceTemplateByServiceType(serviceType string) (*templates.BrokerInstanceTemplate, error) {
	opts := labels.SelectorFromSet(map[string]string{
		templates.FieldServiceTypeName: serviceType,
	})
	results, err := sdk.BrokerInstanceTemplateCache().List(opts)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}

	return results[0].DeepCopy(), nil
}

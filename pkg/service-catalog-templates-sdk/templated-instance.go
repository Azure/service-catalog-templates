// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package servicecatalogtempltesdk

import (
	"fmt"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"github.com/Azure/service-catalog-templates/pkg/service-catalog-templates-sdk/errors"
	servicecatalog "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/svcat/service-catalog"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (sdk *SDK) GetManagedServiceInstance(tinst *templates.TemplatedInstance) (inst *servicecatalog.ServiceInstance, err error) {
	inst, err = sdk.svcatSDK.GetServiceInstance(tinst.Namespace, tinst.Name)
	if err != nil {
		return nil, err
	}

	if !meta.IsControlledBy(inst, tinst) {
		return nil, errors.NewUnmanagedResource()
	}

	return inst, nil
}

// GetTemplatedInstanceFromCache retrieves a TemplatedInstance by name from the informer cache.
func (sdk *SDK) GetInstanceFromCache(namespace, name string) (*templates.TemplatedInstance, error) {
	inst, err := sdk.InstanceCache().TemplatedInstances(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return inst.DeepCopy(), nil
}

// RetrieveTemplatedInstances lists all instances in a namespace.
func (sdk *SDK) RetrieveTemplatedInstances(ns string) (*templates.TemplatedInstanceList, error) {
	instances, err := sdk.Templates().TemplatedInstances(ns).List(meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list templated instances in %s (%s)", ns, err)
	}

	return instances, nil
}

// RetrieveTemplatedInstance gets an instance by its name.
func (sdk *SDK) RetrieveTemplatedInstance(ns, name string) (*templates.TemplatedInstance, error) {
	instance, err := sdk.Templates().TemplatedInstances(ns).Get(name, meta.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get instance '%s.%s' (%s)", ns, name, err)
	}
	return instance, nil
}

// RetrieveTemplatedInstanceByBinding retrieves the parent instance for a binding.
func (sdk *SDK) RetrieveTemplatedInstanceByBinding(b *templates.TemplatedBinding,
) (*templates.TemplatedInstance, error) {
	ns := b.Namespace
	instName := b.Spec.TemplatedInstanceRef.Name
	inst, err := sdk.Templates().TemplatedInstances(ns).Get(instName, meta.GetOptions{})
	if err != nil {
		return nil, err
	}
	return inst, nil
}

// TemplatedInstanceParentHierarchy retrieves all ancestor resources of a templated instance.
func (sdk *SDK) TemplatedInstanceParentHierarchy(tinst *templates.TemplatedInstance,
) (*servicecatalog.ClusterServiceClass, *servicecatalog.ClusterServicePlan, *servicecatalog.ClusterServiceBroker, error) {
	inst, err := sdk.GetManagedServiceInstance(tinst)
	if err != nil {
		return nil, nil, nil, err
	}

	return sdk.svcatSDK.InstanceParentHierarchy(inst)
}

// Provision creates an instance of a service class and plan.
func (sdk *SDK) Provision(namespace, instanceName, serviceType, className, planName string,
	params interface{}, secrets map[string]string) (*templates.TemplatedInstance, error) {

	request := &templates.TemplatedInstance{
		ObjectMeta: meta.ObjectMeta{
			Name:      instanceName,
			Namespace: namespace,
		},
		Spec: templates.TemplatedInstanceSpec{
			ServiceType: serviceType,
			PlanReference: servicecatalog.PlanReference{
				ClusterServiceClassExternalName: className,
				ClusterServicePlanExternalName:  planName,
			},
			Parameters:     svcat.BuildParameters(params),
			ParametersFrom: svcat.BuildParametersFrom(secrets),
		},
	}

	result, err := sdk.Templates().TemplatedInstances(namespace).Create(request)
	if err != nil {
		return nil, fmt.Errorf("provision request failed (%s)", err)
	}
	return result, nil
}

// Deprovision deletes an instance.
func (sdk *SDK) Deprovision(namespace, instanceName string) error {
	err := sdk.Templates().TemplatedInstances(namespace).Delete(instanceName, &meta.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("deprovision request failed (%s)", err)
	}
	return nil
}

// TouchTemplatedInstance increments the updateRequests field on an instance to make
// service process it again (might be an update, delete, or noop)
func (sdk *SDK) TouchTemplatedInstance(ns, name string, retries int) error {
	for j := 0; j < retries; j++ {
		inst, err := sdk.RetrieveTemplatedInstance(ns, name)
		if err != nil {
			return err
		}

		inst.Spec.UpdateRequests = inst.Spec.UpdateRequests + 1

		_, err = sdk.Templates().TemplatedInstances(ns).Update(inst)
		if err == nil {
			return nil
		}
		// if we didn't get a conflict, no idea what happened
		if !apierrors.IsConflict(err) {
			return fmt.Errorf("could not touch templated instance (%s)", err)
		}
	}

	// conflict after `retries` tries
	return fmt.Errorf("could not update templated instance after %d tries", retries)
}

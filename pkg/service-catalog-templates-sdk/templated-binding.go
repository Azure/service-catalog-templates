// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package servicecatalogtempltesdk

import (
	"fmt"
	"strings"
	"sync"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	"github.com/Azure/service-catalog-templates/pkg/service-catalog-templates-sdk/errors"
	"github.com/hashicorp/go-multierror"
	servicecatalog "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/svcat/service-catalog"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (sdk *SDK) GetManagedServiceBinding(tbnd *templates.TemplatedBinding) (bnd *servicecatalog.ServiceBinding, err error) {
	bnd, err = sdk.svcatSDK.GetServiceBinding(tbnd.Namespace, tbnd.Name)
	if err != nil {
		return nil, err
	}

	if !meta.IsControlledBy(bnd, tbnd) {
		return nil, errors.NewUnmanagedResource()
	}

	return bnd, nil
}

// GetTemplatedBindingFromCache retrieves a TemplatedBinding by name from the informer cache.
func (sdk *SDK) GetBindingFromCache(namespace, name string) (*templates.TemplatedBinding, error) {
	bnd, err := sdk.BindingCache().TemplatedBindings(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return bnd.DeepCopy(), nil
}

func (sdk *SDK) GetBindingOwner(svcBnd *servicecatalog.ServiceBinding) (*templates.TemplatedBinding, error) {
	ownerBnd := meta.GetControllerOf(svcBnd)
	if ownerBnd == nil {
		// Ignore unmanaged resources
		return nil, nil
	}
	tbnd, err := sdk.GetBindingFromCache(svcBnd.Namespace, ownerBnd.Name)
	return tbnd, err
}

// RetrieveBindings lists all bindings in a namespace.
func (sdk *SDK) RetrieveTemplatedBindings(ns string) (*templates.TemplatedBindingList, error) {
	bindings, err := sdk.Templates().TemplatedBindings(ns).List(meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list bindings in %s (%s)", ns, err)
	}

	return bindings, nil
}

// RetrieveBinding gets a binding by its name.
func (sdk *SDK) RetrieveTemplatedBinding(ns, name string) (*templates.TemplatedBinding, error) {
	binding, err := sdk.Templates().TemplatedBindings(ns).Get(name, meta.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get binding '%s.%s' (%+v)", ns, name, err)
	}
	return binding, nil
}

// RetrieveBindingsByInstance gets all child bindings for an instance.
func (sdk *SDK) RetrieveTemplatedBindingsByInstance(instance *templates.TemplatedInstance,
) ([]templates.TemplatedBinding, error) {
	// Not using a filtered list operation because it's not supported yet.
	results, err := sdk.Templates().TemplatedBindings(instance.Namespace).List(meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to search bindings (%s)", err)
	}

	var bindings []templates.TemplatedBinding
	for _, binding := range results.Items {
		if binding.Spec.TemplatedInstanceRef.Name == instance.Name {
			bindings = append(bindings, binding)
		}
	}

	return bindings, nil
}

// Bind an instance to a secret.
func (sdk *SDK) Bind(namespace, bindingName, instanceName, secretName string,
	params map[string]string, secrets map[string]string) (*templates.TemplatedBinding, error) {

	// Manually defaulting the name of the binding
	// I'm not doing the same for the secret since the API handles defaulting that value.
	if bindingName == "" {
		bindingName = instanceName
	}

	request := &templates.TemplatedBinding{
		ObjectMeta: meta.ObjectMeta{
			Name:      bindingName,
			Namespace: namespace,
		},
		Spec: templates.TemplatedBindingSpec{
			TemplatedInstanceRef: servicecatalog.LocalObjectReference{
				Name: instanceName,
			},
			SecretName:     secretName,
			Parameters:     svcat.BuildParameters(params),
			ParametersFrom: svcat.BuildParametersFrom(secrets),
		},
	}

	result, err := sdk.Templates().TemplatedBindings(namespace).Create(request)
	if err != nil {
		return nil, fmt.Errorf("bind request failed (%s)", err)
	}

	return result, nil
}

// Unbind deletes all bindings associated to an instance.
func (sdk *SDK) Unbind(ns, instanceName string) ([]templates.TemplatedBinding, error) {
	instance, err := sdk.RetrieveTemplatedInstance(ns, instanceName)
	if err != nil {
		return nil, err
	}
	bindings, err := sdk.RetrieveTemplatedBindingsByInstance(instance)
	if err != nil {
		return nil, err
	}
	var g sync.WaitGroup
	errs := make(chan error, len(bindings))
	deletedBindings := make(chan templates.TemplatedBinding, len(bindings))
	for _, binding := range bindings {
		g.Add(1)
		go func(binding templates.TemplatedBinding) {
			defer g.Done()
			err := sdk.DeleteTemplatedBinding(binding.Namespace, binding.Name)
			if err == nil {
				deletedBindings <- binding
			}
			errs <- err
		}(binding)
	}

	g.Wait()
	close(errs)
	close(deletedBindings)

	// Collect any errors that occurred into a single formatted error
	bindErr := &multierror.Error{
		ErrorFormat: func(errors []error) string {
			return joinErrors("could not remove some bindings:", errors, "\n  ")
		},
	}
	for err := range errs {
		bindErr = multierror.Append(bindErr, err)
	}

	//Range over the deleted bindings to build a slice to return
	deleted := []templates.TemplatedBinding{}
	for b := range deletedBindings {
		deleted = append(deleted, b)
	}
	return deleted, bindErr.ErrorOrNil()
}

// DeleteBinding by name.
func (sdk *SDK) DeleteTemplatedBinding(ns, bindingName string) error {
	err := sdk.Templates().TemplatedBindings(ns).Delete(bindingName, &meta.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("remove binding %s/%s failed (%s)", ns, bindingName, err)
	}
	return nil
}

func joinErrors(groupMsg string, errors []error, sep string, a ...interface{}) string {
	if len(errors) == 0 {
		return ""
	}

	msgs := make([]string, 0, len(errors)+1)
	msgs = append(msgs, fmt.Sprintf(groupMsg, a...))
	for _, err := range errors {
		msgs = append(msgs, err.Error())
	}

	return strings.Join(msgs, sep)
}

// BindingParentHierarchy retrieves all ancestor resources of a binding.
func (sdk *SDK) TemplatedBindingParentHierarchy(tbnd *templates.TemplatedBinding,
) (*templates.TemplatedInstance, *servicecatalog.ClusterServiceClass, *servicecatalog.ClusterServicePlan, *servicecatalog.ClusterServiceBroker, error) {
	tinst, err := sdk.RetrieveTemplatedInstance(tbnd.Namespace, tbnd.Spec.TemplatedInstanceRef.Name)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	bnd, err := sdk.GetManagedServiceBinding(tbnd)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	_, class, plan, broker, err := sdk.svcatSDK.BindingParentHierarchy(bnd)
	return tinst, class, plan, broker, err
}

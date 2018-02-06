// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package servicecatalogtempltesdk

import (
	"fmt"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (sdk *SDK) GetBindingTemplates(ns string, cluster, broker bool, serviceType string) ([]templates.BindingTemplateInterface, error) {
	var t []templates.BindingTemplateInterface

	if broker {
		bbndts, err := sdk.GetBrokerBindingTemplatesByServiceType(serviceType)
		if err != nil {
			return nil, err
		}
		for _, bbndt := range bbndts.Items {
			t = append(t, &bbndt)
		}
		return t, err
	}

	if cluster {
		cbndt, err := sdk.GetClusterBindingTemplateByServiceType(serviceType)
		if err != nil {
			return nil, err
		}
		if cbndt != nil {
			t = append(t, cbndt)
		}
		return t, err
	}

	bndt, err := sdk.GetBindingTemplateByServiceType(serviceType, ns)
	if err != nil {
		return nil, err
	}
	if bndt != nil {
		t = append(t, bndt)
	}
	return t, err
}

func (sdk *SDK) GetBindingTemplate(ns string, cluster, broker bool, name string) (t templates.BindingTemplateInterface, err error) {
	if broker {
		return sdk.RetrieveBrokerBindingTemplate(name)
	}

	if cluster {
		return sdk.RetrieveClusterBindingTemplate(name)
	}

	return sdk.RetrieveBindingTemplate(ns, name)
}

func (sdk *SDK) GetBindingTemplateByServiceType(serviceType, namespace string) (*templates.BindingTemplate, error) {
	results, err := sdk.Templates().BindingTemplates(namespace).List(filterByServiceTypeLabel(serviceType))
	if err != nil {
		return nil, err
	}
	if len(results.Items) == 0 {
		return nil, nil
	}

	return &results.Items[0], nil
}

func (sdk *SDK) GetClusterBindingTemplateByServiceType(serviceType string) (*templates.ClusterBindingTemplate, error) {
	results, err := sdk.Templates().ClusterBindingTemplates().List(filterByServiceTypeLabel(serviceType))
	if err != nil {
		return nil, err
	}
	if len(results.Items) == 0 {
		return nil, nil
	}

	return &results.Items[0], nil
}

func (sdk *SDK) GetBrokerBindingTemplatesByServiceType(serviceType string) (*templates.BrokerBindingTemplateList, error) {
	return sdk.Templates().BrokerBindingTemplates().List(filterByServiceTypeLabel(serviceType))
}

// RetrieveBindingTemplates lists all binding templates in a namespace.
func (sdk *SDK) RetrieveBindingTemplates(ns, serviceType string) (*templates.BindingTemplateList, error) {
	bndts, err := sdk.Templates().BindingTemplates(ns).List(meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list binding templates in %q (%s)", ns, err)
	}

	return bndts, nil
}

// RetrieveClusterBindingTemplates lists all binding templates in a namespace.
func (sdk *SDK) RetrieveClusterBindingTemplates() (*templates.ClusterBindingTemplateList, error) {
	cbndts, err := sdk.Templates().ClusterBindingTemplates().List(meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list cluster binding templates (%s)", err)
	}

	return cbndts, nil
}

// RetrieveBrokerBindingTemplates lists all binding templates in a namespace.
func (sdk *SDK) RetrieveBrokerBindingTemplates() (*templates.BrokerBindingTemplateList, error) {
	bbndts, err := sdk.Templates().BrokerBindingTemplates().List(meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list broker binding templates (%s)", err)
	}

	return bbndts, nil
}

// RetrieveTemplatedInstance gets an binding template by its name.
func (sdk *SDK) RetrieveBindingTemplate(ns, name string) (*templates.BindingTemplate, error) {
	bndt, err := sdk.Templates().BindingTemplates(ns).Get(name, meta.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get binding template '%s.%s' (%s)", ns, name, err)
	}
	return bndt, nil
}

// RetrieveClusterBindingTemplate gets a cluster binding template by its name.
func (sdk *SDK) RetrieveClusterBindingTemplate(name string) (*templates.ClusterBindingTemplate, error) {
	cbndt, err := sdk.Templates().ClusterBindingTemplates().Get(name, meta.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get cluster binding template %q (%s)", name, err)
	}
	return cbndt, nil
}

// RetrieveBrokerBindingTemplate gets a cluster binding template by its name.
func (sdk *SDK) RetrieveBrokerBindingTemplate(name string) (*templates.BrokerBindingTemplate, error) {
	bbndt, err := sdk.Templates().BrokerBindingTemplates().Get(name, meta.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get broker binding template %q (%s)", name, err)
	}
	return bbndt, nil
}

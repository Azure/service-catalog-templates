// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package servicecatalogtempltesdk

import (
	"fmt"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (sdk *SDK) GetInstanceTemplates(ns string, cluster, broker bool, serviceType string) ([]templates.InstanceTemplateInterface, error) {
	var t []templates.InstanceTemplateInterface

	if broker {
		binstts, err := sdk.GetBrokerInstanceTemplatesByServiceType(serviceType)
		if err != nil {
			return nil, err
		}
		for _, binstt := range binstts.Items {
			t = append(t, &binstt)
		}
		return t, err

	}

	if cluster {
		cinstt, err := sdk.GetClusterInstanceTemplateByServiceType(serviceType)
		if err != nil {
			return nil, err
		}
		if cinstt != nil {
			t = append(t, cinstt)
		}
		return t, err
	}

	instt, err := sdk.GetInstanceTemplateByServiceType(serviceType, ns)
	if err != nil {
		return nil, err
	}
	if instt != nil {
		t = append(t, instt)
	}
	return t, err
}

func (sdk *SDK) GetInstanceTemplate(ns string, cluster, broker bool, name string) (t templates.InstanceTemplateInterface, err error) {
	if broker {
		t, err = sdk.RetrieveBrokerInstanceTemplate(name)
		return t, err
	}

	if cluster {
		t, err = sdk.RetrieveClusterInstanceTemplate(name)
		return t, err
	}

	t, err = sdk.RetrieveInstanceTemplate(ns, name)
	return t, err
}

func (sdk *SDK) GetInstanceTemplateByServiceType(serviceType, namespace string) (*templates.InstanceTemplate, error) {
	results, err := sdk.Templates().InstanceTemplates(namespace).List(filterByServiceTypeLabel(serviceType))
	if err != nil {
		return nil, err
	}
	if len(results.Items) == 0 {
		return nil, nil
	}

	return &results.Items[0], nil
}

func (sdk *SDK) GetClusterInstanceTemplateByServiceType(serviceType string) (*templates.ClusterInstanceTemplate, error) {
	results, err := sdk.Templates().ClusterInstanceTemplates().List(filterByServiceTypeLabel(serviceType))
	if err != nil {
		return nil, err
	}
	if len(results.Items) == 0 {
		return nil, nil
	}

	return &results.Items[0], nil
}

func (sdk *SDK) GetBrokerInstanceTemplatesByServiceType(serviceType string) (*templates.BrokerInstanceTemplateList, error) {
	return sdk.Templates().BrokerInstanceTemplates().List(filterByServiceTypeLabel(serviceType))
}

// RetrieveInstanceTemplates lists all instance templates in a namespace.
func (sdk *SDK) RetrieveInstanceTemplates(ns, serviceType string) (*templates.InstanceTemplateList, error) {
	instts, err := sdk.Templates().InstanceTemplates(ns).List(meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list instance templates in %q (%s)", ns, err)
	}

	return instts, nil
}

// RetrieveClusterInstanceTemplates lists all instance templates in a namespace.
func (sdk *SDK) RetrieveClusterInstanceTemplates() (*templates.ClusterInstanceTemplateList, error) {
	cinstts, err := sdk.Templates().ClusterInstanceTemplates().List(meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list cluster instance templates (%s)", err)
	}

	return cinstts, nil
}

// RetrieveBrokerInstanceTemplates lists all instance templates in a namespace.
func (sdk *SDK) RetrieveBrokerInstanceTemplatesByServiceType() (*templates.BrokerInstanceTemplateList, error) {
	binstts, err := sdk.Templates().BrokerInstanceTemplates().List(meta.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list broker instance templates (%s)", err)
	}

	return binstts, nil
}

// RetrieveTemplatedInstance gets an instance template by its name.
func (sdk *SDK) RetrieveInstanceTemplate(ns, name string) (*templates.InstanceTemplate, error) {
	instt, err := sdk.Templates().InstanceTemplates(ns).Get(name, meta.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get instance template '%s.%s' (%s)", ns, name, err)
	}
	return instt, nil
}

// RetrieveClusterInstanceTemplate gets a cluster instance template by its name.
func (sdk *SDK) RetrieveClusterInstanceTemplate(name string) (*templates.ClusterInstanceTemplate, error) {
	cinstt, err := sdk.Templates().ClusterInstanceTemplates().Get(name, meta.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get cluster instance template %q (%s)", name, err)
	}
	return cinstt, nil
}

// RetrieveBrokerInstanceTemplate gets a cluster instance template by its name.
func (sdk *SDK) RetrieveBrokerInstanceTemplate(name string) (*templates.BrokerInstanceTemplate, error) {
	binstt, err := sdk.Templates().BrokerInstanceTemplates().Get(name, meta.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get broker instance template %q (%s)", name, err)
	}
	return binstt, nil
}

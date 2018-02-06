// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package servicecatalogtemplates

import (
	"fmt"

	"github.com/Azure/service-catalog-templates/pkg/service-catalog-sdk"
	"github.com/Azure/service-catalog-templates/pkg/service-catalog-templates-sdk"
	"github.com/Azure/service-catalog-templates/pkg/service-catalog-templates/builder"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

type resolver struct {
	sdk      *servicecatalogtempltesdk.SDK
	svcatSDK *servicecatalogsdk.SDK
}

func newResolver(sdk *servicecatalogtempltesdk.SDK, svcatSDK *servicecatalogsdk.SDK) *resolver {
	return &resolver{
		sdk:      sdk,
		svcatSDK: svcatSDK,
	}
}

func (r *resolver) ResolveInstanceTemplate(tinst *templates.TemplatedInstance) (templates.InstanceTemplateInterface, error) {
	nsTemplate, err := r.sdk.GetInstanceTemplateByServiceType(tinst.Spec.ServiceType, tinst.Namespace)
	if err != nil {
		return nil, err
	}

	clusterTemplate, err := r.sdk.GetClusterInstanceTemplateByServiceType(tinst.Spec.ServiceType)
	if err != nil {
		return nil, err
	}

	brokerTemplates, err := r.sdk.GetBrokerInstanceTemplatesByServiceType(tinst.Spec.ServiceType)
	if err != nil {
		return nil, err
	}
	var brokerTemplate *templates.BrokerInstanceTemplate
	if len(brokerTemplates.Items) == 1 {
		brokerTemplate = &brokerTemplates.Items[0]
	}

	var template templates.InstanceTemplateInterface
	if nsTemplate == nil && clusterTemplate == nil && brokerTemplate == nil {
		if r.requiresInstanceTemplate(tinst) {
			if len(brokerTemplates.Items) > 1 {
				return nil, fmt.Errorf("more than one broker-level instance template is defined for service type: %s and more specific templates do not exist",
					tinst.Spec.ServiceType)
			}
			return nil, fmt.Errorf("unable to resolve an instance template for service type: %s in namespace: %s",
				tinst.Spec.ServiceType, tinst.Namespace)
		}

		// Just use a blank template since it's okay to use a TemplatedInstance even when you don't need us to resolve a plan
		// i.e. they used to use it and now have picked a plan, or maybe still need it for mapping secret keys, etc.
		template = &templates.InstanceTemplate{}
	} else {
		template, err = r.mergeInstanceTemplates(nsTemplate, clusterTemplate, brokerTemplate)
		if err != nil {
			return nil, err
		}
	}

	// TODO: if a plan selector is specified, pick a different plan from the template's default
	if tinst.Spec.PlanSelector != nil {
		resolvedClass, resolvedPlan, err := r.ResolvePlan(tinst)
		if err != nil {
			return nil, err
		}
		template.SetPlanReference(svcat.PlanReference{
			ClusterServiceClassName: resolvedClass.Name,
			ClusterServicePlanName:  resolvedPlan.Name,
		})
	}

	return template, nil
}

func (r *resolver) requiresInstanceTemplate(inst *templates.TemplatedInstance) bool {
	if (inst.Spec.ClusterServiceClassName != "" || inst.Spec.ClusterServiceClassExternalName != "") &&
		(inst.Spec.ClusterServicePlanName != "" || inst.Spec.ClusterServicePlanExternalName != "") {
		return false
	}

	return true
}

func (r *resolver) ResolveBindingTemplate(tbnd templates.TemplatedBinding) (templates.BindingTemplateInterface, error) {
	tinst, err := r.sdk.GetInstanceFromCache(tbnd.Namespace, tbnd.Spec.TemplatedInstanceRef.Name)
	if err != nil {
		return nil, err
	}

	nsTemplate, err := r.sdk.GetBindingTemplateByServiceType(tinst.Spec.ServiceType, tinst.Namespace)
	if err != nil {
		return nil, err
	}

	clusterTemplate, err := r.sdk.GetClusterBindingTemplateByServiceType(tinst.Spec.ServiceType)
	if err != nil {
		return nil, err
	}

	brokerTemplates, err := r.sdk.GetBrokerBindingTemplatesByServiceType(tinst.Spec.ServiceType)
	if err != nil {
		return nil, err
	}
	var brokerTemplate *templates.BrokerBindingTemplate
	if len(brokerTemplates.Items) == 1 {
		brokerTemplate = &brokerTemplates.Items[0]
	}

	var template templates.BindingTemplateInterface
	if nsTemplate == nil && clusterTemplate == nil && brokerTemplate == nil {
		// Just use a blank template
		template = &templates.BindingTemplate{}
	} else {
		template, err = r.mergeBindingTemplates(nsTemplate, clusterTemplate, brokerTemplate)
		if err != nil {
			return nil, err
		}
	}

	return template, nil
}

func (r *resolver) ResolvePlan(instance *templates.TemplatedInstance) (*svcat.ClusterServiceClass, *svcat.ClusterServicePlan, error) {
	// TODO: using the plan selector and type select a matching plan
	return nil, nil, nil
}

func (r *resolver) mergeInstanceTemplates(namespaceTemplate *templates.InstanceTemplate,
	clusterTemplate *templates.ClusterInstanceTemplate, brokerTemplate *templates.BrokerInstanceTemplate,
) (templates.InstanceTemplateInterface, error) {
	template := &templates.InstanceTemplate{}

	if brokerTemplate != nil {
		template.Spec.PlanReference = brokerTemplate.Spec.PlanReference
		template.Spec.Parameters = brokerTemplate.Spec.Parameters
		template.Spec.ParametersFrom = brokerTemplate.Spec.ParametersFrom
	}

	var err error
	if clusterTemplate != nil {
		template.Spec.Parameters, err = builder.MergeParameters(template.Spec.Parameters, clusterTemplate.Spec.Parameters)
		if err != nil {
			return nil, err
		}
		template.Spec.ParametersFrom = builder.MergeParametersFromSource(template.Spec.ParametersFrom, clusterTemplate.Spec.ParametersFrom)
		template.Spec.PlanReference = builder.MergePlanReference(template.Spec.PlanReference, clusterTemplate.Spec.PlanReference)
	}

	if namespaceTemplate != nil {
		template.Spec.Parameters, err = builder.MergeParameters(template.Spec.Parameters, namespaceTemplate.Spec.Parameters)
		if err != nil {
			return nil, err
		}
		template.Spec.ParametersFrom = builder.MergeParametersFromSource(template.Spec.ParametersFrom, namespaceTemplate.Spec.ParametersFrom)
		template.Spec.PlanReference = builder.MergePlanReference(template.Spec.PlanReference, namespaceTemplate.Spec.PlanReference)
	}

	return template, nil
}

func (r *resolver) mergeBindingTemplates(namespaceTemplate *templates.BindingTemplate,
	clusterTemplate *templates.ClusterBindingTemplate, brokerTemplate *templates.BrokerBindingTemplate,
) (templates.BindingTemplateInterface, error) {
	template := &templates.BindingTemplate{}

	if brokerTemplate != nil {
		template.Spec.Parameters = brokerTemplate.Spec.Parameters
		template.Spec.ParametersFrom = brokerTemplate.Spec.ParametersFrom
		template.Spec.SecretKeys = brokerTemplate.Spec.SecretKeys
	}

	var err error
	if clusterTemplate != nil {
		template.Spec.Parameters, err = builder.MergeParameters(template.Spec.Parameters, clusterTemplate.Spec.Parameters)
		if err != nil {
			return nil, err
		}
		template.Spec.ParametersFrom = builder.MergeParametersFromSource(template.Spec.ParametersFrom, clusterTemplate.Spec.ParametersFrom)
		template.Spec.SecretKeys = builder.MergeSecretKeys(template.Spec.SecretKeys, clusterTemplate.Spec.SecretKeys)
	}

	if namespaceTemplate != nil {
		template.Spec.Parameters, err = builder.MergeParameters(template.Spec.Parameters, namespaceTemplate.Spec.Parameters)
		if err != nil {
			return nil, err
		}
		template.Spec.ParametersFrom = builder.MergeParametersFromSource(template.Spec.ParametersFrom, namespaceTemplate.Spec.ParametersFrom)
		template.Spec.SecretKeys = builder.MergeSecretKeys(template.Spec.SecretKeys, namespaceTemplate.Spec.SecretKeys)
	}

	return template, nil
}

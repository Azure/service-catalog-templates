package svcatt

import (
	"fmt"

	"github.com/Azure/service-catalog-templates/pkg/sdk"
	"github.com/Azure/service-catalog-templates/pkg/svcatsdk"
	"github.com/Azure/service-catalog-templates/pkg/svcatt/builder"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

type resolver struct {
	sdk      *sdk.SDK
	svcatSDK *svcatsdk.SDK
}

func newResolver(sdk *sdk.SDK, svcatSDK *svcatsdk.SDK) *resolver {
	return &resolver{
		sdk:      sdk,
		svcatSDK: svcatSDK,
	}
}

func (r *resolver) ResolveInstanceTemplate(instance templates.TemplatedInstance) (templates.InstanceTemplateInterface, error) {
	nsTemplate, err := r.sdk.GetInstanceTemplateByServiceType(instance.Spec.ServiceType, instance.Namespace)
	if err != nil {
		return nil, err
	}

	clusterTemplate, err := r.sdk.GetClusterInstanceTemplateByServiceType(instance.Spec.ServiceType)
	if err != nil {
		return nil, err
	}

	brokerTemplate, err := r.sdk.GetBrokerInstanceTemplateByServiceType(instance.Spec.ServiceType)
	if err != nil {
		return nil, err
	}

	if nsTemplate == nil && clusterTemplate == nil && brokerTemplate == nil {
		return nil, fmt.Errorf("unable to resolve an instance template for service type: %s in namespace: %s",
			instance.Spec.ServiceType, instance.Namespace)
	}

	template, err := r.mergeInstanceTemplates(nsTemplate, clusterTemplate, brokerTemplate)
	if err != nil {
		return nil, err
	}

	// TODO: if a plan selector is specified, pick a different plan from the template's default
	if instance.Spec.PlanSelector != nil {
		resolvedClass, resolvedPlan, err := r.ResolvePlan(instance)
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

func (r *resolver) ResolveBindingTemplate(tbnd templates.TemplatedBinding) (template templates.BindingTemplateInterface, err error) {
	tinst, err := r.sdk.GetInstanceFromCache(tbnd.Namespace, tbnd.Spec.InstanceRef.Name)
	if err != nil {
		return nil, err
	}

	template, err = r.sdk.GetBindingTemplateByServiceType(tinst.Spec.ServiceType, tbnd.Namespace)
	if err != nil {
		return nil, err
	}
	if template == nil {
		template, err = r.sdk.GetClusterBindingTemplateByServiceType(tinst.Spec.ServiceType)
		if err != nil {
			return nil, err
		}
	}
	if template == nil {
		template, err = r.sdk.GetBrokerBindingTemplateByServiceType(tinst.Spec.ServiceType)
		if err != nil {
			return nil, err
		}
	}

	if template == nil {
		return nil, fmt.Errorf("unable to resolve an tbnd template for service type: %s in namespace: %s",
			tinst.Spec.ServiceType, tbnd.Namespace)
	}

	return template, nil
}

func (r *resolver) ResolvePlan(instance templates.TemplatedInstance) (*svcat.ClusterServiceClass, *svcat.ClusterServicePlan, error) {
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

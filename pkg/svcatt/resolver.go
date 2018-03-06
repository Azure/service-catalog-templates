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

func (r *resolver) ResolveInstanceTemplate(instance *templates.TemplatedInstance) (templates.InstanceTemplateInterface, error) {
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

	var template templates.InstanceTemplateInterface
	if nsTemplate == nil && clusterTemplate == nil && brokerTemplate == nil {
		if r.requiresInstanceTemplate(instance) {
			return nil, fmt.Errorf("unable to resolve an instance template for service type: %s in namespace: %s",
				instance.Spec.ServiceType, instance.Namespace)
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

func (r *resolver) requiresInstanceTemplate(inst *templates.TemplatedInstance) bool {
	if (inst.Spec.ClusterServiceClassName != "" || inst.Spec.ClusterServiceClassExternalName != "") &&
		(inst.Spec.ClusterServicePlanName != "" || inst.Spec.ClusterServicePlanExternalName != "") {
		return false
	}

	return true
}

func (r *resolver) ResolveBindingTemplate(tbnd templates.TemplatedBinding) (templates.BindingTemplateInterface, error) {
	tinst, err := r.sdk.GetInstanceFromCache(tbnd.Namespace, tbnd.Spec.InstanceRef.Name)
	if err != nil {
		return nil, err
	}

	var template templates.BindingTemplateInterface

	bndt, err := r.sdk.GetBindingTemplateByServiceType(tinst.Spec.ServiceType, tbnd.Namespace)
	if err != nil {
		return nil, err
	}
	if bndt != nil {
		template = bndt
	}

	if template == nil {
		cbndt, err := r.sdk.GetClusterBindingTemplateByServiceType(tinst.Spec.ServiceType)
		if err != nil {
			return nil, err
		}
		if cbndt != nil {
			template = cbndt
		}
	}

	if template == nil {
		bbndt, err := r.sdk.GetBrokerBindingTemplateByServiceType(tinst.Spec.ServiceType)
		if err != nil {
			return nil, err
		}
		if bbndt != nil {
			template = bbndt
		}
	}

	if template == nil {
		return nil, fmt.Errorf("unable to resolve an binding template for service type: %s in namespace: %s",
			tinst.Spec.ServiceType, tbnd.Namespace)
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

package svcatt

import (
	"fmt"

	"github.com/Azure/service-catalog-templates/pkg/sdk"
	"github.com/Azure/service-catalog-templates/pkg/svcatsdk"
	"k8s.io/apimachinery/pkg/labels"

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

func (r *resolver) ResolveInstanceTemplate(instance templates.TemplatedInstance) (*templates.InstanceTemplate, error) {
	opts := labels.SelectorFromSet(map[string]string{templates.FieldServiceTypeName: instance.Spec.ServiceType})
	results, err := r.sdk.InstanceTemplateCache().InstanceTemplates(instance.Namespace).List(opts)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("unable to resolve an instance template for service type: %s", instance.Spec.ServiceType)
	}

	template := results[0].DeepCopy()

	// TODO: if a plan selector is specified, pick a different plan from the template's default
	if instance.Spec.PlanSelector != nil {
		resolvedClass, resolvedPlan, err := r.ResolvePlan(instance)
		if err != nil {
			return nil, err
		}
		// TODO: track the uuid instead of resolving then forcing a second lookup
		template.Spec.ClassExternalName = resolvedClass.Spec.ExternalName
		template.Spec.PlanExternalName = resolvedPlan.Spec.ExternalName
	}

	return template, nil
}

func (r *resolver) ResolveBindingTemplate(binding templates.TemplatedBinding) (*templates.BindingTemplate, error) {
	inst, err := r.sdk.GetInstanceFromCache(binding.Namespace, binding.Spec.InstanceRef.Name)
	if err != nil {
		return nil, err
	}

	opts := labels.SelectorFromSet(map[string]string{templates.FieldServiceTypeName: inst.Spec.ServiceType})
	results, err := r.sdk.BindingTemplateCache().BindingTemplates(binding.Namespace).List(opts)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("unable to resolve a binding template for service type: %s", inst.Spec.ServiceType)
	}

	template := results[0].DeepCopy()
	return template, nil
}

func (r *resolver) ResolvePlan(instance templates.TemplatedInstance) (*svcat.ClusterServiceClass, *svcat.ClusterServicePlan, error) {
	// TODO: using the plan selector and type select a matching plan
	return nil, nil, nil
}

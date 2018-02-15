package svcatt

import (
	"fmt"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	templates "github.com/Azure/service-catalog-templates/pkg/apis/templates/experimental"
	templatesclient "github.com/Azure/service-catalog-templates/pkg/client/clientset/versioned"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	svcatclient "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
)

type resolver struct {
	templatesClient templatesclient.Interface
	svcatClient     svcatclient.Interface
}

func newResolver(templatesClient templatesclient.Interface, svcatClient svcatclient.Interface) *resolver {
	return &resolver{
		templatesClient: templatesClient,
		svcatClient:     svcatClient,
	}
}

func (r *resolver) ResolveTemplate(instance templates.Instance) (*templates.InstanceTemplate, error) {
	opts := meta.ListOptions{
		LabelSelector: labels.FormatLabels(map[string]string{templates.FieldServiceTypeName: instance.Spec.ServiceType}),
	}
	results, err := r.templatesClient.TemplatesExperimental().InstanceTemplates(instance.Namespace).List(opts)
	if err != nil {
		return nil, err
	}
	if len(results.Items) == 0 {
		return nil, fmt.Errorf("unable to resolve a template for service type: %s", instance.Spec.ServiceType)
	}

	template := results.Items[0]

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

	return &template, nil
}

func (r *resolver) ResolvePlan(instance templates.Instance) (*svcat.ClusterServiceClass, *svcat.ClusterServicePlan, error) {
	// TODO: using the plan selector and type select a matching plan
	return nil, nil, nil
}

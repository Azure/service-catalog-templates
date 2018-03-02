package builder

import (
	"encoding/json"
	"fmt"

	"github.com/peterbourgon/mergemap"
	"k8s.io/apimachinery/pkg/runtime"

	svcat "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
)

func MergeParameters(instParams *runtime.RawExtension, tmplParams *runtime.RawExtension) (*runtime.RawExtension, error) {
	if tmplParams == nil {
		return instParams, nil
	}

	if instParams == nil {
		return tmplParams, nil
	}

	var instMap, tmplMap map[string]interface{}
	json.Unmarshal(instParams.Raw, &instMap)
	json.Unmarshal(tmplParams.Raw, &tmplMap)

	merged := mergemap.Merge(instMap, tmplMap)

	result, err := json.Marshal(merged)
	if err != nil {
		return nil, fmt.Errorf("could not merge the instance and template parameters: %s", err)
	}

	return &runtime.RawExtension{Raw: result}, nil
}

func MergeParametersFromSource(instParams []svcat.ParametersFromSource, tmplParams []svcat.ParametersFromSource) []svcat.ParametersFromSource {
	// TODO: I don't believe that merging is the right thing, so I'm only using the template if the instance didn't define anything
	if len(instParams) == 0 {
		return tmplParams
	}

	return instParams
}

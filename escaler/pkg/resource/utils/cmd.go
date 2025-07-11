package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/Emerging-AI/ENOVA/escaler/pkg/config"

	"github.com/Emerging-AI/ENOVA/escaler/pkg/meta"
)

func shouldAppend(v interface{}) bool {
	switch v := v.(type) {
	case int, int32, int64:
		return v != 0
	case float32, float64:
		return v != 0.0
	case string:
		return v != ""
	case bool:
		return v // no need to check, because false is the zero value and means "not set"
	default:
		// This case is for types not explicitly checked above; assumes non-zero by default
		return !reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
	}
}

func BuildCmdFromTaskSpec(spec meta.TaskSpec) []string {

	cmd := []string{
		"enova", "serving", "run", "--model", spec.Model, "--port", strconv.Itoa(spec.Port), "--host", spec.Host,
		"--backend", spec.Backend,
		"--exporter_service_name", spec.ExporterServiceName,
	}
	if config.GetEConfig().ResourceBackend.Type == config.ResourceBackendTypeK8s {
		cmd = append(cmd, "--exporter_endpoint", spec.Name+"-collector."+spec.Namespace+".svc.cluster.local:4317")
	} else {
		cmd = append(cmd, "--exporter_endpoint", spec.ExporterEndpoint)
	}

	switch spec.Backend {
	case "vllm":
		cmd = UpdateCmdByBackendConfig[*meta.VllmBackendConfig](cmd, spec)
	case "sglang":
		cmd = UpdateCmdByBackendConfig[*meta.SglangBackendConfig](cmd, spec)
	}
	// Add extra serving params
	for k, v := range spec.BackendExtraConfig {
		cmd = append(cmd, []string{fmt.Sprintf("--%s", k), fmt.Sprintf("%v", v)}...)
	}
	return cmd
}

func UpdateCmdByBackendConfig[B interface{}](cmd []string, spec meta.TaskSpec) []string {
	backendConfig, ok := spec.BackendConfig.(B)
	if ok {
		jsonBytes, err := json.Marshal(backendConfig)
		if err != nil {

		} else {
			var backendConfigMap map[string]interface{}
			err = json.Unmarshal(jsonBytes, &backendConfigMap)
			if err != nil {

			} else {
				// if there is not valid value, dont append to cmd params
				for k, v := range backendConfigMap {
					if shouldAppend(v) {
						cmd = append(cmd, []string{fmt.Sprintf("--%s", k), fmt.Sprintf("%v", v)}...)
					}
				}
			}
		}
	}
	return cmd
}

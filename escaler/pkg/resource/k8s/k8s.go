package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	rscutils "github.com/Emerging-AI/ENOVA/escaler/pkg/resource/utils"

	"sigs.k8s.io/controller-runtime/pkg/client"

	k8sresource "k8s.io/apimachinery/pkg/api/resource"

	"github.com/Emerging-AI/ENOVA/escaler/pkg/logger"
	"github.com/Emerging-AI/ENOVA/escaler/pkg/meta"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/mitchellh/mapstructure"
	otalv1 "github.com/open-telemetry/opentelemetry-operator/apis/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

const annotationRestarted = "restartedAt"

var collectorServiceAccount = "otel-collector"

type K8sCli struct {
	K8sClient     *kubernetes.Clientset
	DynamicClient *dynamic.DynamicClient
	Ctx           context.Context
}

type Workload struct {
	K8sCli *K8sCli
	Spec   *meta.TaskSpec
}

func (w *Workload) CreateOrUpdate() {
	dp, err := w.GetWorkload()
	if err != nil {
		if apierrors.IsNotFound(err) {
			_, err := w.CreateWorkload()
			if err != nil {
				logger.Error(err, "Failed to create workload")
				return
			}
		}
		logger.Error(err, "workload GetWorkload failed")
	} else {
		_, err := w.UpdateWorkload(dp)
		if err != nil {
			logger.Error(err, "Failed to update workload")
			return
		}
	}

	_, err = w.GetService()
	if err != nil {
		logger.Debugf("K8sResourceClient DeployTask check service get error: %v", err)
		_, err = w.CreateService()
		if err != nil {
			return
		}
	} else {
		_, err = w.UpdateService()
		if err != nil {
			return
		}
	}

	if w.Spec.Collector.Enable && !w.isCustomized() {
		// TODO:resourceVersion not found
		ot, err := w.GetCollector()
		if err != nil {
			_, err = w.CreateCollector()
			if err != nil {
				return
			}
		} else {
			_, err = w.UpdateCollector(ot)
			if err != nil {
				return
			}
		}
	}
}

// Create 1. create workload, 2. create service
func (w *Workload) Create() {
	_, err := w.CreateWorkload()
	if err != nil {
		return
	}
	_, err = w.CreateService()
	if err != nil {
		return
	}
}

func (w *Workload) Update() {
	dp, err := w.GetWorkload()
	if err != nil {
		return
	}
	_, err = w.UpdateWorkload(dp)
	if err != nil {
		return
	}
	_, err = w.UpdateService()
	if err != nil {
		return
	}
}

func (w *Workload) Delete() {
	_ = w.DeleteWorkload()
	_ = w.DeleteService()
	_ = w.DeleteCollector()
}

func (w *Workload) CreateWorkload() (*v1.Deployment, error) {
	deployment := w.buildDeployment()
	taskSpecJson, err := json.Marshal(deployment)
	logger.Infof("Create deployment: %s", string(taskSpecJson))
	opts := metav1.CreateOptions{}
	ret, err := w.K8sCli.K8sClient.AppsV1().Deployments(w.Spec.Namespace).Create(w.K8sCli.Ctx, &deployment, opts)
	if err != nil {
		logger.Errorf("Workload CreateWorkload error: %v", err)
		return ret, err
	}
	return ret, nil
}

func (w *Workload) GetWorkload() (*v1.Deployment, error) {
	deployment := w.buildDeployment()
	ret, err := w.K8sCli.K8sClient.AppsV1().Deployments(w.Spec.Namespace).Get(w.K8sCli.Ctx, deployment.Name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("Workload GetWorkload error: %v", err)
		return ret, err
	}
	return ret, nil
}

func (w *Workload) UpdateWorkload(dp *v1.Deployment) (*v1.Deployment, error) {
	deployment := w.buildDeployment()
	taskSpecJson, err := json.Marshal(deployment)
	logger.Infof("Update deployment: %s", string(taskSpecJson))
	opts := metav1.UpdateOptions{}
	deployment.ResourceVersion = dp.ResourceVersion
	ret, err := w.K8sCli.K8sClient.AppsV1().Deployments(w.Spec.Namespace).Update(w.K8sCli.Ctx, &deployment, opts)
	if err != nil {
		logger.Errorf("Workload UpdateWorkload error: %v", err)
		return ret, err
	}
	return ret, nil
}

func (w *Workload) UpdateWorkloadWithPatch(patchData []byte) (*v1.Deployment, error) {
	dp, err := w.GetWorkload()
	if err != nil {
		return nil, err
	}

	ret, err := w.K8sCli.K8sClient.AppsV1().Deployments(w.Spec.Namespace).Patch(
		w.K8sCli.Ctx,
		dp.Name,
		types.JSONPatchType,
		patchData,
		metav1.PatchOptions{},
	)
	if err != nil {
		logger.Errorf("Workload UpdateWorkload error: %v", err)
		return ret, err
	}
	return ret, nil
}

func (w *Workload) DeleteWorkload() error {
	if err := w.K8sCli.K8sClient.AppsV1().Deployments(w.Spec.Namespace).Delete(w.K8sCli.Ctx, w.Spec.Name, metav1.DeleteOptions{}); client.IgnoreNotFound(err) != nil {
		logger.Errorf("Workload DeleteWorkload error: %v", err)
		return err
	}
	return nil
}

func (w *Workload) CreateService() (*corev1.Service, error) {
	opts := metav1.CreateOptions{}
	service := w.buildService()
	ret, err := w.K8sCli.K8sClient.CoreV1().Services(w.Spec.Namespace).Create(w.K8sCli.Ctx, &service, opts)
	if err != nil {
		logger.Errorf("Workload CreateService error: %v", err)
		return ret, err
	}
	return ret, nil
}

func (w *Workload) UpdateService() (*corev1.Service, error) {
	opts := metav1.UpdateOptions{}
	service := w.buildService()
	ret, err := w.K8sCli.K8sClient.CoreV1().Services(w.Spec.Namespace).Update(w.K8sCli.Ctx, &service, opts)
	if err != nil {
		logger.Errorf("Workload UpdateService error: %v", err)
		return ret, err
	}
	return ret, nil
}

func (w *Workload) DeleteService() error {
	if err := w.K8sCli.K8sClient.CoreV1().Services(w.Spec.Namespace).Delete(w.K8sCli.Ctx, w.buildService().Name, metav1.DeleteOptions{}); client.IgnoreNotFound(err) != nil {
		logger.Errorf("Workload DeleteService error: %v", err)
		return err
	}
	return nil
}

func (w *Workload) GetConfigMap(name string) (*corev1.ConfigMap, error) {
	ret, err := w.K8sCli.K8sClient.CoreV1().ConfigMaps(w.Spec.Namespace).Get(w.K8sCli.Ctx, name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("Workload GetConfigMap error: %v", err)
		return ret, err
	}
	return ret, nil
}

func (w *Workload) buildDeployment() v1.Deployment {
	replicas := int32(w.Spec.Replica)
	matchLabels := make(map[string]string)
	matchLabels["enovaserving-name"] = w.Spec.Name
	matchLabels["app"] = w.Spec.Name
	matchLabels["version"] = "v1.0.0"

	env := make([]corev1.EnvVar, len(w.Spec.Envs))
	for i, e := range w.Spec.Envs {
		env[i] = corev1.EnvVar{
			Name:  e.Name,
			Value: e.Value,
		}
	}

	// imagePullSecrets
	var imagePullSecrets []corev1.LocalObjectReference
	for _, s := range w.Spec.ImagePullSecrets {
		if s != "" {
			imagePullSecrets = append(imagePullSecrets, corev1.LocalObjectReference{Name: s})
		}
	}

	livenessProbe := corev1.Probe{}
	readinessProbe := corev1.Probe{}
	probe := corev1.Probe{ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{Path: "/health",
		Port: intstr.IntOrString{IntVal: int32(w.Spec.Port)}}}, InitialDelaySeconds: 30}
	switch w.Spec.Backend {
	case "vllm", "sglang":
		// TODO: custom health
		if !w.isCustomized() {
			livenessProbe = probe
			livenessProbe.FailureThreshold = 3
			livenessProbe.InitialDelaySeconds = 60
			livenessProbe.TimeoutSeconds = 5
			readinessProbe = probe
			readinessProbe.FailureThreshold = 3
			readinessProbe.InitialDelaySeconds = 60
			readinessProbe.TimeoutSeconds = 5
		}
	}

	// cpu,memory request & limit
	resources := corev1.ResourceRequirements{
		Requests: corev1.ResourceList{},
		Limits:   corev1.ResourceList{},
	}

	// cpu,memory request & limit
	if w.Spec.Resources.CPU != "" {
		q := k8sresource.MustParse(w.Spec.Resources.CPU)
		resources.Requests[corev1.ResourceCPU] = q
		resources.Limits[corev1.ResourceCPU] = q
	}

	if w.Spec.Resources.Memory != "" {
		q := k8sresource.MustParse(w.Spec.Resources.Memory)
		resources.Requests[corev1.ResourceMemory] = q
		resources.Limits[corev1.ResourceMemory] = q
	}

	if w.Spec.Resources.GPU != "" {
		q := k8sresource.MustParse(w.Spec.Resources.GPU)
		resources.Requests[corev1.ResourceName("nvidia.com/gpu")] = q
		resources.Limits[corev1.ResourceName("nvidia.com/gpu")] = q
	}

	// default mount ~/.cache to host data disk
	deployment := v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        w.Spec.Name,
			Namespace:   w.Spec.Namespace,
			Labels:      matchLabels,
			Annotations: w.Spec.Annotations,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: matchLabels,
			},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image:           w.Spec.Image,
							ImagePullPolicy: corev1.PullAlways,
							Name:            w.Spec.Name,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: int32(w.Spec.Port),
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env:       env,
							Resources: resources,
						},
					},
					ImagePullSecrets: imagePullSecrets,
				},
			},
		},
	}

	if w.isCustomized() {
		if len(w.Spec.Command) > 0 {
			deployment.Spec.Template.Spec.Containers[0].Command = w.Spec.Command
			deployment.Spec.Template.Spec.Containers[0].Args = w.Spec.Args
		}
	} else {
		cmd := rscutils.BuildCmdFromTaskSpec(*w.Spec)
		deployment.Spec.Template.Spec.Containers[0].Command = cmd[:1]
		deployment.Spec.Template.Spec.Containers[0].Args = cmd[1:]
	}

	deployment.Spec.Template.Annotations = w.Spec.Annotations
	volumes := make([]corev1.Volume, 0)
	volumeMounts := make([]corev1.VolumeMount, 0)

	for _, v := range w.Spec.Volumes {
		volumeSource := corev1.VolumeSource{}
		switch v.Type {
		case "emptyDir":
			volumeSource.EmptyDir = &corev1.EmptyDirVolumeSource{}
			break
		case "hostPath":
			volumeSource.HostPath = &corev1.HostPathVolumeSource{Path: v.Path}
			break
		case "NFS":
			volumeSource.NFS = &corev1.NFSVolumeSource{Server: v.Value, Path: v.Path}
			break
		default:
			continue
		}
		volumes = append(volumes, corev1.Volume{Name: v.Name, VolumeSource: volumeSource})
	}

	for _, v := range w.Spec.VolumeMounts {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      v.Name,
			MountPath: v.Path,
			ReadOnly:  v.ReadOnly,
		})
	}

	if !w.isCustomized() {
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &readinessProbe
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &livenessProbe
	}
	// it will add shm by default
	shmLimitSize := k8sresource.MustParse("1Gi")
	volumes = append(volumes, corev1.Volume{
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium:    corev1.StorageMediumMemory,
				SizeLimit: &shmLimitSize,
			},
		},
		Name: "shm",
	})
	volumeMounts = append(volumeMounts, corev1.VolumeMount{
		Name:      "shm",
		MountPath: "/dev/shm",
	})

	// ConfigMaps
	for _, v := range w.Spec.ConfigMaps {
		if configMap, err := w.GetConfigMap(v.Name); err == nil {
			volumes = append(volumes, corev1.Volume{
				Name: configMap.Name,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{Name: configMap.Name},
					},
				},
			})
			volumeMounts = append(volumeMounts, corev1.VolumeMount{
				Name:      configMap.Name,
				MountPath: v.Path,
			})
		}
	}

	maxUnavailable := intstr.FromString("25%")
	maxSurge := intstr.FromString("0%")
	deployment.Spec.Strategy = v1.DeploymentStrategy{
		Type:          v1.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &v1.RollingUpdateDeployment{MaxUnavailable: &maxUnavailable, MaxSurge: &maxSurge},
	}
	deployment.Spec.Template.Spec.Volumes = volumes
	deployment.Spec.Template.Spec.Containers[0].VolumeMounts = volumeMounts
	if len(w.Spec.NodeSelector) > 0 {
		deployment.Spec.Template.Spec.NodeSelector = w.Spec.NodeSelector
	}

	deployment.Spec.Template.Labels = matchLabels
	deployment.Labels = matchLabels
	return deployment
}

func (w *Workload) buildService() corev1.Service {
	selector := make(map[string]string)
	selector["enovaserving-name"] = w.Spec.Name
	ports := make([]corev1.ServicePort, len(w.Spec.Service.Ports))
	for i, p := range w.Spec.Service.Ports {
		ports[i] = corev1.ServicePort{
			Name:     fmt.Sprintf("tcp%d", i),
			Protocol: corev1.ProtocolTCP,
			Port:     p.Number,
			TargetPort: intstr.IntOrString{
				IntVal: p.Number,
			},
		}
	}

	service := corev1.Service{
		Spec: corev1.ServiceSpec{
			Selector: selector,
			Ports:    ports,
		},
	}
	service.Name = fmt.Sprintf("%s-svc", w.Spec.Name)
	service.Namespace = w.Spec.Namespace
	return service
}

func formatBrokers(brokers []string) string {
	if len(brokers) == 0 {
		return "[]"
	}
	formatted := `["`
	for i, broker := range brokers {
		if i > 0 {
			formatted += `", "`
		}
		formatted += broker
	}
	formatted += `"]`
	return formatted
}

func (w *Workload) buildCollector() otalv1.OpenTelemetryCollector {
	actions := []interface{}{
		map[string]interface{}{
			"key":    "cluster_id",
			"action": "insert",
			"value":  w.Spec.Collector.ClusterId,
		},
	}

	if len(w.Spec.Collector.CustomMetricsAdd) > 0 {
		for k, v := range w.Spec.Collector.CustomMetricsAdd {
			actions = append(actions, map[string]interface{}{
				"key":    k,
				"action": "insert",
				"value":  v,
			})
		}
	}
	processors := otalv1.AnyConfig{Object: map[string]interface{}{
		"batch": map[string]interface{}{},
		"attributes/metrics": map[string]interface{}{
			"actions": actions,
		},
		"attributes/http": map[string]interface{}{
			"actions": []interface{}{
				map[string]interface{}{
					"action": "delete",
					"key":    "http.server_name",
				},
				map[string]interface{}{
					"action": "delete",
					"key":    "http.host",
				},
			},
		},
	},
	}

	serviceProcessors := []string{"attributes/metrics", "attributes/http", "batch"}
	if w.Spec.Backend == "sglang" {
		processors.Object["metricstransform"] = map[string]interface{}{
			"transforms": []interface{}{
				map[string]interface{}{
					"action":     "update",
					"include":    "^sglang:num_queue_reqs$$",
					"match_type": "regexp",
					"new_name":   "vllm_num_requests_waiting",
				},
				map[string]interface{}{
					"action":     "update",
					"include":    "^sglang:num_running_reqs$$",
					"match_type": "regexp",
					"new_name":   "vllm_num_requests_running",
				},
				map[string]interface{}{
					"action":     "update",
					"include":    "^sglang:(.*)$$",
					"match_type": "regexp",
					"new_name":   "vllm:$${1}]",
				},
			},
		}
		serviceProcessors = []string{"attributes/metrics", "attributes/http", "metricstransform", "batch"}
	}

	service := otalv1.Service{
		Extensions: nil,
		Telemetry:  nil,
		Pipelines: map[string]*otalv1.Pipeline{
			"traces": {
				Receivers:  []string{"otlp"},
				Processors: []string{"batch"},
				Exporters:  []string{"kafka"},
			},
			"metrics": {
				Receivers:  []string{"prometheus", "otlp"},
				Processors: serviceProcessors,
				Exporters:  []string{"kafka"},
			},
		},
	}

	collector := otalv1.OpenTelemetryCollector{
		Spec: otalv1.OpenTelemetryCollectorSpec{
			OpenTelemetryCommonFields: otalv1.OpenTelemetryCommonFields{ServiceAccount: collectorServiceAccount},
			Config:                    otalv1.Config{Processors: &processors, Service: service},
		},
	}
	collector.Name = w.Spec.Name
	collector.Namespace = w.Spec.Namespace
	return collector
}

func (w *Workload) GetDeployment() (*v1.Deployment, error) {
	return w.K8sCli.K8sClient.AppsV1().Deployments(w.Spec.Namespace).Get(w.K8sCli.Ctx, w.Spec.Name, metav1.GetOptions{})
}

func (w *Workload) GetPodsList() (*corev1.PodList, error) {
	opts := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("enovaserving-name=%s", w.Spec.Name),
	}
	return w.K8sCli.K8sClient.CoreV1().Pods(w.Spec.Namespace).List(w.K8sCli.Ctx, opts)
}

func (w *Workload) GetService() (*corev1.Service, error) {
	opts := metav1.GetOptions{}
	service := w.buildService()
	ret, err := w.K8sCli.K8sClient.CoreV1().Services(w.Spec.Namespace).Get(w.K8sCli.Ctx, service.Name, opts)
	if err != nil {
		logger.Errorf("Workload GetService error: %v", err)
		return ret, err
	}
	return ret, nil
}

func (w *Workload) GetCollector() (*otalv1.OpenTelemetryCollector, error) {
	collector := otalv1.OpenTelemetryCollector{}
	rsc := w.GetOtCollectorResource()
	ret, err := rsc.Namespace(w.Spec.Namespace).Get(w.K8sCli.Ctx, w.Spec.Name, metav1.GetOptions{})
	if err != nil {
		logger.Errorf("GetCollector Get error: %v", err)
		return &collector, err
	}
	_ = mapstructure.Decode(ret.Object, &collector)
	return &collector, err
}

func (w *Workload) CreateCollector() (otalv1.OpenTelemetryCollector, error) {
	collector := w.buildCollector()

	obj := w.buildCollectorUnstructued(collector)

	rsc := w.GetOtCollectorResource()
	ret, err := rsc.Namespace(w.Spec.Namespace).Create(w.K8sCli.Ctx, &obj, metav1.CreateOptions{})
	if err != nil {
		logger.Errorf("CreateCollector Create error: %v", err)
		return collector, err
	}
	_ = mapstructure.Decode(&ret.Object, &collector)
	return collector, err
}

func (w *Workload) DeleteCollector() error {
	rsc := w.GetOtCollectorResource()
	if err := rsc.Namespace(w.Spec.Namespace).Delete(w.K8sCli.Ctx, w.Spec.Name, metav1.DeleteOptions{}); client.IgnoreNotFound(err) != nil {
		return err
	}
	return nil
}

func (w *Workload) UpdateCollector(ot *otalv1.OpenTelemetryCollector) (otalv1.OpenTelemetryCollector, error) {
	collector := w.buildCollector()

	obj := w.buildCollectorUnstructued(collector)
	_ = unstructured.SetNestedField(obj.Object, ot.ResourceVersion, "metadata", "resourceVersion")

	logger.Infof("Update Collector: %+v\n", obj.Object)

	rsc := w.GetOtCollectorResource()
	ret, err := rsc.Namespace(w.Spec.Namespace).Update(w.K8sCli.Ctx, &obj, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("UpdateCollector Update error: %v", err)
		return collector, err
	}
	_ = mapstructure.Decode(&ret.Object, &collector)
	return collector, err
}

func (w *Workload) buildCollectorUnstructued(collector otalv1.OpenTelemetryCollector) unstructured.Unstructured {
	return unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "opentelemetry.io/v1beta1",
			"kind":       "OpenTelemetryCollector",
			"metadata": map[string]interface{}{
				"name":      collector.Name,
				"namespace": collector.Namespace,
			},
			"spec": map[string]interface{}{
				"serviceAccount": collector.Spec.ServiceAccount,
				"config": map[string]interface{}{
					"receivers": map[string]interface{}{
						"otlp": map[string]interface{}{
							"protocols": map[string]interface{}{
								"grpc": map[string]interface{}{
									"endpoint": "0.0.0.0:4317",
								},
								"http": map[string]interface{}{
									"endpoint": "0.0.0.0:4318",
								},
							},
						},
						"prometheus": map[string]interface{}{
							"config": map[string]interface{}{
								"scrape_configs": []interface{}{
									map[string]interface{}{
										"job_name":        "enovaserving",
										"scrape_interval": "5s",
										"static_configs": []interface{}{
											map[string]interface{}{
												"targets": []string{fmt.Sprintf("%s-svc", w.Spec.Name) + ".emergingai.svc.cluster.local:9199"}},
										},
									},
								},
							},
						},
					},
					"exporters": map[string]interface{}{
						"kafka": map[string]interface{}{
							"brokers":          w.Spec.Collector.Kafka.Brokers,
							"topic":            "k8s-common-collector",
							"protocol_version": "2.0.0",
							"auth": map[string]interface{}{
								"sasl": map[string]interface{}{
									"mechanism": "PLAIN",
									"username":  w.Spec.Collector.Kafka.Username,
									"password":  w.Spec.Collector.Kafka.Password,
								},
							},
						},
					},
					"processors": collector.Spec.Config.Processors,
					"service":    collector.Spec.Config.Service,
				},
			},
		},
	}
}

func (w *Workload) GetOtCollectorResource() dynamic.NamespaceableResourceInterface {
	gvr := schema.GroupVersionResource{
		Group:    "opentelemetry.io",
		Version:  "v1beta1",
		Resource: "opentelemetrycollectors",
	}
	return w.K8sCli.DynamicClient.Resource(gvr)
}

func (w *Workload) isCustomized() bool {
	return !strings.HasSuffix(strings.Split(w.Spec.Image, ":")[0], "/enova")
}

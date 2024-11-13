package resource

import "github.com/Emerging-AI/ENOVA/escaler/pkg/meta"

type ClientInterface interface {
	DeployTask(spec meta.TaskSpec)
	DeleteTask(spec meta.TaskSpec)
	IsTaskExist(spec meta.TaskSpec) bool
	IsTaskRunning(spec meta.TaskSpec) bool
	GetRuntimeInfos(spec meta.TaskSpec) *meta.RuntimeInfo
}

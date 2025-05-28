package resource

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/docker/docker/client"

	"github.com/Emerging-AI/ENOVA/escaler/pkg/meta"

	"github.com/Emerging-AI/ENOVA/escaler/pkg/config"
	"github.com/Emerging-AI/ENOVA/escaler/pkg/logger"
	"github.com/Emerging-AI/ENOVA/escaler/pkg/redis"
	"github.com/Emerging-AI/ENOVA/escaler/pkg/resource/docker"
	rscutils "github.com/Emerging-AI/ENOVA/escaler/pkg/resource/utils"

	"github.com/google/uuid"
)

type TaskManager struct {
	RedisClient *redis.RedisClient
}

type GpuStatsInfo struct {
	GpuId  int
	Status string
}

type DockerResourceClient struct {
	DockerClient       *docker.DockerCli
	TaskManager        *TaskManager
	ContainerIDGpusMap map[string][]string
	LocalGpuStats      []*GpuStatsInfo
}

func (d *DockerResourceClient) DeployTask(task meta.TaskSpec) {
	logger.Infof("start local deploy, Replica: %d", task.Replica)

	containerIds, ok := d.TaskManager.getTaskContainerIds(task)
	logger.Infof("deployByDocker getTaskContainerIds taskName: %s, containerIds: %v", task.Name, containerIds)
	if !ok {
		for i := 0; i < task.Replica; i++ {
			containerID, err := d.singleDeployByDocker(&task)
			if err != nil {
				logger.Errorf("deployByDocker err: %v", err)
			} else {
				logger.Infof("sucees deploy task: %v, containerId: %s", task, containerID)
				containerIds = append(containerIds, containerID)
			}
		}
		_ = d.TaskManager.setTaskContainerIds(task, containerIds)
	} else {
		// delete all container when replica = 0
		if task.Replica == 0 {
			for _, containerId := range containerIds {
				if err := d.DockerClient.StopContainer(containerId); err != nil {
					logger.Errorf("deleteSingleServing containerId: %s, err: %v", containerId, err)
				}
				gpuStrList := d.ContainerIDGpusMap[containerId]
				for _, gpuIdStr := range gpuStrList {
					gpuId, err := strconv.Atoi(gpuIdStr)
					if err != nil {
						logger.Errorf("reset GpuStat error, containerId: %s, err: %v", containerId, err)
					}
					for _, gpuStat := range d.LocalGpuStats {
						if gpuStat.GpuId == gpuId {
							gpuStrList = append(gpuStrList, fmt.Sprintf("%d", gpuStat.GpuId))
							gpuStat.Status = "Available"
							break
						}
					}
				}
				delete(d.ContainerIDGpusMap, containerId)
			}
			_, err := d.TaskManager.RedisClient.DelList(task.Name)
			if err != nil {
				logger.Errorf("deleteTaskContainerIds err: %v", err)
			}
		} else {
			// first checkout replica number
			// second scale up or down when replica is not match
			if len(containerIds) > task.Replica {
				removeCnt := len(containerIds) - task.Replica
				logger.Infof("start to scale down task: %s, removeCnt: %d", task.Name, removeCnt)
				for i := 0; i < removeCnt; i++ {
					containerId := containerIds[i]
					if err := d.DockerClient.StopContainer(containerId); err != nil {
						logger.Errorf("deleteSingleServing containerId: %s, err: %v", containerIds[0], err)
					}
					gpuStrList := d.ContainerIDGpusMap[containerId]
					for _, gpuIdStr := range gpuStrList {
						gpuId, err := strconv.Atoi(gpuIdStr)
						if err != nil {
							logger.Errorf("reset GpuStat error, containerId: %s, err: %v", containerId, err)
						}
						for _, gpuStat := range d.LocalGpuStats {
							if gpuStat.GpuId == gpuId {
								gpuStrList = append(gpuStrList, fmt.Sprintf("%d", gpuStat.GpuId))
								gpuStat.Status = "Available"
								break
							}
						}
					}
				}
				containerIds = containerIds[removeCnt:]
				// d.TaskContainerIdMap[task.Name] = d.TaskContainerIdMap[task.Name][removeCnt:]
			} else if len(containerIds) < task.Replica {
				scaleoutCnt := task.Replica - len(containerIds)
				logger.Infof("start to scale up task: %s, scaleoutCnt: %d", task.Name, scaleoutCnt)
				for i := 0; i < scaleoutCnt; i++ {
					// Record Gpu
					containerID, err := d.singleDeployByDocker(&task)
					if err != nil {
						logger.Errorf("deployByDocker err: %v", err)
					} else {
						logger.Infof("sucees deploy task: %v, containerId: %s", task, containerID)
						containerIds = append(containerIds, containerID)
					}
				}
			}
			_ = d.TaskManager.setTaskContainerIds(task, containerIds)
		}
	}
}

func (d *DockerResourceClient) DeleteTask(spec meta.TaskSpec) {
	_ = d.DockerClient.StopContainer(spec.Name)
}

func (d *DockerResourceClient) IsTaskExist(spec meta.TaskSpec) bool {
	// TODO: Not implemented
	return false
}

// IsTaskRunning check all containers running
func (d *DockerResourceClient) IsTaskRunning(task meta.TaskSpec) bool {
	containerIds, ok := d.TaskManager.getTaskContainerIds(task)
	if !ok {
		logger.Infof("IsTaskRunning getTaskContainerIds failed")
		return false
	}
	ret := false
	for _, containerId := range containerIds {
		status, err := d.DockerClient.GetContainerStatus(containerId)
		if err != nil {
			logger.Errorf("IsTaskRunning GetContainerStatus error: %v", err)
			ret = false
		}
		logger.Infof("IsTaskRunning GetContainerStatus, taskName: %s, containerId: %s, status: %s", task.Name, containerId, status)
		if status == "running" {
			ret = true
		} else {
			ret = false
		}
	}
	return ret
}

func (d *DockerResourceClient) GetRuntimeInfos(spec meta.TaskSpec) *meta.RuntimeInfo {
	ret := &meta.RuntimeInfo{Source: meta.DockerSource}
	containerIds, ok := d.TaskManager.getTaskContainerIds(spec)
	if !ok {
		logger.Infof("GetTaskInfo getTaskContainerIds failed")
		return ret
	}
	for _, containerId := range containerIds {
		containerJson, err := d.DockerClient.GetContainerInfo(containerId)
		if err != nil {
			logger.Errorf("IsTaskRunning GetContainerStatus error: %v", err)
			continue
		}
		*ret.Containers = append(*ret.Containers, containerJson)
	}
	return ret
}

func (t *TaskManager) getTaskContainerIds(task meta.TaskSpec) ([]string, bool) {
	containerIds, err := t.RedisClient.GetList(task.Name)
	if err != nil {
		logger.Infof("GetList err: %v", err)
		return []string{}, false
	}
	return containerIds, true
}

func (d *DockerResourceClient) singleDeployByDocker(task *meta.TaskSpec) (string, error) {
	selectGpuNum := 0
	preferGpuNum := task.GetPreferGpuNum()
	gpuStrList := []string{}
	logger.Infof("before deploy d.LocalGpuStats: %v", d.LocalGpuStats)
	for _, gpuStat := range d.LocalGpuStats {
		if gpuStat.Status == "Available" {
			gpuStrList = append(gpuStrList, fmt.Sprintf("%d", gpuStat.GpuId))
			selectGpuNum += 1
			gpuStat.Status = "InUsed"
		}
		if selectGpuNum >= preferGpuNum {
			break
		}
	}
	logger.Infof("after deploy d.LocalGpuStats: %v, gpuStrList: %v", d.LocalGpuStats, gpuStrList)
	task.Gpus = strings.Join(gpuStrList, ",")
	containerID, err := d.createSingleServing(*task, fmt.Sprintf("%s-replica-%s", task.ExporterServiceName, uuid.New().String()[:4]))
	if err != nil {
		logger.Errorf("singleDeployByDocker err: %v", err)
	} else {
		d.ContainerIDGpusMap[containerID] = gpuStrList
	}
	return containerID, err
}

func (d *DockerResourceClient) createSingleServing(spec meta.TaskSpec, containerName string) (string, error) {
	cmd := rscutils.BuildCmdFromTaskSpec(spec)
	params := docker.CreateContainerParams{
		ImageName:     config.GetEConfig().Serving.Image,
		Cmd:           cmd,
		NetworkName:   config.GetEConfig().Serving.Network,
		Ports:         []int{spec.Port},
		NetworkAlias:  config.GetEConfig().Serving.NetworkAlias,
		ContainerName: containerName,
		Envs:          buildDockerEnvs(spec.Envs),
		Gpus:          spec.Gpus,
		Volumes:       buildDockerVolumes(spec.Volumes),
	}
	return d.DockerClient.CreateContainer(params)
}

func buildDockerEnvs(envs []meta.Env) []string {
	ret := []string{}
	for _, e := range envs {
		ret = append(ret, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}
	return ret
}

// TODO: will fix
func buildDockerVolumes(volumes []meta.Volume) []string {
	ret := []string{}
	for _, e := range volumes {
		ret = append(ret, fmt.Sprintf("%s:%s", e.Value, e.Path))
	}
	return ret
}

func getGpuStats() ([]*GpuStatsInfo, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=index", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	var gpus []*GpuStatsInfo

	for scanner.Scan() {
		line := scanner.Text()

		gpuId, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}

		gpus = append(gpus, &GpuStatsInfo{
			GpuId:  gpuId,
			Status: "Available",
		})
	}

	return gpus, nil
}

func (t *TaskManager) setTaskContainerIds(task meta.TaskSpec, containerIds []string) error {
	return t.RedisClient.SetList(task.Name, containerIds)
}

func NewDockerResourceClient() *DockerResourceClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	dockerCli := docker.DockerCli{
		Cli: cli,
		Ctx: context.Background(),
	}
	localGpuStats, err := getGpuStats()
	if err != nil {
		logger.Errorf("DockerResourceClient getGpuStats err: %v", err)
	}
	return &DockerResourceClient{
		DockerClient: &dockerCli,
		TaskManager: &TaskManager{
			RedisClient: redis.NewRedisClient(
				config.GetEConfig().Redis.Addr, config.GetEConfig().Redis.Password, config.GetEConfig().Redis.Db),
		},
		ContainerIDGpusMap: make(map[string][]string),
		LocalGpuStats:      localGpuStats,
	}
}

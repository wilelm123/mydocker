package command

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wilelm123/mydocker/pkg/cgroups"
	"github.com/wilelm123/mydocker/pkg/cgroups/subsystems"
	"github.com/wilelm123/mydocker/pkg/container"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func Run(tty bool, cmdArr []string, res *subsystems.ResourceConfig, containerName, volume, imageName string, envSlice []string, nw string, portmapping []string) {
	containerId := randStringBytes(10)
	if containerName == "" {
		containerName = containerId
	}

	parent, writePipe := container.NewParentProcess(tty, containerName, volume, imageName, envSlice)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}

	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	containerName, err := recordContainerInfo(parent.Process.Pid, cmdArr, containerName, containerId, volume)
	if err != nil {
		log.Errorf("Record container info error, %v", err)
		return
	}

	cgroupManager := cgroups.NewCgroupManager(containerId)
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	if nw != "" {
		network.Init()
	}
}

func sendInitCommand(cmdArr []string, writePipe *os.File) {
	command := strings.Join(cmdArr, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func recordContainerInfo(containerPID int, commandArray []string, containerName, id, volume string) (string, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")
	containerInfo := &container.ContainerInfo{
		Id:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      container.RUNNING,
		Name:        containerName,
		Volume:      volume,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Errorf("Record container info error, %v", err)
		return "", err
	}

	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Errorf("Mkdir error %s error %v", dirUrl, err)
		return "", err
	}

	fileName := dirUrl + "/" + container.ConfigName
	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		log.Errorf("Create file %s error, %v", fileName, err)
		return "", err
	}

	if _, err := file.WriteString(jsonStr); err != nil {
		log.Errorf("File write string error, %v", err)
		return "", err
	}

	return containerName, nil
}

func deleteContainerInfo(containerId string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerId)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Errorf("Remove dir %s error, %v", dirURL, err)
	}
}

func randStringBytes(n int) string {
	nums := "1234567890"
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = nums[rand.Intn(len(nums))]
	}
	return string(b)
}
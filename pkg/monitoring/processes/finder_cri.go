/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2018 Red Hat, Inc.
 *
 */

package processes

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	pb "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
	"k8s.io/kubernetes/pkg/kubelet/util"

	"github.com/fromanirh/kubevirt-metrics-collector/internal/pkg/log"
	"github.com/fromanirh/kubevirt-metrics-collector/pkg/procscanner"
)

const (
	DefaultTimeout = 10 * time.Second
	DefaultProcDir = "/proc"
)

type CRIPodFinder struct {
	ProcDir        string
	conn           *grpc.ClientConn
	client         pb.RuntimeServiceClient
	containerToPod map[string]string
	podInfos       map[string]string
	scanner        procscanner.ProcScanner
}

func NewCRIPodFinder(runtimeEndPoint string, timeout time.Duration, scanner procscanner.ProcScanner) (*CRIPodFinder, error) {
	log.Log.Infof("connecting to '%v'...", runtimeEndPoint)
	pr := &CRIPodFinder{
		ProcDir: DefaultProcDir,
		scanner: scanner,
	}

	addr, dialer, err := getAddressAndDialer(runtimeEndPoint)
	if err != nil {
		return nil, err
	}

	pr.conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(timeout), grpc.WithDialer(dialer))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	pr.client = pb.NewRuntimeServiceClient(pr.conn)
	log.Log.Infof("connected to '%v'!", runtimeEndPoint)
	return pr, nil
}

func (cpf *CRIPodFinder) FindPods() (map[string]*PodInfo, error) {
	var err error
	pods := make(PodMap)

	procs, err := cpf.scanner.Scan(cpf.ProcDir)
	if err != nil {
		log.Log.Warningf("error scanning for pods in %v: %v", cpf.ProcDir, err)
		return pods, err
	}
	err = cpf.updateCRIInfo()
	if err != nil {
		log.Log.Warningf("error updating CRI infos: %v", err)
		return pods, err
	}

	return pods.MapProcsToPods(cpf, procs)
}

func (cpf *CRIPodFinder) updateCRIInfo() error {
	var err error

	err = cpf.updateInfoContainers()
	if err != nil {
		return err
	}

	err = cpf.updateInfoPods()
	if err != nil {
		return err
	}

	return nil
}

func (cpf *CRIPodFinder) updateInfoContainers() error {
	st := &pb.ContainerStateValue{}
	st.State = pb.ContainerState_CONTAINER_RUNNING
	filter := &pb.ContainerFilter{}
	filter.State = st

	request := &pb.ListContainersRequest{
		Filter: filter,
	}

	r, err := cpf.client.ListContainers(context.Background(), request)
	if err != nil {
		return err
	}

	cpf.containerToPod = make(map[string]string)
	for _, c := range r.GetContainers() {
		cpf.containerToPod[c.Id] = c.PodSandboxId
	}

	return nil
}

func (cpf *CRIPodFinder) updateInfoPods() error {
	st := &pb.PodSandboxStateValue{}
	st.State = pb.PodSandboxState_SANDBOX_READY
	filter := &pb.PodSandboxFilter{}
	filter.State = st

	request := &pb.ListPodSandboxRequest{
		Filter: filter,
	}

	r, err := cpf.client.ListPodSandbox(context.Background(), request)

	if err != nil {
		return err
	}

	cpf.podInfos = make(map[string]string)
	for _, p := range r.GetItems() {
		if domainName, ok := p.Annotations["kubevirt.io/domain"]; ok {
			cpf.podInfos[p.Id] = domainName
		} else {
			cpf.podInfos[p.Id] = p.Metadata.Name
		}
	}

	return nil
}

func (cpf *CRIPodFinder) FindPodByPID(pid int32) (string, error) {
	containerId, cgroupStyle := FindContainerIDByCGroup(pid)
	if cgroupStyle != DockerCGroup {
		return "", fmt.Errorf("unsupported cgroup style: %v", cgroupStyle)
	}
	podId, ok := cpf.containerToPod[containerId]
	if !ok {
		return "", fmt.Errorf("no POD found for pid %v on container %v", pid, containerId)
	}
	podName, ok := cpf.podInfos[podId]
	if !ok {
		return "", fmt.Errorf("no info for pid %v on container %v on pod %v", pid, containerId, podId)
	}
	return podName, nil
}

func getAddressAndDialer(endpoint string) (string, func(addr string, timeout time.Duration) (net.Conn, error), error) {
	return util.GetAddressAndDialer(endpoint)
}

/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloud

import (
	"fmt"
	"strings"

	"k8s.io/klog/v2"
)

var _ MetadataService = &Metadata{}

const (
	ProviderIDValidLength = 6
)

// Metadata is info about the instance on which the driver is running
type Metadata struct {
	region          string
	zone            string
	cloudInstanceId string
	pvmInstanceId   string
}

// GetRegion returns region of the instance
func (m *Metadata) GetRegion() string {
	return m.region
}

// GetZone returns zone of the instance
func (m *Metadata) GetZone() string {
	return m.zone
}

// GetCloudInstanceId returns cloud instance id of the instance
func (m *Metadata) GetCloudInstanceId() string {
	return m.cloudInstanceId
}

// GetPvmInstanceId returns pvm instance id of the instance
func (m *Metadata) GetPvmInstanceId() string {
	return m.pvmInstanceId
}

// TokenizeProviderID tokenize the provider id into Metadata structure
// ProviderID format: ibmpowervs://<region>/<zone>/<service_instance_id>/<powervs_machine_id>
func TokenizeProviderID(providerID string) (*Metadata, error) {
	data := strings.Split(providerID, "/")
	errFormat := "invalid ProviderID format - %v, expected format - ibmpowervs://<region>/<zone>/<service_instance_id>/<powervs_machine_id>, err: %s"
	if len(data) != ProviderIDValidLength {
		return nil, fmt.Errorf(errFormat, providerID, "invalid length")
	}
	if data[2] == "" {
		return nil, fmt.Errorf(errFormat, providerID, "region can't be empty")
	}
	if data[3] == "" {
		return nil, fmt.Errorf(errFormat, providerID, "zone can't be empty")
	}
	if data[4] == "" {
		return nil, fmt.Errorf(errFormat, providerID, "service_instance_id can't be empty")
	}
	if data[5] == "" {
		return nil, fmt.Errorf(errFormat, providerID, "powervs_machine_id can't be empty")
	}
	return &Metadata{
		region:          data[2],
		zone:            data[3],
		cloudInstanceId: data[4],
		pvmInstanceId:   data[5],
	}, nil
}

// Get New Metadata Service
func NewMetadataService(k8sAPIClient KubernetesAPIClient) (MetadataService, error) {
	klog.Infof("retrieving instance data from kubernetes api")
	clientset, err := k8sAPIClient()
	if err != nil {
		klog.Warningf("error creating kubernetes api client: %v", err)
	} else {
		klog.Infof("kubernetes api is available")
		return KubernetesAPIInstanceInfo(clientset)
	}
	return nil, fmt.Errorf("error getting instance data from ec2 metadata or kubernetes api")
}

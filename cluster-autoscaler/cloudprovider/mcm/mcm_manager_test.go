/*
Copyright 2017 The Kubernetes Authors.

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

package mcm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	kubeletapis "k8s.io/kubernetes/pkg/kubelet/apis"
)

func TestBuildGenericLabels(t *testing.T) {
	labels := buildGenericLabels(&nodeTemplate{
		InstanceType: &instanceType{
			InstanceType: "c4.large",
			VCPU:         2,
			MemoryMb:     3840,
		},
		Region: "us-east-1",
	}, "sillyname")
	assert.Equal(t, "us-east-1", labels[apiv1.LabelZoneRegion])
	assert.Equal(t, "sillyname", labels[apiv1.LabelHostname])
	assert.Equal(t, "c4.large", labels[apiv1.LabelInstanceType])
	assert.Equal(t, cloudprovider.DefaultArch, labels[kubeletapis.LabelArch])
	assert.Equal(t, cloudprovider.DefaultOS, labels[kubeletapis.LabelOS])
}

func TestGenerationOfCorrectZoneValueFromMCLabel(t *testing.T) {
	var (
		resultingZone string
		zoneA         = "zone-a"
		zoneB         = "zone-b"
		randomKey     = "random-key"
		randomValue   = "random-value"
	)

	// Basic test to get zone value
	resultingZone = getZoneValueFromMCLabels(map[string]string{
		labelTopologyZone: zoneA,
	})
	assert.Equal(t, resultingZone, zoneA)

	// Prefer labelTopologyZone label over labelFailureDomainZone
	resultingZone = getZoneValueFromMCLabels(map[string]string{
		labelTopologyZone:      zoneA,
		labelFailureDomainZone: zoneB,
	})
	assert.Equal(t, resultingZone, zoneA)

	// Fallback to labelFailureDomainZone whne labelTopologyZone is not found
	resultingZone = getZoneValueFromMCLabels(map[string]string{
		randomKey:              randomValue,
		labelFailureDomainZone: zoneB,
	})
	assert.Equal(t, resultingZone, zoneB)

	// When neither of the labels are found
	resultingZone = getZoneValueFromMCLabels(map[string]string{
		randomKey: randomValue,
	})
	assert.Equal(t, resultingZone, "")
}

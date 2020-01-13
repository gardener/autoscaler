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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	apiv1 "k8s.io/api/core/v1"
	"github.com/gardener/autoscaler/cluster-autoscaler/cloudprovider"
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
	assert.Equal(t, "us-east-1", labels[kubeletapis.LabelZoneRegion])
	assert.Equal(t, "sillyname", labels[kubeletapis.LabelHostname])
	assert.Equal(t, "c4.large", labels[kubeletapis.LabelInstanceType])
	assert.Equal(t, cloudprovider.DefaultArch, labels[kubeletapis.LabelArch])
	assert.Equal(t, cloudprovider.DefaultOS, labels[kubeletapis.LabelOS])
}

func TestExtractLabelsFromMachineDeployment(t *testing.T) {
	tags := map[string]string{
		{
			"k8s.io/cluster-autoscaler/node-template/label/foo" : "bar",
		},
		{
			"bar" : "baz",
		},
	}

	labels := extractTaintsFromMachineDeployment(tags)

	assert.Equal(t, 1, len(labels))
	assert.Equal(t, "bar", labels["foo"])
}

func TestExtractTaintsFromMachineDeployment(t *testing.T) {
	tags := map[string]string{
		{
			"k8s.io/cluster-autoscaler/node-template/taint/dedicated" : "foo:NoSchedule",
		},
		{
			"k8s.io/cluster-autoscaler/node-template/taint/group" : "bar:NoExecute",
		},
		{
			"k8s.io/cluster-autoscaler/node-template/taint/app" : "fizz:PreferNoSchedule",
		},
		{
			"bar" : "baz",
		},
		{
			"k8s.io/cluster-autoscaler/node-template/taint/blank" :  ""
		},
		{
			"k8s.io/cluster-autoscaler/node-template/taint/nosplit" : "some_value",
		},
	}

	expectedTaints := []apiv1.Taint{
		{
			Key:    "dedicated",
			Value:  "foo",
			Effect: apiv1.TaintEffectNoSchedule,
		},
		{
			Key:    "group",
			Value:  "bar",
			Effect: apiv1.TaintEffectNoExecute,
		},
		{
			Key:    "app",
			Value:  "fizz",
			Effect: apiv1.TaintEffectPreferNoSchedule,
		},
	}

	taints := extractTaintsFromMachineDeployment(tags)
	assert.Equal(t, 3, len(taints))
	assert.Equal(t, makeTaintSet(expectedTaints), makeTaintSet(taints))
}

func makeTaintSet(taints []apiv1.Taint) map[apiv1.Taint]bool {
	set := make(map[apiv1.Taint]bool)
	for _, taint := range taints {
		set[taint] = true
	}
	return set
}
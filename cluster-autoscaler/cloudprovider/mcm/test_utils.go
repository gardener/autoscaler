/*
Copyright (c) 2017 SAP SE or an SAP affiliate company. All rights reserved.

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
	"k8s.io/autoscaler/cluster-autoscaler/cloudprovider"
	"k8s.io/client-go/tools/cache"
	"time"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	coreinformers "k8s.io/client-go/informers"
	"k8s.io/utils/pointer"

	machine_internal "github.com/gardener/machine-controller-manager/pkg/apis/machine"
	"github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	faketyped "github.com/gardener/machine-controller-manager/pkg/client/clientset/versioned/typed/machine/v1alpha1/fake"
	machineinformers "github.com/gardener/machine-controller-manager/pkg/client/informers/externalversions"
	mcmcache "github.com/gardener/machine-controller-manager/pkg/util/provider/cache"
	customfake "k8s.io/autoscaler/cluster-autoscaler/cloudprovider/mcm/fakeclient"
)

//func TestMachineControllerManagerSuite(t *testing.T) {
//	// for filtering out warning logs. Reflector short watch warning logs won't print now
//	klog.SetOutput(io.Discard)
//	flags := &flag.FlagSet{}
//	klog.InitFlags(flags)
//	if err := flags.Set("logtostderr", "false"); err != nil {
//		t.Errorf("failed to set flags: %v", err)
//	}
//
//	RegisterFailHandler(Fail)
//	RunSpecs(t, "Machine Controller Manager Suite")
//}

var (
	TestMachineClass = "machineClass-0"
	testNamespace    = "test-namespace"
)

func newMachineDeployment(
	specTemplate *v1alpha1.MachineTemplateSpec,
	replicas int32,
	minReadySeconds int32,
	maxSurge int,
	maxUnavailable int,
	statusTemplate *v1alpha1.MachineDeploymentStatus,
	owner *metav1.OwnerReference,
	annotations map[string]string,
	labels map[string]string,
) *v1alpha1.MachineDeployment {
	md := newMachineDeployments(1, specTemplate, replicas, minReadySeconds, statusTemplate, owner, annotations, labels)[0]
	intStrMaxSurge := intstr.FromInt(maxSurge)
	intStrMaxUnavailable := intstr.FromInt(maxUnavailable)
	md.Spec.Strategy.RollingUpdate.MaxSurge = &intStrMaxSurge
	md.Spec.Strategy.RollingUpdate.MaxUnavailable = &intStrMaxUnavailable

	return md
}

func newMachineDeployments(
	machineDeploymentCount int,
	specTemplate *v1alpha1.MachineTemplateSpec,
	replicas int32,
	minReadySeconds int32,
	statusTemplate *v1alpha1.MachineDeploymentStatus,
	owner *metav1.OwnerReference,
	annotations map[string]string,
	labels map[string]string,
) []*v1alpha1.MachineDeployment {

	intStr1 := intstr.FromInt(1)
	machineDeployments := make([]*v1alpha1.MachineDeployment, machineDeploymentCount)
	for i := range machineDeployments {
		machineDeployment := &v1alpha1.MachineDeployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.sapcloud.io",
				Kind:       "MachineDeployment",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("machinedeployment-%d", i),
				Namespace: testNamespace,
				Labels:    labels,
			},
			Spec: v1alpha1.MachineDeploymentSpec{
				MinReadySeconds: minReadySeconds,
				Replicas:        replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: deepCopy(specTemplate.ObjectMeta.Labels),
				},
				Strategy: v1alpha1.MachineDeploymentStrategy{
					Type: v1alpha1.RollingUpdateMachineDeploymentStrategyType,
					RollingUpdate: &v1alpha1.RollingUpdateMachineDeployment{
						MaxSurge:       &intStr1,
						MaxUnavailable: &intStr1,
					},
				},
				Template: *specTemplate.DeepCopy(),
			},
		}

		if statusTemplate != nil {
			machineDeployment.Status = *statusTemplate.DeepCopy()
		}

		if owner != nil {
			machineDeployment.OwnerReferences = append(machineDeployment.OwnerReferences, *owner.DeepCopy())
		}

		if annotations != nil {
			machineDeployment.Annotations = annotations
		}

		machineDeployments[i] = machineDeployment
	}
	return machineDeployments
}

func newMachineSet(
	specTemplate *v1alpha1.MachineTemplateSpec,
	name string,
	replicas int32,
	minReadySeconds int32,
	statusTemplate *v1alpha1.MachineSetStatus,
	owner *metav1.OwnerReference,
	annotations map[string]string,
	labels map[string]string,
) *v1alpha1.MachineSet {
	ms := newMachineSets(1, specTemplate, replicas, minReadySeconds, statusTemplate, owner, annotations, labels)[0]
	ms.Name = name
	return ms
}

func newMachineSets(
	machineSetCount int,
	specTemplate *v1alpha1.MachineTemplateSpec,
	replicas int32,
	minReadySeconds int32,
	statusTemplate *v1alpha1.MachineSetStatus,
	owner *metav1.OwnerReference,
	annotations map[string]string,
	labels map[string]string,
) []*v1alpha1.MachineSet {

	machineSets := make([]*v1alpha1.MachineSet, machineSetCount)
	for i := range machineSets {
		ms := &v1alpha1.MachineSet{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.sapcloud.io",
				Kind:       "MachineSet",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("machineset-%d", i),
				Namespace: testNamespace,
				Labels:    labels,
			},
			Spec: v1alpha1.MachineSetSpec{
				MachineClass:    *specTemplate.Spec.Class.DeepCopy(),
				MinReadySeconds: minReadySeconds,
				Replicas:        replicas,
				Selector: &metav1.LabelSelector{
					MatchLabels: deepCopy(specTemplate.ObjectMeta.Labels),
				},
				Template: *specTemplate.DeepCopy(),
			},
		}

		if statusTemplate != nil {
			ms.Status = *statusTemplate.DeepCopy()
		}

		if owner != nil {
			ms.OwnerReferences = append(ms.OwnerReferences, *owner.DeepCopy())
		}

		if annotations != nil {
			ms.Annotations = annotations
		}

		machineSets[i] = ms
	}
	return machineSets
}

func deepCopy(m map[string]string) map[string]string {
	r := make(map[string]string, len(m))
	for k := range m {
		r[k] = m[k]
	}
	return r
}

func newMachinesFromMachineSet(
	machineCount int,
	machineSet *v1alpha1.MachineSet,
	statusTemplate *v1alpha1.MachineStatus,
	annotations map[string]string,
	labels map[string]string,
) []*v1alpha1.Machine {
	t := &machineSet.TypeMeta

	finalLabels := make(map[string]string, 0)
	for k, v := range labels {
		finalLabels[k] = v
	}
	for k, v := range machineSet.Spec.Template.Labels {
		finalLabels[k] = v
	}

	return newMachines(
		machineCount,
		&machineSet.Spec.Template,
		statusTemplate,
		&metav1.OwnerReference{
			APIVersion:         t.APIVersion,
			Kind:               t.Kind,
			Name:               machineSet.Name,
			UID:                machineSet.UID,
			BlockOwnerDeletion: pointer.BoolPtr(true),
			Controller:         pointer.BoolPtr(true),
		},
		annotations,
		finalLabels,
	)
}

func newMachine(
	specTemplate *v1alpha1.MachineTemplateSpec,
	statusTemplate *v1alpha1.MachineStatus,
	owner *metav1.OwnerReference,
	annotations map[string]string,
	labels map[string]string,
) *v1alpha1.Machine {
	return newMachines(1, specTemplate, statusTemplate, owner, annotations, labels)[0]
}

func newMachines(
	machineCount int,
	specTemplate *v1alpha1.MachineTemplateSpec,
	statusTemplate *v1alpha1.MachineStatus,
	owner *metav1.OwnerReference,
	annotations map[string]string,
	labels map[string]string,
) []*v1alpha1.Machine {
	machines := make([]*v1alpha1.Machine, machineCount)

	if annotations == nil {
		annotations = make(map[string]string, 0)
	}
	if labels == nil {
		labels = make(map[string]string, 0)
	}

	currentTime := metav1.Now()

	for i := range machines {
		m := &v1alpha1.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.sapcloud.io",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:              fmt.Sprintf("machine-%d", i),
				Namespace:         testNamespace,
				Labels:            labels,
				Annotations:       annotations,
				CreationTimestamp: metav1.Now(),
				DeletionTimestamp: &currentTime,
			},
			Spec: *newMachineSpec(&specTemplate.Spec, i),
		}

		if statusTemplate != nil {
			m.Status = *newMachineStatus(statusTemplate)
		}

		if owner != nil {
			m.OwnerReferences = append(m.OwnerReferences, *owner.DeepCopy())
		}

		machines[i] = m
	}
	return machines
}

func newNode(
	_ int,
	nodeSpec *corev1.NodeSpec,
	nodeStatus *corev1.NodeStatus,
) *corev1.Node {
	return newNodes(1, nodeSpec, nodeStatus)[0]
}

func newNodes(
	nodeCount int,
	nodeSpec *corev1.NodeSpec,
	_ *corev1.NodeStatus,
) []*corev1.Node {

	nodes := make([]*corev1.Node, nodeCount)
	for i := range nodes {
		node := &corev1.Node{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "v1",
				Kind:       "Node",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("node-%d", i),
			},
			Spec: *nodeSpec.DeepCopy(),
		}

		nodes[i] = node
	}
	return nodes
}

func newObjectMeta(meta *metav1.ObjectMeta, index int) *metav1.ObjectMeta {
	r := meta.DeepCopy()

	if r.Name != "" {
		return r
	}

	if r.GenerateName != "" {
		r.Name = fmt.Sprintf("%s-%d", r.GenerateName, index)
		return r
	}

	r.Name = fmt.Sprintf("machine-%d", index)
	return r
}

func newMachineSpec(specTemplate *v1alpha1.MachineSpec, index int) *v1alpha1.MachineSpec {
	r := specTemplate.DeepCopy()
	if r.ProviderID == "" {
		return r
	}

	r.ProviderID = fmt.Sprintf("%s-%d", r.ProviderID, index)
	return r
}

func newMachineStatus(statusTemplate *v1alpha1.MachineStatus) *v1alpha1.MachineStatus {
	if statusTemplate == nil {
		return &v1alpha1.MachineStatus{}
	}

	return statusTemplate.DeepCopy()
}

func createMcmManager(
	stop <-chan struct{},
	namespace string,
	nodeGroups []string, controlMachineObjects, targetCoreObjects []runtime.Object,
) (*McmManager, *customfake.FakeObjectTrackers, []cache.InformerSynced) {

	fakeControlMachineClient, controlMachineObjectTracker := customfake.NewMachineClientSet(controlMachineObjects...)
	fakeTypedMachineClient := &faketyped.FakeMachineV1alpha1{
		Fake: &fakeControlMachineClient.Fake,
	}
	fakeTargetCoreClient, targetCoreObjectTracker := customfake.NewCoreClientSet(targetCoreObjects...)
	fakeObjectTrackers := customfake.NewFakeObjectTrackers(
		controlMachineObjectTracker,
		targetCoreObjectTracker,
	)
	fakeObjectTrackers.Start()

	coreTargetInformerFactory := coreinformers.NewFilteredSharedInformerFactory(
		fakeTargetCoreClient,
		100*time.Millisecond,
		namespace,
		nil,
	)
	defer coreTargetInformerFactory.Start(stop)
	coreTargetSharedInformers := coreTargetInformerFactory.Core().V1()
	nodes := coreTargetSharedInformers.Nodes()

	controlMachineInformerFactory := machineinformers.NewFilteredSharedInformerFactory(
		fakeControlMachineClient,
		100*time.Millisecond,
		namespace,
		nil,
	)
	defer controlMachineInformerFactory.Start(stop)

	machineSharedInformers := controlMachineInformerFactory.Machine().V1alpha1()
	machines := machineSharedInformers.Machines()
	machineSets := machineSharedInformers.MachineSets()
	machineDeployments := machineSharedInformers.MachineDeployments()
	machineClasses := machineSharedInformers.MachineClasses()

	internalExternalScheme := runtime.NewScheme()
	Expect(machine_internal.AddToScheme(internalExternalScheme)).To(Succeed())
	Expect(v1alpha1.AddToScheme(internalExternalScheme)).To(Succeed())

	mcmManager := McmManager{
		namespace: namespace,
		interrupt: make(chan struct{}),
		discoveryOpts: cloudprovider.NodeGroupDiscoveryOptions{
			NodeGroupSpecs: nodeGroups,
		},
		machineClient:           fakeTypedMachineClient,
		machineDeploymentLister: machineDeployments.Lister(),
		machineSetLister:        machineSets.Lister(),
		machineLister:           machines.Lister(),
		machineClassLister:      machineClasses.Lister(),
		nodeLister:              nodes.Lister(),
	}

	hasSyncedCachesFns := []cache.InformerSynced{
		nodes.Informer().HasSynced,
		machines.Informer().HasSynced,
		machineSets.Informer().HasSynced,
		machineDeployments.Informer().HasSynced,
		machineClasses.Informer().HasSynced,
	}

	return &mcmManager, fakeObjectTrackers, hasSyncedCachesFns
}

func waitForCacheSync(stop <-chan struct{}, hasSyncedCachesFns []cache.InformerSynced) {
	Expect(mcmcache.WaitForCacheSync(stop, hasSyncedCachesFns...)).To(BeTrue())
}

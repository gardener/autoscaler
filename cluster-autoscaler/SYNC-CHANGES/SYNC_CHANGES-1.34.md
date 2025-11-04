<!--- For help refer to https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md?plain=1 as example --->

- [v1.34.0](#v1340)
    - [Synced with which upstream CA](#synced-with-which-upstream-ca)
    - [Changes made](#changes-made)
        - [Changes during go.mod update](#changes-during-gomod-update)
        - [Others](#others)

# v1.34.0

## Synced with which upstream CA
[v1.34.0](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.34.0/cluster-autoscaler)

## Changes made
- See general release notes of 1.34.0: https://github.com/kubernetes/autoscaler/releases/tag/cluster-autoscaler-1.34.0
- Binpacking simulator will now consider old nodes when trying to pack pods with topology spread constraints in order to avoid unnecessary scale ups.
- Default value for `--cordon-node-before-terminating` changed to `True`.
- Introducing `--force-delete-failed-nodes` to force delete nodes with errors. This prevents CA from entering in error loop if scale up fails when node group is below its min size.
- New field `KeepPartiallyFailedZeroOrMaxScalingNodeGroups` added to `NodeGroupAutoScalingOptions` to prevent nodes from the node-groups being deleted if they have creation errors i.e. the node-group is left unmodified unless all their nodes fail.
- Extends Helm charts with support for configuring `dnsConfig`.
- New CRD `CapacityBuffer` added which defines spare capacity per workload or set workloads. This is intended to allow users to express the need for spare capacity in the cluster. Additional details for the changes can be found [here](https://github.com/kubernetes/autoscaler/blob/cluster-autoscaler-1.34.0/cluster-autoscaler/proposals/buffers.md)
- Deprecated `ProvisioningRequest v1beta1`.

### Changes during go.mod update
- mcm v0.59.0 -> v0.60.0
- mcm-provider-aws v0.25.0 -> v0.26.0
- mcm-provider-azure v0.16.0 -> v0.17.0
- k8s.io/api v0.33.0 -> v0.34.1

### Others
- [Release matrix](../README.md#releases-gardenerautoscaler) of Gardener Autoscaler updated.
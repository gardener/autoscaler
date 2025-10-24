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
- Binpacking simulator will now consider old nodes when trying to pack pods with topology spread contraints in order to avoid unnecessary scale ups.
- Default value for `--cordon-node-before-terminating` changed to `True`.
- The KeepPartiallyFailedZeroOrMaxScalingNodeGroups field has been added to NodeGroupAutoScalingOptions to not allow nodes from the node-groups to be deleted if any node has creation errors i.e. the node-group is left unmodified unless all their nodes fail.
- Adds Helm chart support for configuring dnsConfig. When dnsConfig is provided, it gets set for the Deployment Pod template spec. Also, backward compatibility is preserved; there are no template changes with default values.
- In order to allow users to express the need for spare capacity in the cluster a new kubernetes object will be introduced called a CapacityBuffer, which will define spare capacity per workload or set workloads. Configuration would be translated to pod specs that could be injected in memory by autoscaler to drive scaling decisions for the cluster.
- Added CapacityBuffer controller loop along with the main needed skeleton for buffers with podTemplateRef.
    - Filters: CapacityBuffers provisioning strategy and status filtering.
    - Translators: podTemplateRef translator that updates buffer status accordingly.
    - Updater: updates buffer status via capacity buffer client.
    - Controller: initiates the needed components and contains the reconciliation loops.
- Changed the internal name for the annotation config to fromNodeAnnotationKey to match the same format as fromNodeLabelKey.
- Deprecated ProvisioningRequest v1beta1.

### Changes during go.mod update
- mcm v0.59.0 -> v0.60.0
- mcm-provider-aws v0.25.0 -> v0.26.0
- mcm-provider-azure v0.16.0 -> v0.17.0
- k8s.io/api v0.33.0 -> v0.34.1

### Others
- [Release matrix](../README.md#releases-gardenerautoscaler) of Gardener Autoscaler updated.
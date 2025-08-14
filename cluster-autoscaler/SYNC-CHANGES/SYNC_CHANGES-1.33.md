<!--- For help refer to https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md?plain=1 as example --->

- [v1.33.0](#v1330)
    - [Synced with which upstream CA](#synced-with-which-upstream-ca)
    - [Changes made](#changes-made)
        - [Changes during go.mod update](#changes-during-gomod-update)
        - [Others](#others)

# v1.33.0

## Synced with which upstream CA
[v1.33.0](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.33.0/cluster-autoscaler)

## Changes made
- See general release notes of 1.33.0: https://github.com/kubernetes/autoscaler/releases/tag/cluster-autoscaler-1.33.0
- The default cluster-autoscaler expander changed from `random` to `least-waste` to prevent misused expensive nodes.
- DRA Structured Parameters feature set now production-ready. Experimental version shipped with v1.32.0.
- Added checkCapacityProvisioningRequestProcessorInstance option that allows to filter check capacity ProvisioningRequests to process by specific Cluster Autoscaler instance.
- Adding pods when creating a cluster snapshot for scale up is now parallelized, improving performance. Parallelization can be adjusted or disabled by setting the cluster-snapshot-parallelization flag.
- Improved frequent loops when there is only one of scale up activity or ProvisioningRequest processing is productive.
- ProvisioningRequests can now be processed without delay if they have just been created.
- Deprecated flags removed:
    - `max-autoprovisioned-node-group-count`
    - `node-autoprovisioning-enabled`
    - `max-empty-bulk-delete`
- Deployments with topology spread constraints will no longer block scale down if the removable node changes the global minimum
- Updated pipeline to target `go v1.24.0`
- New `ScaleDownNoCandidates` status emitted instead of existing `ScaleDownInCooldown` when there are no candidates.
  `last_activity{activity=scaleDown}` metric will be updated even when there are no candidates.

### Changes during go.mod update
- mcm v0.57.0 -> v0.59.0
- mcm-provider-aws v0.23.0 -> v0.25.0
- mcm-provider-azure v0.15.1 -> v0.16.0
- k8s.io/api v0.32.0 -> v0.33.3
- github.com/onsi/ginkgo/v2 v2.21.0 -> v2.23.0

### Others
- [Release matrix](../README.md#releases-gardenerautoscaler) of Gardener Autoscaler updated.
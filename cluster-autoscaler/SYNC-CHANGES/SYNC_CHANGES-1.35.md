<!--- For help refer to https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md?plain=1 as example --->

- [v1.35.0](#v1350)
    - [Synced with which upstream CA](#synced-with-which-upstream-ca)
    - [Changes made](#changes-made)
        - [Changes during go.mod update](#changes-during-gomod-update)
        - [Others](#others)

# v1.35.0

## Synced with which upstream CA
[v1.35.0](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.35.0/cluster-autoscaler)

## Changes made
- See general release notes of 1.35.0: https://github.com/kubernetes/autoscaler/releases/tag/cluster-autoscaler-1.35.0
- Azure cloud provider SDK updated to v2.
- CapacityBuffers graduated to v1beta1 and integrated with resource quotas.
- `--scale-down-enabled` flag deprecated; use `--scale-down-delay-after-add=0s` instead.
- Pod listing now filters out unschedulable and scheduler-unprocessed pods not in the allowed schedulers list.
- Introduced `TemplateNodeInfoRegistry` for caching and exposing template NodeInfos; DRA processor uses cached templates.
- CSI node awareness enabled in scale-up path.
- AWS: updated to `k8s.io/cloud-provider-aws v1.35.1`; migrated from `aws-sdk-go` v1 to v2; added g7e EC2 instance types.
- Log stream handling corrected for stderr/file-only/file+stderr modes.
- Latency tracker: prevent metric flapping and negative latency values.
- `MixedTemplateNodeInfoProvider` now handles `NodeGroupForNode` errors gracefully.
- Go version bumped to 1.25.0 (required by k8s v1.35).

### Changes during go.mod update
- mcm v0.60.0 -> v0.61.3
- mcm-provider-aws v0.26.0 -> v0.27.3
- mcm-provider-azure v0.17.0 -> v0.19.0
- k8s.io/api v0.34.1 -> v0.35.0
- k8s.io/kubernetes v1.34.1 -> v1.35.0

### Others
- [Release matrix](../README.md#releases-gardenerautoscaler) of Gardener Autoscaler updated.

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
- CA now supports overriding `max-node-startup-time` per node group via `NodeGroupAutoscalingOptions`.
- New flag `--max-node-startup-time` (default 15m): the maximum time from the moment the node is registered to the time the node is ready. Can be overridden per node group.
- New flag `--scale-from-unschedulable` (default false): when set, CA ignores a node's `.spec.unschedulable` field in node templates when considering whether to scale a node group.
- CapacityBuffers graduated to v1beta1 and integrated with resource quotas.
- `--scale-down-enabled` flag deprecated; use `--scale-down-delay-after-add=0s` instead.
- Pod listing now filters out unschedulable and scheduler-unprocessed pods not in the allowed schedulers list.
- Introduced `TemplateNodeInfoRegistry` for caching and exposing template NodeInfos; DRA processor uses cached templates.
- CSI node awareness enabled in scale-up path (controlled by `--enable-csi-node-aware-scheduling` flag).
- Log stream handling corrected for stderr/file-only/file+stderr modes.
- Latency tracker: prevent metric flapping and negative latency values.
- `MixedTemplateNodeInfoProvider` now handles `NodeGroupForNode` errors gracefully.
- Go version bumped to 1.26.2 (required by mcm v0.62.0).
- `WatchListClient` feature gate set to `false` in MCM test files to prevent fake client panics (k8s 1.35 enables `WatchListClient` by default, which is incompatible with the fake clientsets used in tests).

### Changes during go.mod update
- mcm v0.60.0 -> v0.62.0
- mcm-provider-aws v0.26.0 -> v0.27.3
- mcm-provider-azure v0.17.0 -> v0.19.0
- k8s.io/api v0.34.1 -> v0.35.0
- k8s.io/kubernetes v1.34.1 -> v1.35.0

### Others
- [Release matrix](../README.md#releases-gardenerautoscaler) of Gardener Autoscaler updated.

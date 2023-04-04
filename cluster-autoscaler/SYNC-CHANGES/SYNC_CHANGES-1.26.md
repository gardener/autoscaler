<!--- For help refer to https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md?plain=1 as example --->

- [v1.26.0](#v1260)
  - [Synced with which upstream CA](#synced-with-which-upstream-ca)
  - [Changes made](#changes-made)
    - [During merging](#during-merging)
    - [During vendoring k8s](#during-vendoring-k8s)
    - [Others](#others)
- [v1.26.1](#v1261)
  - [Synced with which upstream CA](#synced-with-which-upstream-ca-1)
  - [Changes made](#changes-made-1)
    - [During merging](#during-merging-1)
    - [During vendoring k8s](#during-vendoring-k8s-1)


# v1.26.0

## Synced with which upstream CA

[v1.26.0](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.26.0/cluster-autoscaler)

## Changes made
 - See general release notes of 1.26.0: https://github.com/kubernetes/autoscaler/releases/tag/cluster-autoscaler-1.26.0

### During merging
  - `cloudprovider/cloud_provider.go` has new method in interface: `HasInstance`. Introduced implementation of this method
    in `cloudprovider/mcm/mcm_cloud_provider.go`

### During vendoring k8s
- mcm v0.47.0 -> 0.48.1
- mcm-provider-aws v0.15.0 -> v0.16.0
- mcm-provider-azure v0.9.0 -> v0.10.0

### Others
- [Release matrix](../README.md#releases-gardenerautoscaler) of Gardener Autoscaler updated.


# v1.26.1

## Synced with which upstream CA

[v1.26.1](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.26.1/cluster-autoscaler)

## Changes made
 - See general release notes of 1.26.1: https://github.com/kubernetes/autoscaler/releases/tag/cluster-autoscaler-1.26.1

### During merging

### During vendoring k8s
- mcm v0.48.1 -> 0.48.2
- mcm-provider-aws v0.16.0 -> v0.17.0
- mcm-provider-azure remains at v0.10.0

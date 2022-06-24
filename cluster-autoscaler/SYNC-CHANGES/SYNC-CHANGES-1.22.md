<!--- For help refer to https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md?plain=1 as example --->

- [v1.22.0](#v1220)
    - [Synced with which upstream CA](#synced-with-which-upstream-ca)
    - [Changes made](#changes-made)
        - [To FAQ](#to-faq)
        - [During merging](#during-merging)
        - [During vendoring k8s](#during-vendoring-k8s)
        - [Others](#others)

- [v1.22.1](#v1221)
    - [Synced with which upstream CA](#synced-with-which-upstream-ca-1)
    - [Changes made](#changes-made-1)
        - [To FAQ](#to-faq-1)
        - [During merging](#during-merging-1)
        - [During vendoring k8s](#during-vendoring-k8s-1)
        - [Others](#others-1)


# v1.22.0


## Synced with which upstream CA

[v1.22.0](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.22.0/cluster-autoscaler)

## Changes made

### To FAQ

- ignored integration test package files during running unit tests.

### During merging
- // go:build tags were removed from `builder_all.go`
- `Architecture` field added to InstanceType struct. This field is used to set label `Kubernetes.io/arch` label during forming node-template


### During vendoring k8s
- mcm 0.42 -> mcm 0.44.2
- mcm-provider-aws 0.9.0 -> 0.10.0
- mcm-provider-azure 0.6.0 -> 0.6.1

### Others
- [Release matrix](../README.md#releases-gardenerautoscaler) of Gardener Autoscaler updated.
- GO111MODULE=off in .ci/build and .ci/test.


# v1.22.1


## Synced with which upstream CA

[v1.22.3](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.22.3/cluster-autoscaler)

## Changes made

### To FAQ

- updated with new question and answers

### During merging
_None_

### During vendoring k8s
- didn't vendor in k8s 1.25 which is vendored in upstream CA 1.22.3

### Others


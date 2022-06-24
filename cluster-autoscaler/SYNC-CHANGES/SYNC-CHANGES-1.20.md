<!--- For help refer to https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md?plain=1 as example --->

- [v1.20.0](#v1200)
    - [Synced with which upstream CA](#synced-with-which-upstream-ca)
    - [Changes made](#changes-made)
        - [To FAQ](#to-faq)
        - [During merging](#during-merging)
        - [During vendoring k8s](#during-vendoring-k8s)
        - [Others](#others)
- [v1.20.1](#v1201)
    - [Synced with which upstream CA](#synced-with-which-upstream-ca-1)
    - [Changes made](#changes-made-1)
        - [To FAQ](#to-faq-1)
        - [During merging](#during-merging-1)
        - [During vendoring k8s](#during-vendoring-k8s-1)
        - [Others](#others-1)


# v1.20.0


## Synced with which upstream CA

[v1.20.0](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.20.0/cluster-autoscaler)

## Changes made

### To FAQ

- Updated steps for `How to sync with upstream autoscaler`
    - broke steps into two parts
        - syncing with upstream CA < v1.21.0
        - sycing with upstream CA >= v1.21.0
    - how to use `update-vendor.sh` updated
- Warning enhanded for `How to vendor new MCM version`
### During merging
_None_
### During vendoring k8s
- Used the old `update_vendor.sh` to vendor.
### Others
- Updated README.md for cluster-autoscaler repo to contain new [release matrix](../README.md#releases-gardenerautoscaler) of Gardener Autoscaler

# v1.20.1


## Synced with which upstream CA

[v1.20.3](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.20.3/cluster-autoscaler)

## Changes made

### To FAQ
_None_
### During merging
- didn't update `cluster-autoscaler/cloudprovider/aws/ec2_instance_types.go` as its not used anymore by mcm, the one used
now (in worst case) is `cluster-autoscaler/cloudprovider/mcm/ec2_instance_types.go`
### During vendoring k8s
_None_
### Others
_None_

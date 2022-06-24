<!--- For help refer to https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md?plain=1 as example --->

- [v1.21.0](#v1210)
    - [Synced with which upstream CA](#synced-with-which-upstream-ca)
    - [Changes made](#changes-made)
        - [To FAQ](#to-faq)
        - [During merging](#during-merging)
        - [During vendoring k8s](#during-vendoring-k8s)
        - [Others](#others)
- [v1.21.1](#v1211)
    - [Synced with which upstream CA](#synced-with-which-upstream-ca-1)
    - [Changes made](#changes-made-1)
        - [To FAQ](#to-faq-1)
        - [During merging](#during-merging-1)
        - [During vendoring k8s](#during-vendoring-k8s-1)
        - [Others](#others-1)
        

# v1.21.0


## Synced with which upstream CA

[v1.21.0](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.21.0/cluster-autoscaler)

## Changes made

### To FAQ

- ignored integration test package files during running unit tests.
### During merging

- Removed `clusterapi.ProviderName` in builder_all.go
### During vendoring k8s
- Used the new `update_vendor.sh` to vendor.
- “k8s.io/kubernetes/pkg/kubelet/apis" => “k8s.io/kubelet/pkg/apis"
- mcm 0.42 -> mcm 0.44.1
- mcm-provider-aws 0.8.0 -> 0.9.0
- mcm-provider-azure 0.5.0 -> 0.6.0
- google.golang.org/grpc -> v1.29.0
- nodegroup interface GetOptions() trivial implementation done
### Others
- [Release matrix](../README.md#releases-gardenerautoscaler) of Gardener Autoscaler updated.
- skipped following directories for boilerplate checking in `hack/boilerplate/boilerplate.py`
    - cluster-autoscaler/cloudprovider/builder/builder_all.go
    - cluster-autoscaler/cloudprovider/mcm
    - cluster-autoscaler/integration


# v1.21.1


## Synced with which upstream CA

[v1.21.3](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.21.3/cluster-autoscaler)

## Changes made

### To FAQ

- included new questions and answers
- included new steps to vendor new MCM version

### During merging

- included new ec2 instances in `cluster-autoscaler/cloudprovider/aws/ec2_instance_types.go`

### During vendoring k8s
Still vendoring k8s 1.21.0 in this fork but the upstream 1.21.3 is vendoring k8s 1.25.0

### Others
_None_
<!--- For help refer to https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG/CHANGELOG-1.20.md?plain=1 as example --->

- [v1.36.0](#v1360)
    - [Synced with which upstream CA](#synced-with-which-upstream-ca)
    - [Changes made](#changes-made)
        - [Changes during go.mod update](#changes-during-gomod-update)
        - [Others](#others)

# v1.36.0

## Synced with which upstream CA
[v1.36.0](https://github.com/kubernetes/autoscaler/tree/cluster-autoscaler-1.36.0/cluster-autoscaler)

## Changes made
- See general release notes of 1.36.0: https://github.com/kubernetes/autoscaler/releases/tag/cluster-autoscaler-1.36.0
- Notable upstream features in 1.36.0:
    - New `CapacityQuota` CRD and reconciliation loop (`--capacity-quotas-enabled`) for limiting resources scaled up by CA; resource-quota tracking refactored to support both minimum and maximum limits (used by the scale-down-prevention-on-quota feature).
    - New experimental `--salvo-scale-up` flag (with `--salvo-scale-up-budget`) allowing multiple scale-ups in a single CA loop.
    - New `--startup-taint-prefix` flag to treat taints with matching key prefixes as startup taints.
    - New `--max-startup-time` flag for a separate liveness-probe timeout until the first successful loop.
    - Fastpath binpacking optimization for scale-ups (`FastpathBinpackingEnabled`).
    - DRA: partitionable devices support, DRA scale-up/scale-down metrics, and a `dra_node_template_resources_mismatch` metric.
    - Azure SDK updated to v2; `AtomicIncreaseSize` for VMSS/VMS; optional ETag optimistic concurrency on VMSS capacity updates.
    - `cluster-autoscaler-status` config map gains a `suspended` state.
    - New provider registration layout: upstream introduced `cloudprovider/router/` alongside the existing `cloudprovider/builder/`. The fork's MCM provider continues to live in `cloudprovider/builder/builder_mcm.go`.
    - Versioning is now derived at build time via `-ldflags` (the Makefile computes `VERSION` from the git tag); `version.ClusterAutoscalerVersion` defaults to `dev` and is overridden during the build.

### Fork-specific changes preserved during the sync
The following Gardener fork patches (marked `FORK-CHANGE`) were re-applied on top of upstream 1.36.0:
- **MCM provider registration adapted to upstream 1.36's self-registration pattern.** Upstream 1.36 replaced the compile-time `buildCloudProvider` switch in `cloudprovider/builder/cloud_provider_builder.go` with a dynamic registry (`RegisterCloudProvider` / `GetCloudProviderBuilder`). The fork's MCM provider was adapted to match: `cloudprovider/mcm/mcm_cloud_provider.go` now registers itself via `init()` calling `builder.RegisterCloudProvider(ProviderName, ...)` and `builder.SetDefaultCloudProvider(ProviderName)`. A new `cloudprovider/router/router_mcm.go` (build tag: `mcm`) blank-imports the mcm package to trigger the `init()`. A minimal `cloudprovider/builder/builder_mcm.go` stub (build tag: `mcm`) is retained so that `SUPPORTED_BUILD_TAGS` detection continues to work. `make start` and `.ci/build` updated to pass `-tags mcm`.
- `core/static_autoscaler.go`: `fixNodeGroupSize` remains commented out (prevents removal of "registered but long not ready" nodes, which caused scaling below minimum size and removal of ready nodes during meltdown); scale-down status log kept at V(2); force-delete logging retained (upstream extracted the scale-down block into a new `scaleDown` method — the fork changes were relocated accordingly).
- `processors/nodeinfosprovider/mixed_nodeinfos_processor.go` (+ test): node-info cache is consulted only when the node group's min size is not zero (avoids using a stale cached template after an instance-type change).
- `config/autoscaling_options.go`: default node-group difference ratios kept at `0.5` (gardener node groups share similar labels; balancing during scale-from-zero).
- `processors/nodegroupset/compare_nodegroups.go`: Gardener-specific label constants and ignored-label handling retained; adopted upstream's `node.GetRequested()` API change.
- `core/scaledown/eligibility/eligibility.go`, `core/scaleup/orchestrator/orchestrator.go`, `simulator/scheduling/hinting_simulator.go`, `simulator/utilization/info.go`: fork debug/diagnostic logging and log-level bumps retained.
- `core/scaledown/resource/limits.go`: **removed** — upstream deleted the legacy `LimitsFinder` (the fork's Go 1.26 `go vet` fix was on now-deleted code).
- Go 1.26 `go vet` format fixes retained (`simulator/clustersnapshot/scheduling_error.go`, `provisioningrequest/provreqclient/client.go`).
- `#nosec` SAST suppressions retained (`cloudprovider/utho/*`, `capacitybuffer/controller/*`).
- `main.go`: logs the Gardener autoscaler version (from the `VERSION` file) in addition to the upstream CA version.
- Fork build/CI/doc files preserved: `Dockerfile` (gardener `.ci/build` flow), `Makefile` (fork `builder/` build-tag scan, release-validate, `--test.short`; upstream `test-ci`/`test-controllers`/envtest tooling adopted), `.github/workflows/ci.yaml`/`release.yaml`/`ca-test.yaml`, `.gitignore`, `hack/for-go-proj.sh`, `hack/boilerplate/boilerplate.py`, `FAQ.md` (sync runbook), `README.md`.

### Changes during go.mod update
- mcm v0.62.0 -> v0.62.1
- mcm-provider-aws v0.27.3 -> v0.28.1
- mcm-provider-azure v0.19.0 -> v0.20.1
- k8s.io/* and k8s.io/kubernetes updated to the 1.36.0 line via `./hack/update-deps.sh 1.36.0 1.36.0` (k8s.io/api -> v0.36.2, k8s.io/kubernetes -> v1.36.0).

### Others
- The merge base with upstream was stale (~1.26-era) because prior syncs were squash-merged, so the 1.36 merge surfaced a large number of conflicts that were mostly artifacts (fork side identical to the last-synced upstream). Genuine fork patches were identified by triaging each conflicted file against the last-synced upstream release and hand-merged; the rest took upstream.
- [Release matrix](../README.md#releases-gardenerautoscaler) of Gardener Autoscaler updated.

---
name: Upstream Autoscaler Sync 
about: Sync with Upstream Autoscaler
labels: kind/enhancement
title: Issue for - Sync with Upstream v1.2x
---

## What would you like to be added
Gardener autoscaler should be synced with [Kubernetes Autoscaler 1.2x](https://github.com/kubernetes/autoscaler/releases/tag/cluster-autoscaler-1.2x.y)

## Why is this needed

* To keep the fork in sync with upstream.
* To ensure the shoot control planes are using respective version of CA for a given K8S version

## Steps

NOTE: replace `1.2x.0` with the specific autoscaler upstream version that you are trying to sync to. Replace `user` with your github user.

- [ ] Fork `github.com/gardener/autoscaler`  into your github account say `user`. 
- [ ] Check out `github.com/gardener/autoscaler` under `$GOPATH/src/k8s.io`
- [ ] Change to checked-out `autoscaler` dir and set upstream, origin and elankath (fork)
	- [ ] `git remote add user https://github.com/user/autoscaler`
	- [ ] `git remote add upstream https://github.com/kubernetes/autoscaler.git`
- [ ] `git fetch –all`
- [ ] Determine which release branch you wish to sync with
- [ ] `git checkout -b upstream-release-1.2x upstream/cluster-autoscaler-release-1.2x`
- [ ] Go to [autoscaler releases](https://github.com/kubernetes/autoscaler/releases) and get the commit id for release `1.2x.0` . Note down the commitId.
- [ ] `git reset --hard commitId`
- [ ] `git checkout machine-controller-manager-provider`
- [ ] `git pull origin machine-controller-manager-provider`
- [ ] `git checkout -b sync-upstream-v1.2x.0`
- [ ] `git merge upstream-release-1.2x`
	- [ ] This is where changes of master get merged in. advantage of merging is it won’t change the commit hashes of already existing commits. 
	- [ ] Accept all vpa changes.
	- [ ] Accept our changes in `go.mod` and `go.sum` and `vendor` directory.  (The `hack/update-vendor.sh` script below is expected to fix things)
- [ ] Upgrade versions of `machine-controller-manager`, `machine-controller-manager-provider-aws`  and `machine-controller-manager-provider-azure` to the latest release in `go.mod`
- [ ] Run update vendor script after changing to `cluster-autoscaler` directory:  `./hack/update-vendor.sh 1.2x.0`
	- [ ] If the above still gives test issues then use `rsync` or `diff -rq` to figure out differences in `vendor` directory and synchronize them manually.
- [ ] Create a new file `SYNC_CHANGES-1.2x.md`  summarily describing the changes done. Follow template. Use upstream release notes as a guide.
- [ ]  Run Tests
	- [ ] `cd cluster-autoscaler`
	- [ ] `go test $(go list ./... | grep -v cloudprovider | grep -v vendor | grep -v integration)` # core autoscaler unit-tests
	- [ ] `go test $(go list ./cloudprovider/mcm/... | grep -v vendor)` - run MCM cloud provider test
	- [ ] `go test -race $(go list ./... | grep -v /vendor/ | grep -v vertical-pod-autoscaler/e2e | grep -v cluster-autoscaler/integration)`
	- [ ] Verify that binary can be created using: `../.ci/build`
	- [ ] Image can be created using the `Dockerfile` :
	- [ ] Run IT test locallys. 
		- [ ] Follow instructions at: https://github.com/gardener/autoscaler/blob/machine-controller-manager-provider/cluster-autoscaler/integration/usage.md  
		- [ ] Before running `make download-kubeconfigs`, create a folder `mkdir -p dev/kubeconfigs`
		- [ ] This target will print out a list of shell variable `EXPORT` statements. Copy paste thema and execute them
		- [ ] Now run `make test-integration`
- [ ] Update the Availability Matrix 
- [ ] Update any RBAC rules if required
- [ ] Now make a commit.
- [ ] Push the sync branch to fork repo: `git push user sync-upstream-v1.2x.0`
- [ ] Create a PR.
- [ ] Add the `ok-to-test` label to the PR
- [ ] Ensure that integrationt test passes on CICD.
- [ ] Give PR for review.



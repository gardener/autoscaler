# Virtual Pod Autoscaler Recommender Rewrite

Crash loops due to spike in pod resource usage is an observed behaviour for targets autoscaled by vertical pod autoscaler.
This document proposes different solutions for resource recommendation by virtual pod autoscaler to handle spike usage.

## Problem
The current VPA recommender has the following problems which calls for enhancement in the recommending scheme for the VPA to try and auto-heal a component that go into crash loops due to spikes in resource usage.

1. Current VPA recommender implementation considers historic usage of resources by a pod for recommending scale of memory and CPU. This method is not efficient in handling spike in resources usage. When a spike occurs, the recommender lacks information about it in historic usage to bump up the resources and scale up. This can result in frequent component outages.
2. Stagnant CPU recommendations w.r.t increasing limits during scale up and decreasing them during scale down scenarios.
3. OOM Kills and possible livenessProbe failures due to inefficient CPU and memory recommendation. One observation is that in a component that is in crash loop, memory recommendation is increased much more than CPU recommendations especially if OOMKills occur. But many a times, OOMKills occur due to CPU being throttled and recommending only memory scale up in such cases is not efficient. The below recommendation given from a sample where 100 periodic iterations with 10s intervals, of increasing the CPU millicores consumptions by ![resource-consumer](https://github.com/dj80hd/resource-consumer) shows clearly how the CPU recommendation has not scaled according the given requirement. It also shows how the memory is limits are increasing while CPU is remaining same for requests and limits.

```
Note: The experiment to simulate OOM Kill scenario was conducted using ![resource-consumer](https://github.com/dj80hd/resource-consumer).
      The graphs below capture details of a container journey through constant CPU resource consumption. "consume-cpu" API from resouce-consumer was
      used for this experiment. It was bombarded against the service endpoint of the deployment for 100 times with 5 seconds interval between every iteration to generate a sustained high load.
```
We can see the pod restarting once it faces resource crunch at around 14:21 - 14:45 local time in the graph timeline.
![Container Resource Restarting](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/ctr_restarting_resource_crunch.png)

4. For both CPU and memory the recommender suggests a percentile resource unit upgrade from the historical value of resource usage obtained from memory, VeriticalPodAutoScalerCheckpoint or Prometheus client. This does not take into account the current sudden bump in the resource which might be increasing rapidly (exponentially) over a short duration and might be outside any observed range captured in the history. VPA recommender does not suggest a better scale up and scale down for such spikes (it often recommends less than the current usage if limits are configured higher than the requests) resulting in crashes.

The Limits of the container during this 
![Max Container Limits CPU](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/max_ctr_limits_cpu_cores.png)

![Max Container Limits Memory](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/max_ctr_limits_mem_bytes.png)

Please observe the CPU and memory hikes resulting in crash in the below graphs.

![Max Container Requests CPU](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/max_ctr_requests_cpu_cores.png)

![Max Container Requests Memory](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/max_ctr_requests_mem_bytes.png)

Once the sustainance of such high load continues through the time, the container ends up in `OOMKilled` status as shown below.
```
Containers:
  compute:
    Container ID:   docker://d633e681818d3cadd4d336fff5ad03e0862223b6ec5838772b60529e25531538
    Image:          dj80hd/resource-consumer
    Image ID:       docker-pullable://dj80hd/resource-consumer@sha256:4d1bf26a187211d8f2429b6a1c2e6eab3b6cfa2866d4788ccd3331c91dd3931b
    Port:           <none>
    Host Port:      <none>
    State:          Running
      Started:      Fri, 15 Jan 2021 14:45:38 +0530
    Last State:     Terminated
      Reason:       OOMKilled
      Exit Code:    137
      Started:      Fri, 15 Jan 2021 14:31:31 +0530
      Finished:     Fri, 15 Jan 2021 14:45:35 +0530
    Ready:          True
    Restart Count:  1
    Limits:
      cpu:     22m
      memory:  216968644500m
    Requests:
      cpu:     11m
      memory:  144645763
```

One of the recommendations given by the recommender at around the time matching the OOM killing is as follows:
```
Recommendation:
    Container Recommendations:
      Container Name:  compute
      Lower Bound:
        Cpu:     10m
        Memory:  11423184
      Target:
        Cpu:     11m
        Memory:  248153480
      Uncapped Target:
        Cpu:     11m
        Memory:  248153480
      Upper Bound:
        Cpu:     47m
        Memory:  1081116210
```

The VPA recommendation for target limits at around during this time of operation is as follows. We can clearly observe that the CPU recommendation stays constant with the increasing load.

#### VPA Target Recommendation
![VPA target CPU recommendation](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/vpa_target_recom_cpu.png)

![VPA target memory recommendation](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/vpa_target_recom_mem.png)

#### VPA Upper Bound Recommendation
![VPA upper bound CPU recommendation](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/vpa_ub_recom_cpu.png)

![VPA upper bound memory recommendation](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/vpa_ub_recom_mem.png)

#### VPA Lower Bound Recommendation
![VPA lower bound CPU recommendation](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/vpa_lb_recom_cpu.png)

![VPA lower bound memory recommendation](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/vpa_lb_recom_mem.png)

## Goals
- Pluggable new recommender in place of existing VPA recommender.
- React better to CPU and memory usage spikes (especially, avoid recommending lower than resource usage)
- React better to crashloops and possibly achieve automatic recovery if the crashes are triggered by lack of resources
- Scale down faster if usage is consistently lower for some period of time (for the time-being configurable globally via CLI flags because we are re-using the VPA CRD)
- Optionally support recommendations in annotations so that this can be used in parallel with VPA recommender and compared by enhancing vpa-exporter to additionally export recommendations from the annotations.


## Current VPA Recommender Implementation
### Scale Up and Scale Down
- VPA recommender gets its historical resource usage details from memory by default if the pod has not restarted.
- If the pod has restarted information about resource usage is stored in `VerticalPodAutoscalerCheckpoint` structure defined at. 

![VerticalPodAutoscalerCheckpoint](https://github.com/kallurbsk/autoscaler/blob/8cf6e2698a17cf9e7beb928136980cc78d1b9dfb/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1/types.go#L295)
- If a prometheus client is used, A history of metrics for CPU and memory requests usage is obtained from prometheus for a history length of 8 days
- Resource requests metrics for CPU and memory are obtained for each container in a pod from the prometheus client
- The history obtained is initialized with the recommender to start estimating the target, lower and upper bounds for the containers, and there by the pods
- Target Estimator, Lower Bound Requests and Upper Bound Requests are the three main estimations we determine using recommender
- A percentile value for each of these estimators are set as ![follows](https://github.com/kallurbsk/autoscaler/blob/b46a54192da8ab4181346beba13ae1579db82063/vertical-pod-autoscaler/pkg/recommender/logic/recommender.go#L101)
  - 0.9 percentile for target CPU and Memory
  - 0.5 for lower bound requests of CPU and Memory
  - 0.95 for upper bound requests of CPU and Memory
- For Target Estimator we only recommend Percentile Estimation with a safety margin of 15% by default
- An exponential variation in the lower bound requests, upper bound requests which is determined by history of metrics monitored and pod life span. There is a "confidence" value on which the entire resource scaling is determined.
- For lower bound and upper bound estimation a configurable safety margin (by default 15%) is applied
- The current lower bound and upper bound estimation formula is as follows
```
scaledResource = originalResource * (1 + 1/confidence)^exponent

confidence = (LastSampleTimeStart - FirstSampleTimeStart) / (no. of hours per day i.e 24)
```

Where,
- **scaledResource**  : newly scaled integer value of the amount of the 
                    resource  (CPU or Memory)
- **originalResource**: current integer value of the amount of the
                    resource (CPU or Memory)
- **confidence**      : An non-negative integer number that heuristically measures 
                    how much confidence the history aggregated in the AggregateContainerState provides. For a workload producing a steady stream of samples over N days at the rate of 1 sample per minute, this metric is equal to N. It is calcualted as the total count of samples and the time between the first and the last sample 
- **exponent**        : A non-zero float exponent value which is determined by the 
                    targetEstimator, lowedBoundEstimator or upperBoundEstimator. Each of these estimators are set through percentile estimations. targetEstimator is set to percentileEstimator while lowerBoundEstimator and upperBoundEstimator are set as confidence multipliers. For lowerBoundEstimator the exponent is set to -2 and for upperBoundEstimator it is set to +1

## Proposal

To address the concerns mentioned we are proposing this change in the recommendation engine of VPA recommender for autoscaler component.

#### Scale Up:
Doubling of resources based on current resource usage
  1. Today we depend on the resource requests history from prometheus client for a history length of 8 days (default)
  2. Instead, we could obtain current and past 10 minutes usage from prometheus client to get a nearest estimate of how much requests for resources is a pod making
  3. In the event of a spike in resource, we can recommend the target estimator, requests as double the current resource usage. The lower bound and upper bound become the function of target estimator instead of directly recommending based on resource usage using percentile approach. Lower bound can be set to 0.5 times the target estimator value and upper bound can be set to 2 times the target estimator value.
  4. This will ensure handling the spike as it is dependent on the current usage than a hysteresis of resource usage pattern and doubling will react exponentially to a potentially exponential spike in resouce usage.
  5. We consider the following 3 variables (temporary names) threshold_no_of_crashes, last_crash_count_recorded_time, time_length_for_crash_check.
  6. `threshold_no_of_crashes` talk about value being set during the initiation of recommender to threshold number of restarts due to crash beyond which recommender will scale up both CPU and memory. (Will have a default mostly as 3)
  7. `last_crash_count_recorded_time` indicates the last crash count time recorded on the node. (Default: 0)
  8. `time_length_for_crash_check` is `current machine time - last_crash_count_recorded_time`. (Default: 0)
  9. We can set threshold_no_of_crashes to 3 as the default value. If there is a resource crunch even after bumping one of CPU or memory based on current usage and has resulted in a crash, threshold_no_of_crashes serves a threshold beyond which both resources are to be scaled up irrespective of current usage consumption.
  10. The variables `threshold_no_of_crashes, last_crash_count_recorded_time, time_length_for_crash_check` are reset to default if the VPA recommender itself restarts. Example: Let's say we recommended the memory doubling but had constant crashes for more than 3 times. Then we get into else condition and double both memory and CPU irrespective of whatever the recommendation for individual resource is.

#### Scale Down:
Scale down to a variable between 1 and 2 times the local maxima within a time window
  1. Two new variables called scale_down_safety_margin (default to 1.2), scale_down_monitor_time_window (default to 1 hour) are used.
  2. We find the local maxima w.r.t resource usage within scale_down_monitor_time_window.
  3. We only scale down to `scale_down_safety_margin` times the maximum resource usage recorded within `scale_down_monitor_time_window`.
  4. `1 < scale_down_safety_margin < 2` (defaulted to 1.2). This ensures we always scale down to, slightly higher than previously seen highest resource usage within the time window but also lesser than scale up limit of doubling.
  Note: It could also be given as a best practise to ensure `1 < scale_down_mutliple < 1.5` so that we do not over provision the resources while scaling down. As the parameter scale_down_safety_margin is configurable between 1 and 2, there is room for trial and error and eventually arrive at a custom combination of parmater values for each component based on their usage of VPA.
  5. It also is crucial to scale down stepwise with shorter time windows. This ensures:
      1. Scale Down is gradual and does not vary drastically unlike percentile method during spike on and spike off events
      2. Scale Down is always done with nearest local maxima of resource usage which ensures we are not unncessarily taking history into account to avoid over or under provision resources.


### Annotation based recommender variables
  1. To configure the VPA recommender to work in parallel with the existing VPA recommender a view only mode recommendation is provided through these annotation variables.
  2. The variables depicting the recommendations given by the new recommender will be exported by the vpa-exporter to add values directly to these annotations.

### Example simulation of proposed target recommendation
As you can see below, the container went into a `CrashLoopBackOff` by 18:51 around. Once the VPA recommender was edited to double the current peak witnessed in the usage, the container came back up from `CrashLoopBackOff` serving traffic.
One can see from the below image that the usage of memory slowly increased peaking to 20GB. The new VPA recommender for memory suggested target to be set to 40GB for tackling this scenario.
![CPU usage during the peak resulting in crash loop](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/proposal_manual_simulation/simu_ctr_cpu_usage_seconds.png)

![Memory usage during the peak resulting in crash loop](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/proposal_manual_simulation/simu_ctr_memory_working_bytes.png)

![CrashLoopBackOff with multiple restarts](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/proposal_manual_simulation/simu_ctr_restarts.png)

For the above container which had a stress at 20GB memory and 600 millicores CPU, the container crashed and ended up in backOff loop. To remove it from there a manual simulation of CPU 1200 millicores and 40GB memory were initiated. 

![Target CPU Recommendation](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/proposal_manual_simulation/simu_vpa_target_cpu.png)

![Target Memory Recommendation](https://github.com/kallurbsk/autoscaler/blob/machine-controller-manager-provider/vertical-pod-autoscaler/docs/images/vpa-recommender/proposal_manual_simulation/simu_vpa_target_cpu.png)

### Concerns:
  1. Doubling of target might be too generous and might hit upper limit of the resource availability.
      1. As we are monitoring current resource usage over historical data, this should ideally not be a problem for scale up.
      2. For scale down, we have the proposed scale_down_monitor_time_window post which the scale down recommendation for the pod is issued gracefully.
  2. Will the new VPA recommender handle all the cases of VPA
      1. The current VPA recommender is failing primarily during spikes causing OOM KIlls as it is based on historical data.
      2. The new VPA recommender which scales up primarily on the current usage and doubling the target estimation should cover most of the worst case spikes.
  3. Doubling of only memory or CPU
      1. This is avoided by using threshold_no_of_crashes variable as mentioned in Scale Up proposal.
      2. threshold_no_of_crashes is a customisable variable and can be set according to the application time over which VPA is used for
      3. So if only one resource is getting scaled up and even then we see restart of pods due to OOMKills, then we can double both the resources beyond a certain threshold for threshold_no_of_crashes
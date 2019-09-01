---
title: "Squeezing GKE System Resources In Small Clusters"
date: 2019-08-28T18:00:00-01:00
author: "@alexdmoss"
description: "Trying to get as much out of a tiny GKE cluster as possible by squeezing the system resources"
banner: "/images/squeeze-1.jpg"
tags: [ "GKE", "Kubernetes", "GCP", "VPA", "Autoscaling", "Cost" ]
categories: [ "GKE", "Kubernetes", "Cost" ]
draft: true
---

{{< figure src="/images/squeeze-1.jpg?width=600px&classes=shadow" attr="Photo by Davide Ragusa on Unsplash" attrlink="https://unsplash.com/photos/cDwZ40Lj9eo" >}}

**Spoiler Alert!** This blog is *really* about Vertical Pod Autoscaling. It just might not sound like it at the start :smile: If you're not interested in how I got there and just want to jump to the good stuff - I've summarised it at the bottom of the post.

---

So, this all started because I wanted to try out Elasticsearch. I'm not the biggest fan of GCP's Stackdriver Logging, and I wanted to poke around with some alternatives. I have a "personal" GKE cluster that I use to run a few websites (including this blog!), and where I also experiment with these sort of things outside of work. Should be straight-forward I thought to myself. Sure, it's Java - but I don't need it to do that much, I'll just shrink it down enough that it'll start up ok, play around a bit, and then maybe see if I want to do more with it rather than sign up with some sort of managed logging service instead.

Well, I was dead wrong about that. I quickly realised that there's no way it was going to fit - even with the tiny workloads, too much of the machine resource was gobbled up by things like processes running in `kube-system` and I really didn't want to extend beyond my 3 x g1-small instances if I could possibly help it.

So I could make them temporarily larger while I experiment. But where's the fun in that?! What follows is a look at what's gobbling up all the resource - because let's be honest, the handful of services I run on my GKE aren't going to be doing a huge amount ...

---

## So, What's Going On?

*Visualise the Problem*. No, not some sort of zen meditative technique. What are we actually dealing with here? Surely a dozen or so deployments of tiny web-apps isn't gobbling up 3x1.7Gb? Have I been a responsible engineer and set some resource requests/limits (erm ...)?

As with all things in life, we start with `bash`:

```bash
kubectl get po $1 \
-o custom-columns=NAME:.metadata.name,CPU_REQ:.spec.containers[].resources.requests.cpu,CPU_LIMIT:.spec.containers[].resources.limits.cpu,MEMORY_REQ:.spec.containers[].resources.requests.memory,MEMORY_LIMIT:.spec.containers[].resources.limits.memory
```

**add note about krew & resource-capacity here**

This ugly-as-hell one-liner produces something that looks a bit nicer than you'd think:

```bash
[~ (⎈ |my-cluster:default)]$ get-pod-resource-requests.sh kube-system
NAME                                                        CPU_REQ   CPU_LIMIT   MEMORY_REQ   MEMORY_LIMIT
event-exporter-v0.2.5-7d99d74cf8-9jdn4                      <none>    <none>      <none>       <none>
fluentd-gcp-scaler-55bdf597c-xsxwr                          <none>    <none>      <none>       <none>
fluentd-gcp-v3.1.1-6r6x5                                    10m       100m        100Mi        250Mi
fluentd-gcp-v3.1.1-mzgdn                                    10m       100m        100Mi        250Mi
heapster-v1.6.1-5b5df74474-vjvhd                            13m       13m         120Mi        120Mi
kube-dns-6987857fdb-42kzx                                   100m      <none>      70Mi         170Mi
kube-dns-6987857fdb-w8bd4                                   100m      <none>      70Mi         170Mi
kube-dns-autoscaler-bb58c6784-rsklz                         20m       <none>      10Mi         <none>
kube-proxy-gke-mw-prod-np-1-2cfd6748-k8kg                   100m      <none>      <none>       <none>
kube-proxy-gke-mw-prod-np-1-68476602-7zlg                   100m      <none>      <none>       <none>
l7-default-backend-fd59995cd-vcmd4                          10m       10m         20Mi         20Mi
metrics-server-v0.3.1-57c75779f-zmbt6                       43m       43m         55Mi         55Mi
prometheus-to-sd-rgdvx                                      1m        3m          20Mi         20Mi
prometheus-to-sd-smrbm                                      1m        3m          20Mi         20Mi
stackdriver-metadata-agent-cluster-level-55dfd764dd-8q674   40m       <none>      50Mi         <none>
```

Useful, but not sexy. I wanted a view of the packing onto my nodes. For this I dusted down a copy of https://github.com/hjacobs/kube-ops-view/, which is something I'd played around with before. My deployment of it is largely based on the sample yaml's with a view tweaks to secure it and make the Redis pod a little happier: https://github.com/alexdmoss/kube-ops-view

This gives something visually more appealing:

![Kube Ops View Before Changes](/images/squeeze-gke-ops-view-before.png?width=600px)

This is for a brand new node pool. With those amber bars (mostly for memory), we're not going to be able to fit something like Elastic on here. Actually, even with a third node, we'd struggle because the kube-system stuff expands over time - it doesn't have caps on its resource utilisation as we see above from the sexy-bash.

When we try to deploy a small EFK stack (already shrunk from sizeable defaults), we end up with something like this:

![EFK won't fit!](/images/squeeze-gke-efk-wont-fit.png?width=600px)

```sh
[~/personal/github/elastic-on-gke (⎈ |mw-prod:elastic)]$ k get po
NAME                          READY   STATUS     RESTARTS   AGE
efk-elasticsearch-0           0/2     Pending    0          7m34s
efk-fluentd-es-bjp4m          1/1     Running    0          7m34s
efk-fluentd-es-zn5r7          1/1     Running    0          7m34s
efk-kibana-5bf6df5fc7-6p9h5   1/1     Running    0          7m34s
efk-kibana-7877d95c-nc4cp     0/1     Running    1          7m34s
efk-kibana-init-job-f2jph     0/1     Init:0/1   0          7m34s
```

```sh
k describe po efk-elasticsearch-0

Events:
  Type     Reason            Age                  From               Message
  ----     ------            ----                 ----               -------
  Warning  FailedScheduling  10m (x5 over 10m)    default-scheduler  pod has unbound immediate PersistentVolumeClaims
  Warning  FailedScheduling  108s (x21 over 10m)  default-scheduler  0/2 nodes are available: 1 Insufficient memory, 1 node(s) had volume node affinity conflict
```

What happens if we terraform ourselves a 3rd node? The answer is yes it does manage to fit on, but it's a tight fit, and we've squeezed the Elasticsearch pod down to a fairly measly 500Mb or so. Adding a second replica shows that it still won't fit on one of the other two nodes:

![Three nodes lets us get a single replica running](/images/squeeze-gke-efk-one-replica.png?width=600px)

We also know that this only works on a brand new cluster. When I tried this before, it was no such luck, as some of the pods without requests/limits had expanded under use, leaving insufficient room. With that in mind, lets get some guardrails in.

---

## Someone Else Has Solved This, Right?

As is standard in our lovely industry, the first thing I did was of course Google it. Someone else must've had this problem before, right? Turns out, Google themselves - or rather, they've been asked the question before and produced [a handy document as a starting point](https://cloud.google.com/kubernetes-engine/docs/how-to/small-cluster-tuning).

This sounded perfect! I wouldn't have been writing a blog post if that were true however :sad

There *is* some good advice in here - but a lot of it didn't really apply for me. I wanted to keep Stackdriver Logging/Monitoring switched on in the meantime while I was experimenting - and might need the metric exports anyway later. The Kubernetes Dashboard was already switched off because it has a history of security thingies **LINK**

Some of the advice is great though - I don't need HPA so could switch that off, I pondered scaling down Kube DNS but in the end decided to leave it. The really interesting one was `ScalingPolicy` for fluentd - I followed this recommendation with a ScalingPolicy as follows and it worked a treat for that component:

```yaml
---
apiVersion: scalingpolicy.kope.io/v1alpha1
kind: ScalingPolicy
metadata:
  name: fluentd-gcp-scaling-policy
  namespace: kube-system
spec:
  containers:
  - name: fluentd-gcp
    resources:
      requests:
      - resource: cpu
        base: 10m
      - resource: memory
        base: 100Mi
      limits:
      - resource: cpu
        base: 100m
      - resource: memory
        base: 250Mi
```

Super-impressed I was. I'll just do the same thing for other stuff, like `kube-dns`, `kube-proxy`, etc. NOPE! Turns out this relies on an additional component running in GKE - `fluentd-gcp-scaler` - which is based on this: https://github.com/justinsb/scaler.

Cool concept - so I did spend a little bit of time trying to run my own equivalent of this in my cluster ... and in so doing, realised there was probably a better way ...

## I Need More Control - Enter VPA

When I say "more control", really I mean - relinquish control completely and have something else sort it out for me :smile Google recently announced Vertical Pod Autoscaler (VPA) in Beta - few months back I think. **LINK**. It's an addon for GKE Advanced, so will carry a cost at some point - but the code itself is open sourced too.

The [google_container_cluster Terraform resource](https://www.terraform.io/docs/providers/google/r/container_cluster.html#vertical_pod_autoscaling) already contains the option to enable GKE's VPA addon, so turning this on was a breeze:

```terraform
resource "google_container_cluster" "cluster" {

  # enabling VPA needs Beta
  provider                        = "google-beta"

  # [other important stuff]
  
  vertical_pod_autoscaling {
    enabled                       = true
  }

}
```

As part of this I also flipped on Cluster Autoscaling just to see if the behaviours involved here caused my cluster to flex in size (I'd removed my EFK deployment at this point, so we were back at square one).

Getting VPA to do its thing involves applying some straight-forward policy, which looks a bit like this:

```yaml
apiVersion: autoscaling.k8s.io/v1beta2
kind: VerticalPodAutoscaler
metadata:
  name: kube-dns-vpa
  namespace: kube-system
spec:
  targetRef:
    apiVersion: extensions/v1beta1
        kind: Deployment
        name: kube-dns
  updatePolicy:
    updateMode: "Off"`
```

With this yaml, we create a VPA policy in Recommendation mode - describing the VPA then tells us what it thinks the resource bounds should be. That's all well and good, but more fun is changing this to `updateMode: "auto"` and letting it resizing these pods for you.

Here's the state of kube-system before enabling VPA in update mode:

![Before the Squeeze](/images/squeeze-gke-before-vpa.png?width=600px)

I set some VPAs up for all the things in `kube-system` and left it for a few days to do its thing. I ended up with the following:

Original:

```sh
event-exporter-v0.2.5-7d99d74cf8-6blm5                      <none>    <none>      <none>       <none>
fluentd-gcp-v3.1.1-mzgdn                                    10m       100m        100Mi        250Mi
fluentd-gcp-v3.1.1-sjztp                                    10m       100m        100Mi        250Mi
fluentd-gcp-scaler-55bdf597c-xsxwr                          <none>    <none>      <none>       <none>
heapster-v1.6.1-74885ff99d-wkjl9                            13m       13m         120Mi        120Mi
kube-dns-6987857fdb-7jm2z                                   100m      <none>      70Mi         170Mi
kube-dns-6987857fdb-w8bd4                                   100m      <none>      70Mi         170Mi
metrics-server-v0.3.1-57c75779f-zmbt6                       43m       43m         55Mi         55Mi
```

New:

```sh
event-exporter-v0.2.5-7d99d74cf8-drbv2                      12m       <none>      131072k      <none>
fluentd-gcp-v3.1.1-27gs7                                    23m       230m        225384266    563460665
fluentd-gcp-v3.1.1-8bhv4                                    23m       230m        203699302    509248255
fluentd-gcp-v3.1.1-rmh47                                    23m       230m        203699302    509248255
fluentd-gcp-scaler-55bdf597c-dzxf5                          63m       <none>      262144k      <none>
heapster-v1.6.1-74885ff99d-8gksk                            11m       11m         87381333     87381333
kube-dns-6987857fdb-7jh59                                   11m       <none>      100Mi        254654171428m
kube-dns-6987857fdb-dqtft                                   11m       <none>      100Mi        254654171428m
metrics-server-v0.3.1-57c75779f-mpdvq                       12m       12m         131072k      131072k

```

Not resized:

```sh
kube-dns-autoscaler-bb58c6784-gdxnp                         20m       <none>      10Mi         <none>
kube-proxy-gke-mw-prod-np-1-2cfd6748-kfw1                   100m      <none>      <none>       <none>
kube-proxy-gke-mw-prod-np-1-2cfd6748-z2cc                   100m      <none>      <none>       <none>
kube-proxy-gke-mw-prod-np-1-68476602-p1jq                   100m      <none>      <none>       <none>
l7-default-backend-fd59995cd-gn95w                          10m       10m         20Mi         20Mi
prometheus-to-sd-97cm7                                      1m        3m          20Mi         20Mi
prometheus-to-sd-cd884                                      1m        3m          20Mi         20Mi
prometheus-to-sd-wptpv                                      1m        3m          20Mi         20Mi
stackdriver-metadata-agent-cluster-level-55dfd764dd-z928x   40m       <none>      50Mi         <none>
```

All kinds of wacky units involved! Ironically, you can see from the "empty" node that this has actually had a negative effect, increasing the amount of memory requests (second vertical bar from the left)! This actually reflects the fact that the VPA has only been running for a couple of hours, and set some requests on components that hadn't been around that long and had no resource requests/limits defined. Given a little longer, it will "learn" a more appropriate set of ranges.

A handy extension that can be made here is to set your own upper/lower bounds - particularly useful for situations where workloads can be spiky or extra resource is needed for pod initialisation, for example:

```yaml
spec:
  resourcePolicy:
    containerPolicies:
    - containerName: '*'
      maxAllowed:
        memory: 2Gi
      minAllowed:
        memory: 100Mi
```

Lets see what it looks like after 24 hours ...

** update here **

Sadly, I could not get the setup working for kube-proxy - which appear to be individual pods not spawned from a deployment/daemonset (or at least, not one that's visible to me). Some docs suggest that VPA can target resources based on labels, but this seems to have been deprecated.

```sh
Status:
  Conditions:
    Last Transition Time:  2019-08-30T17:17:31Z
    Message:               Label selector is no longer supported, please migrate to targetRef
    Status:                True
    Type:                  ConfigUnsupported
    Last Transition Time:  2019-08-30T17:18:31Z
    Message:               Fetching history failed: Cannot construct metric filter. Reason: Selector  doesn't want to select anything
    Reason:                2019-08-30T17:18:31Z
    Status:                False
    Type:                  FetchingHistory
    Last Transition Time:  2019-08-30T17:17:31Z
```

## AutoScale ALL THE THINGS

This seemed effective enough that I fancied rolling it out to my actual workloads too (rather than just the stuff in kube-system) - with that in mind I created a lightweight controller (inspired by some of the work my colleagues have done) - code for it is here: https://github.com/alexdmoss/right-sizer. This will skim through Deployments every 10 mins and create VPA Policies for any new workloads it spots. This had rather comedic effects with `updateMode: Auto` :smile

![Uh-oh - BIG cluster!](/images/squeeze-gke-all-auto-vpa.png?width=600px)

As can be seen by the red bars - this happened a few minutes after the VPA policies were created and isn't super-suprising - all those tiny pods of nginx etc were rapidly set with requests of 200-500Mi. For nodes with only 1Gb of spare RAM available, there was no choice but for the Cluster Autoscaler to kick in!

I did this as a bit of an experiment - but it's obvious that we need to be careful with this stuff. The VPA only has limited info to go on, and unless you set `resourcePolicies` to cap it to sensible values (not so practical for a Controller that applies to all Deployments!) it can do some wacky things. For this reason, I switched things back to recommend-only mode for all my workloads, and then used this data to set sensible defaults that I was happy with for my pods.

## Conclusions

{{< figure src="/images/squeeze-2.jpg?width=600px&classes=shadow" attr="Photo by Josh Appel on Unsplash" attrlink="https://unsplash.com/photos/NeTPASr-bmQ" >}}

So, can I run Elastic now?

** add kube ops view view **

## Evolution

I was sufficiently impressed with VPA that it seemed worth a closer look in a work context too. At a larger scale the savings in Compute could become quite significant - depending on the maturity of testing practices in teams (some of ours are *very* good at right-sizing their workloads already).

We've recently enabled it in Recommendation mode and started bringing the results into Prometheus and visualising them against current utilisation in Grafana - early days, but looks really cool and I'd like to replicate it in my home setup. Some of the recommendations are pretty quirky though, so it may need a bit more time to bed in ... and letting it auto-resize may not be viable for us given the amount of JVM-based workloads we run (it can't also set -Xmx ...).

## In Summary

1. Get a tool that helps you visualise your resource requests/limits/utilisation. I like `kube-ops-view` as it's simple but effective
2. In GKE, enable the VerticalPodAutoscaler addon and apply some VPA policies targeting the deployments you are interested in. I started in "Recommend" mode to see what it was going to do first
3. If you're comfortable with the recommendations and that your workloads can tolerate the restarts - switch on update mode and profit!
4. If you'd like a controller to set VPAs for all your deployments, have a nose at this for inspiration: https://github.com/alexdmoss/right-sizer

**Add some visualisations here**

What about kube-proxy - should I write a controller?

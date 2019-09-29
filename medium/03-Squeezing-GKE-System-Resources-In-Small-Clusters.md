---

Squeezing GKE System Resources in Small Clusters
Photo by davide ragusa on Unsplash

Spoiler Alert! This blog is really about Vertical Pod Autoscaling and patching of kube-systemworkloads in GKE. I've summarised the steps I went through at the bottom of the post.
To read about the steps I took in a bit more detail, see my personal blog post here: https://mosstech.io/2019/09/28/squeezing-gke-system-resources/.


---

If, like me, you run a small GKE cluster of your own to try out things - perhaps running a handful of small websites on it - then you may find it uncomfortable of your available Compute resources are used up just keeping Google Kubernetes Engine itself running. This is not surprising - Google are going set things up on the assumption you're running things at a reasonable scale, and want to use many of the extra features they've set up for you. But in this case, you don't.
So, before you run off to AppEngine, there's a few tricks that can be employed to shrink down those resources to make room for more stuff.

---

First Things First - What's Going On?
Before we start, we need some way of seeing what's going on - what's actually chewing up those resources, and by how much?
Now, it's certainly possible to do this with a bit of kubectl foo or kubectl resource-capacity which I recently discovered via the excellent krew kubectl plugin installer, but I wanted a better way of visualising the packing onto the nodes.
For this, I deployed a copy of kube-ops-view, a handy utility I'd played with in the past. It needed a bit of tweaking to get the sample yaml files working, after which I had the following:
Elasticsearch won't fit!The red circle shows the unscheduled Elasticsearch pod. You can also hover over each pod for more info on its resource utilisation, and filter using the form across the top.
The vertical bars - in particular the amber ones, which show memory - tell us that we're not going to be able to fit something this size on here (and Kubernetes won't re-pack the existing workloads as it is honouring their anti-affinity rules). I actually found that even with a third node, it'd struggle because the kube-system workloads expand over time - many of them don't have caps on resource utilisation. This leaves us with pods in a Pending state and FailedScheduling errors due to insufficient memory.

---

Someone Else Has Solved This, Right?
Naturally, the first thing I did was, of course, to Google it. Turns out Google themselves have produced a handy document as a starting point for folks running on GKE. Whilst there is some good advice in here, a lot of it didn't really apply for me. I didn't want to switch off Stackdriver stuff completely, and things like the Kubernetes Dashboard were already switched off.
One thing from this documentation that did help me was using a ScalingPolicy on fluentd - I followed this recommendation with a policy as follows and it worked a treat:

I was really impressed with this - initially thinking I could use the same technique for a bunch of other things in kube-system until I realised it worked through the presence of a component installed by GKE called fluentd-gcp-scaler - which is based on this. Whilst it might be possible to jury-rig my own implementation of this, I instead switched tack towards something I'd wanted an excuse to try for a while …

---

Enter: Vertical Pod Autoscaling
VPA is a feature that Google recently (ish) announced in Beta. It is talked about in the context of GKE Advanced / Anthos, so I may need to keep an eye on whether it becomes a chargeable product 💰 - but in the meantime, it seemed worth experimenting with.
The google_container_cluster Terraform resource already contains the option to enable GKE's VPA addon, so turning this on was a breeze:

```sh
resource "google_container_cluster" "cluster" {
  provider = "google-beta"
  # [... other important stuff ...]
  vertical_pod_autoscaling {
     enabled = true
  }
}
```

Getting VPA to do its thing involves applying some straight-forward policy, which looks a bit like this:

With this yaml, we create a VPA policy in Recommendation mode - describing the VPA then tells us what it thinks the resource bounds should be, and it doesn't seem to take too long to start making recommendations for you.
[~]$ kubectl describe vpa metrics-server-v0.3.1
Name:         metrics-server-v0.3.1
Namespace:    kube-system

Status:
  Recommendation:
    Container Recommendations:
      Container Name:  metrics-server
      Lower Bound:
        Cpu:     12m
        Memory:  131072k
      Target:
        Cpu:     12m
        Memory:  131072k
      Uncapped Target:
        Cpu:     12m
        Memory:  131072k
     Upper Bound:
       Cpu:     12m
       Memory:  131072k
That's all well and good, but more fun is changing this to updateMode: "auto" and letting it actually perform the resizing these pods for you. A handy extension to your VPA definitions that can be made here is to set your own upper/lower bounds - particularly useful for situations where workloads can be spiky or extra resource is needed for pod initialisation. For example:

I set some VPA definitions up for all the things in kube-system and left it for a short while to do its thing. I ended up with the following:
Before:
NAME                         CPU_REQ   CPU_LIMIT   MEM_REQ   MEM_LIMIT
event-exporter-v0.2.5-[..]   <none>    <none>      <none>    <none>
fluentd-gcp-v3.1.1-[..]      10m       100m        100Mi     250Mi
fluentd-gcp-v3.1.1-[..]      10m       100m        100Mi     250Mi
fluentd-gcp-scaler-[..]      <none>    <none>      <none>    <none>
heapster-v1.6.1-[..]         13m       13m         120Mi     120Mi
kube-dns-[..]                100m      <none>      70Mi      170Mi
kube-dns-[..]                100m      <none>      70Mi      170Mi
metrics-server-v0.3.1-[..]   43m       43m         55Mi      55Mi
After:
NAME                        CPU_REQ   CPU_LIMIT   MEM_REQ     MEM_LIMIT
event-exporter-v0.2.5-[..]  12m       <none>      131072k     <none>
fluentd-gcp-v3.1.1-[..]     23m       230m        225384266   563460665
fluentd-gcp-v3.1.1-[..]     23m       230m        203699302   509248255
fluentd-gcp-v3.1.1-[..]     23m       230m        203699302   509248255
fluentd-gcp-scaler-[..]     63m       <none>      262144k     <none>
heapster-v1.6.1-[..]        11m       11m         87381333    87381333
kube-dns-[..]               11m       <none>      100Mi       254654171428m
kube-dns-[..]               11m       <none>      100Mi       254654171428m
metrics-server-v0.3.1-[..]  12m       12m         131072k     131072k
All kinds of wacky units involved! As well as occasionally recommending things at a larger size than I hoped, I also couldn't get it to target certain resources - namely the kube-proxy pods which aren't DaemonSets as expected, but individual pods in the world of GKE (weird, right?). VPA unfortunately only works based on a targetRef field (rather than something like a label selector, which it looks like it used to support but now no longer does).

---

AutoScale ALL THE THINGS … Oops!
Intrigued by where this would lead, I decided to have some fun with the actual workloads running in my cluster too.
With that in mind I created a lightweight controller (heavily inspired by the work some of my colleagues had done) - code for it is here: https://github.com/alexdmoss/right-sizer. This will skim through Deployments every 10 mins and create VPA Policies for any new workloads it spots. This had rather comedic effects with updateMode: "auto", as can be seen by this screenshot from kube-ops-view:
This happened a few minutes after the VPA policies were created and isn't super-surprising when you think about it. All those tiny pods of mostly nginx were getting set with a memory request of 200–500Mi, creating memory pressure on the nodes as can be seen by the red bars. For nodes with only 1Gb of spare RAM available, there was no choice but for the Cluster Autoscaler to kick in.
Naturally I was just experimenting here, and I had a reasonably tight leash on the Cluster Autoscaler, so no harm done - I could quickly revert the VPA to recommend-only mode and reset my workloads. However, it's clear (and unsurprising!) that you need to be careful with this when using it for "proper" workloads in Production (or where you run the risk of running up a really big bill). Setting a resourcePolicies as mentioned earlier can help with this - but of course not so practical when you've cheated with a Controller that sets VPA for all Deployments as I have 😅

---

Right-sizing kube-system
So, I still had a problem. I had some recommendations from VPA, but the beefier workloads on my tiny cluster were mostly the ones residing in kube-system, for which I don't control the PodSpec.
It's here that things get a little bit hacky 😀. I extended my Controller (used earlier to create the VPA policies) to also set some resource requests/limits on the kube-system resources that were on the larger side. The code for it is really quite awful (it was a quick proof-of-concept, honest!) and given it has been ticking away for a few weeks now and seems to be working out okay, I really should clean it up 🐛
It works by periodically (every 10 mins) patching the pods in kube-system with new entries for memory and CPU utilisation. I opted to do this for kube-dns, heapster and metrics-server (fluentd was covered by the ScalingPolicy mentioned earlier).
Fundamentally, the code is not complicated:

As a general approach, I don't think I'd be comfortable doing this sort of thing in Production for an important system. Google will no doubt be reconciling these things themselves - we're effectively overwriting their chosen settings on a more frequent basis - hence my "hacky" comment. But, for a small cluster used for less important things, this approach does seem to have worked out okay.

---

Conclusions
Photo by Josh Appel on UnsplashSo, can I run Elastic now? Yes!
As can be seen above, we've just about got it squeezed onto a pair of n1-standard-1's.
More useful for me really is that I learnt some interesting things about VPA and right-sizing of workloads, and some idiosyncrasies in how GKE manages pods in the kube-system namespace. Fun times! 🐳
Evolution?
I was sufficiently impressed with VPA in an advisory capacity that it seemed worth a closer look in a work context too. At a larger scale the savings in Compute could become quite significant - depending on the maturity of testing practices in your teams (some of our teams at work are very good at right-sizing their workloads already, whereas some could probably use the help).
We've therefore recently enabled it in Recommendation mode and started bringing the results into Prometheus and visualising them against current utilisation through Grafana - early days, but it looks really cool and I'd like to replicate the same in my home setup. Some of the recommendations are pretty quirky though, so it may need a bit more time to bed in … and letting it auto-resize is unlikely to be viable for us given the amount of JVM-based workloads we run (it can't also set -Xmx, for example).

---

In Summary
Get a tool that helps you visualise your resource requests/limits/utilisation. I like kube-ops-view as it's simple but effective
Google have an article with specific advice for GKE. Many of the recommendations won't be suitable, but things like disabling the Kubernetes Dashboard and taking advantage of their fluentd autoscaler are good quick wins
In GKE, enable the VerticalPodAutoscaler addon and apply some VPA policies targeting the deployments you are interested in. I started in "Recommend" mode to see what it was going to do first
If you'd like a custom controller to setup VPA for all your deployments, have a nose at this for inspiration
If you're comfortable with the recommendations and that your workloads can tolerate the restarts - switch on update mode and forget about needing to right-size your workloads (… in theory)
To really squeeze things down, you can update kube-system resources with a custom controller which does the equivalent of kubectl patch on the resource requests/limits on a regular basis

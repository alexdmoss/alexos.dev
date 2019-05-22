---
title: "KubeCon 2019 - Day 1 Reflections"
date: 2019-05-21T09:00:00-01:00
author: "@alexdmoss"
description: "My reflections on a busy first day of KubeCon"
banner: "/images/barcelona-1.jpg"
tags: [ "KubeCon", "CNCF", "opinion", "reflections" ]
categories: [ "Conference", "Kubernetes", "CNCF" ]
---

{{< figure src="/images/barcelona-1.jpg?width=600px&classes=shadow" title="Photo by Enes on Unplash" >}}

<!-- https://unsplash.com/photos/f-DvU93UhTs -->

> Wow - first day of KubeCon and quite an information overload! Don't know how my brain is going to cope with two more days of this - we will see!

In trying to sum things up, I find myself struggling for any stand-out themes. Thinking about this a little more - this is probably fairly natural in comparison to "big vendor" conferences where they have large product & marketing teams driving their goals. Here however, it's a community event of largely open source (or at the very least, open to contribution and collaboration) software. Even the larger vendors embrace the spirit of it all. It's actually really refreshing.

One personal theme I do reflect on is how nice it is to actually hear from real customers and users of this stuff. Sure, I went to a few talks from the big tech companies, but even then they're presenting very experience-focused content, their speakers focusing on how they're engaging with the community as a whole, and so on. Not at all sales-pitch-y. It's pretty brilliant.

---

## The Higgs-Boson

The personal highlight for me was during the end of day keynote, watching a software engineer + physicist from CERN on stage in front of 7,000 people run a live demo recreating the analysis used to detect the Higgs-Boson particle.

{{< figure src="/images/higgs-1.jpg?width=600px&classes=shadow" title="Higgs-Boson CMS Experiment - Kit" >}}

As a long-ago-sadly-largely-forgotten scientist, this was pretty sweet stuff. The scale of the physical experiment kit, the need to basically take a photo of a particle that disappears in a billionth of a trillionth of a second (actually, I think they take a photo of the after-effect of this!) - it's both staggering and fascinating.

{{< figure src="/images/higgs-2.jpg?width=600px&classes=shadow" title="Higgs-Boson - Decay Model" >}}

Perhaps more incredibly, the computing is actually more mainstream (for this audience, at least), than I thought. Just a lot of it! They have 70TB of data across 25k files, running across a Kubernetes cluster with 20k+ cores. They use Ceph, and Redis, and Jupyter - stuff I've heard of and in some cases fiddled around with.

Funnily enough CERN wouldn't let them borrow their actual kit for a demo - so instead they called in a favour from Google to run it on GCP instead.

> Yay Public Cloud!

This turned out to be a great advert for Google - we see them bursting up to 180Gb/s pulling from Object Storage which is crazy rapid. It does make you wonder, given they're running this on GCS + GKE - how are they mounting that object storage to get those speeds?!

To recreate it in an authentic fashion, they use the data from the experiment itself, and the software used back then (2010 - go-go-gadget Docker) - despite this they're returning a successful result in just a couple of minutes.

![Go-Go-Gadget](/images/gadget.png)

> Apparently anyone can do this as the data & source code is public - I wouldn't fancy paying the GCP bill though :wink:

---

Moving on to slightly more real-world things (awwww) - the brief but highly relevant announcement of [SMI](https://smi-spec.io/) (Service Mesh Interface) is welcome. We aren't running a service mesh in my team currently but I'm convinced we will be in not too long, so having these things start to coalesce around a standard can only be a good thing. 

Their goal is ([I'm paraphrasing massively](https://cloudblogs.microsoft.com/opensource/2019/05/21/service-mesh-interface-smi-release/)):

- common portable APIs across different service mesh technologies - works with the three bigger players (Istio, Linkerd, Consul)
- apps, tools, the wider ecosystem will all be able to integrate through these standardised APIs

It is being done to help folks get started quickly - i.e. choice is believed to be a barrier, simpler is better, tapping into the ecosystem is good

> My personal opinion on this is focused more around the "simpler is better" comment. I think the barrier for us when it comes to meshes is the complexity it introduces - we would be able to make a choice, and we would trust the wider ecosystem to adapt, but we (and in particular our wider developer community) are ourselves still only just getting our heads round Kubernetes itself, if we're being honest about it

It was a brief segment in the end-of-day keynote, so I would imagine a "watch this space" situation. The video demo was snazzy though if you want to check it out.

---

## The Kubernetes Project Update

Janet Kuo did a very nice segment on the Kubernetes project in CNCF itself. I liked her timeline narrative, coming from a "It's Kubernetes' 5th birthday!" perspective:

{{< figure src="/images/borg.jpg?width=600px&classes=shadow" title="Uh oh, the Borg!" >}}

- 2003 - Borg was created. No, not in Star Trek, this is Google's collective instead. That was, staggeringly enough, May 1989. Oh man now I feel old
- 2006 - cgroups arrived (Linux Control Groups ==> process containers)
- 2008 - Linux containers adopted the "container" naming
- 2009 - Omega - next-gen Borg at Google -> influences Kubernetes design
- 2013 - Docker. Hooray for Docker, it's the best
- 2013 - Project 7 (within Google, open source container orchestrator - comes from Seven of Nine from Star Trek of course - the Friendly Borg)
- 2014 - Kubernetes announced at DockerCon. There was much rejoicing
- 2015 - v1.0 of Kubernetes, CNCF formed, 1st KubeCon at end of year
- 2016 - SIGs came into being, KubeCon EU, industry adoption, in Production, at Scale. This thing is a big deal
- 2017 - CRD introduced so can define your own Kubernetes to build on top of. Kubernetes becomes a platform for building platforms (if I, or rather Kelsey Hightower, had a pound for every time that was mentioned at this conference ...). Cloud adoption drives Kubernetes to become a *de facto standard*, many new native APIs introduced
- 2018 - Graduated from incubation in the CNCF
- Today - one of the highest velocity open source projects. It is number 2 on Pull Requests on Github (behind ... Linux), #4 on issues/authors (out of 10's of millions)

Now it is v1.14 and the most stable and mature ever. And hell you can even run it on Windows nodes. Now *that's* crazy.

---

## Deeper Dives

In some of the specialist tracks I went to, I learnt several other interesting things too.

For example, I started the day learning about some proposals around a Multi Cluster Ingress - Ingress v2 basically *(yes please, looked great!)*, an intro to SPIFFE and SPIRE *(which sounded useful but a bit early-days for folks like us, perhaps)*, some wildly sensible folks from Lyft & Triller talking about permissions options in Kubernetes *(reassuringly, we are doing a lot of things in this space already - I did learn some funky things about Admission Controllers and Open Policy Agent though)*.

The latter inspired me sufficiently to go to a more thorough talk on [Open Policy Agent Gatekeeper](https://github.com/open-policy-agent/gatekeeper) at the expense of some other interesting stuff - presented by folks from Google & Microsoft. It works through a set of CRDs so feels more config-y than code-y - although there is code buried in the templates you produce. This sounds like something we should get into - it's only in Alpha at the mo. They worked through a demo of four policy enforcement use-cases:

- all namespaces must have a label that lists a point of contact
- all pods must have an upper bound for resource usage
- all images must be from an approved repository
- Services must all have globally unique selectors

These sort of things sounds highly relevant to me. I know the team have tried it out a bit and been put off by the Rego language, but I understand that the tooling support for this is improving considerably, and as the popularity of OPA increases, examples of how to do this well will make our lives easier in this area anyway. They're definitely shooting for making it easier to write, test, automate into pipelines, and re-use - all good things I like to hear!

> Most importantly, I learnt that we should be pronouncing it "Oh-pa" not "Oh-Pea-Ay" ... +1 for more pronunciation guides. Now that's a valuable reflection right there.

Bring on Day 2 ...

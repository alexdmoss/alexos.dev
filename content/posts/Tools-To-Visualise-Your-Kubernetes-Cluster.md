---
title: "Tools To Visualise Your Kubernetes Cluster"
panelTitle: "Tools To Visualise Your Kubernetes Cluster"
date: 2021-03-28T18:00:00-01:00
author: "@alexdmoss"
description: "My experience running a couple of open source tools to help visualise workloads deployed on Kubernetes"
banner: "/images/???"
tags: [ "Kubernetes", "Weave", "Scope", "Ops View", "kube-ops-view", "Graph", "Observability" ]
categories: [ "Kubernetes", "Observability" ]
---

{{< figure src="/images/???.jpg?width=800px&classes=shadow" attr="Photo by ??? on Unsplash" attrlink="???" >}}

In this post I'm going to discuss a couple of tools I've used to try to visualise workloads deployed on my Kubernetes Clusters. The tools I'll be looking at are:

1. Kube Ops View
2. WeaveWorks Scope

I'll also finish by discussing some of the other approaches I am yet to try fully, but which might be better alternatives for larger clusters.

---

## Background

I have been on-off interested in tools to help visualise all the things running on a Kubernetes cluster since I first got involved in Kubernetes. A bash script that hacks together a bit of `kubectl` to dump out the pods on each node (later replaced by the built-in tools in the Google Cloud Console, since I use GKE) was an early creation that did a good enough job of satisfying my curiosity, but this is naturally hard to read.

My initial interest was motivated by trying to better understand how Kubernetes' scheduler makes decisions on where to put workloads, and potentially to try to get better bin packing (I've [written about this before](https://alexos.dev/2019/09/28/squeezing-gke-system-resources-in-small-clusters/)). That interest later evolved into wanting to understand the interconnectedness of workloads running on the cluster - and I'm going to describe two tools I've used to attempt to do this below.

/// medium - invite comments on alternatives here

## So Just How Loaded Are My Nodes, Anyway?


## I Need More. How Does It All Fit Together?

## The Nirvana?

/// check right usage

/// node graphing in grafana?

/// mesh

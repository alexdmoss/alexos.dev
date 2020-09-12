---
title: "Simplifying switching between kubectl contexts on GKE"
panelTitle: "Simplifying switching between kubectl contexts on GKE"
date: 2020-09-11T18:00:00-01:00
author: "@alexdmoss"
description: "I've changed the way I switch between GKE clusters in my shell, and thought I'd share the approach I've taken to deal with multiple clusters across multiple projects with multiple user accounts"
banner: "/images/???"
draft: true
tags: [ "GKE", "Kubernetes", "Shell", "kubectl", "GCP" ]
categories: [ "GKE", "Kubernetes", "Tips n Tricks" ]
---

Something something: for later:

Useful references:

https://blog.sbstp.ca/introducing-kubie/
https://medium.com/google-cloud/kubernetes-engine-kubectl-config-b6270d2b656c

Future - https://medium.com/google-cloud/authing-w-kubernetes-engine-service-accounts-ae752b46ed18
- genuinely need separate accounts

Notes from my readme:

For brand new setup, you will need to `gcloud init` once to set up the default configuration. It will also be necessary to `gcloud auth login` on each account used at least once (and there is expiry on their validity so this may come up again if not refreshed for long enough, I suspect).

With that done, the `switch` script should allow moving between different GCP projects. See the preconfigured configs at the top of that script (`my-bin/switch`).

For setting up kubectl contexts - not automated (can't be bothered dealing with the secrets yet) - `kubie` must be installed first. For each cluster, you can then:

1. Use `get-creds` to authenticate **in a shell that is not using kubie**. This will create a new `~/.kube/config`.
2. Make a copy of an existing config file in `~/.kube/configs/`.
3. Copy the new `cluster.certificate-authority-data` and `cluster.server` into the cloned config file.
4. Change the `clusters.name`, `contexts.name` and `contexts.context.cluster` to match your desired name. Also update the `namespace` if you wish.
5. Delete `~/.kube/config` to avoid a duplicate.

For the very first cluster to set up, the `users.user.auth-provider.config.access-token` will need to be picked out too (see below).

For that first time setup, here's a blueprint for a file to go into `~/.kube/configs/*.yaml`:

```yaml
clusters:
  - cluster:
      certificate-authority-data: [redacted]
      server: [redacted]
    name: my-cluster
contexts:
  - context:
      cluster: my-cluster
      namespace: default
      user: gcloud-account
    name: my-cluster
users:
  - name: gcloud-account
    user:
      auth-provider:
        config:
          access-token: [redacted]
          cmd-args: config config-helper --format=json
          cmd-path: /usr/lib/google-cloud-sdk/bin/gcloud
          expiry: "2020-09-03T09:35:15Z"
          expiry-key: '{.credential.token_expiry}'
          token-key: '{.credential.access_token}'
        name: gcp
```


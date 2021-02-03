---
title: "How to connect to multiple Google Kubernetes clusters easily in parallel"
panelTitle: "How to connect to multiple Google Kubernetes clusters easily in parallel"
date: 2021-02-02T18:00:00-01:00
author: "@alexdmoss"
description: "My approach to connecting to multiple clusters across multiple accounts from my terminal"
banner: "/images/kube.jpg"
tags: [ "GKE", "Kubernetes", "Shell", "kubectl", "GCP" ]
categories: [ "GKE", "Kubernetes", "Tips n Tricks" ]
---

{{< figure src="/images/kube.jpg?width=800px&classes=shadow" attr="Photo by Alvaro Reyes on Unsplash" attrlink="https://unsplash.com/photos/4eTnTQle0Ks" >}}

In this post I'm going to talk through the approach I use to switch between multiple Google Kubernetes Engine clusters on the command line. I'd expect a lot of the stuff in here has some benefit for non-GCP Kubernetes clusters too, but the ones I use on a day-to-day basis are all hosted there.

Key outcomes for me were:

1. To be able to switch from one cluster to another with only a small number of commands, even if I have to authenticate with different user accounts (as primarily a GCP user, this means different email addresses / GCP projects).
2. To be connected to different Kubernetes in different windows - for example tailing logs in a Prod & Non-Prod cluster at the same time.

**Spoiler Alert:** I solved this with a bit of fiddling of `~/.kube/config` plus the marvellous [kubie](https://blog.sbstp.ca/introducing-kubie/).

---

## Background

I switched to this approach probably around six months ago. At work we have a relatively small number of clusters so up until that point I was pretty comfortable using what I think is the most common approach of [kubectx + kubens](https://github.com/ahmetb/kubectx) and this worked well enough. However I found that I was increasingly getting an inconsistent experience when switching between the clusters I use for work and others I was using for running personal websites _(like this one!)_ and for fiddling with things - so I started looking for alternatives.

For the purposes of explaining things, lets assume a hypothetical setup like this:

- A collection of GKE clusters across several GCP projects at work. Access to these is through a common Production email account, but they are split by project/cluster.
- A sandbox GKE cluster in a different Google organisation where the Prod email account doesn't have access.

The hierarchy would therefore look a bit like this:

```yaml
- user:
  - account: alex@prod.work
    gcp-projects:
    - gcp-project:
      - name: prod
        cluster: brie
    - gcp-project:
      - name: staging
        cluster: cheddar
- user:
  - account: alex@dev.work
  - gcp-projects:
    - gcp-project:
      - name: sandbox
        cluster: chutney
```

> I feel I should point out here that, at work, we do not name our clusters after types of cheese. I just really fancied some cheese when writing this, ok?

Enough background - on with how I set things up.

---

## First, Multiple Google Accounts

{{< figure src="/images/monitors.jpg?width=800px&classes=shadow" attr="Not actually me - my desk is not this tidy. Photo by Max Duzij on Unsplash" attrlink="https://unsplash.com/photos/qAjJk-un3BI" >}}

My approach here was massively inspired by [this blog post](https://medium.com/google-cloud/kubernetes-engine-kubectl-config-b6270d2b656c) by Googler Daz Wilkin. I'm not going to repeat what is explained really well there already so go have a read if you want to understand why the following works!

> For brand new setup, you will need to `gcloud init` once to set up the default configuration. It will also be necessary to `gcloud auth login` on each account used at least once, and this may need refreshing once in a while (but not often enough for me to really notice).

I threw away my pre-saved `gcloud` configurations - not gonna need 'em! All I have in `~/.config/gcloud/` is a `config_default`, which gets updated with a simple bash script when I need to switch between Google Accounts/Projects.

The switching script looks [like this](https://gist.github.com/alexdmoss/8ac9eea1f55063e2c99f8b7bbe9564cd) - with the bits wrapped in `<< >>` to be replaced. I alias this so I would simply do `switch dev` for a pre-saved project, or `switch gcp-project user-email` to activate a new/rarely-used project.

{{< gist alexdmoss 8ac9eea1f55063e2c99f8b7bbe9564cd >}}

---

## Second, Multiple Kubernetes Contexts

{{< figure src="/images/tracks.jpg?width=800px&classes=shadow" attr="Photo by Radek Grzybowski on Unsplash" attrlink="https://unsplash.com/photos/KVenyQf7gH0" >}}

Here, I keep `~/.kube/config` completely clear and set up new configs under `~/.kube/configs/` whenever I have a new cluster I need to deal with. For the number I have to worry about, this is quite manageable, but it could be automated if you had a frequently-changing enough list to be worthwhile. The steps look like this (this must be done in a brand new shell not using `kubie` - see below!):

1. Auth to the new cluster as normal: `gcloud container clusters get-credentials ${cluster} --project=${project} --zone=${zone}`. This adds an entry to your blank `~/.kube/config`.
2. Copy a template config file (see below) into `~/.kube/configs/` with a unique name.
3. Take the values for `clusters.cluster.certificate-authority-data` and `clusters.cluster.server` from the no-longer blank `~/.kube/config` and put them into your new file created from the template.
4. Update the `name:` fields for the cluster to reflect what you want it to be known as when you list your contexts - `clusters.name`, `contexts.name` and `contexts.context.cluster`. It does not have to match exactly the cluster name if you want to save typing.
5. Delete `~/.kube/config` (unless you want to have a default cluster for when not using `kubie` - but you'll need to keep this file tidy to avoid confusion!).

The template I mentioned for this looks as follows:

```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: <<<your-cert-goes-here>>> # taken from .kube/config
    server: <<<https://(your-apiserver-ip)>>>             # taken from .kube/config
  name: <<<name-of-context>>>                             # your choice of name
contexts:
- context:
    cluster: <<<name-of-context>>>                        # your choice of name
    user: gcloud-account
  name: <<<name-of-context>>>                             # your choice of name
kind: Config
preferences: {}
users:
- name: gcloud-account
  user:
    auth-provider:
      config:
        cmd-args: config config-helper --format=json
        cmd-path: /usr/lib/google-cloud-sdk/bin/gcloud    # your path may vary
        expiry-key: '{.credential.token_expiry}'
        token-key: '{.credential.access_token}'
      name: gcp
```

While it looks like a few things to change, in practice it only takes a few seconds - and only needs doing for a fresh cluster. Quite ok really.

Using my list of clusters from earlier, I would have `brie.yaml`, `cheddar.yaml` and `chutney.yaml` all in my `.kube/configs/` directory, with valid certs/server details, but no mention of any GCP account details (we will just use the current gcloud config at the time we connect to them, thanks to the `switch` script).

---

## Finally, Loading Parallel Kubernetes Contexts

{{< figure src="/images/red-dice.jpg?width=800px&classes=shadow" attr="Photo by Jonathan Petersson on Unsplash" attrlink="https://unsplash.com/photos/W8V3G-Nk8FE" >}}

To make use of this shiny new config, we bring in [kubie](https://blog.sbstp.ca/introducing-kubie/). This tool works in a similar way to `kubectx` + `kubens` - you specify `kubie ctx` to set your current cluster, and `kubie ns` to select a namespace. The difference being, that when you run `kubie ctx` you spawn a new shell within your terminal window, with the context loaded to that.

What that means in practice is you can have a terminal on the left of your screen connect to e.g. `prod` and a terminal on the right of your screen connected to e.g. `dev`, and both continue to work independently from each other. This is really marvellous.

> I have sufficient muscle memory that I had to `alias kctx='kubie ctx'` and `alias kns='kubie ns'` to save re-learning / more typing

There's also a `kubie exec` to run just one command using a different context without swapping out the whole shell if you prefer - for example `kubie exec cheddar kube-system kubectl get pods`. This is really handy if you want to use this in scripts across multiple clusters.

There's way more info/options available - see [the project on github](https://github.com/sbstp/kubie) for more ideas.

---

## How This Works in Practice

If working with two clusters and a shared user account, then I simply issue `kctx brie` and `kctx cheddar` in separate terminals and I'm away.

If the second cluster needs a separate user account, then I would `switch sandbox alex@dev.work` first, then `kctx chutney`, and I'm sorted. The only thing I need to keep in mind here is that my **gcloud** context has switched globally (no equivalent of `kubie` here), so any gcloud SDK commands are going to be against `sandbox` in both terminal windows (unless I switch again) - but my `kubectl` commands are fine (I suspect until the refresh token expires, but in practice I've never had an issue).

To show this working in practice:

{{< webm src="/casts/kubie-config.webm" width="90%" >}}

> Before you get any ideas, the Cheddar cluster is long-since deleted - that certificate in the video is useless :smile:

And here's an example with two different clusters being watched at the same time:

{{< webm src="/casts/kubie-two-shells.webm" width="90%" >}}

And that's a wrap - hopefully this inspires you to give kubie a try!
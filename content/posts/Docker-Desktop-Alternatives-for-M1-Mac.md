---
title: "Docker Desktop Alternatives for M1 Mac"
panelTitle: "Docker Desktop Alternatives for M1 Mac"
date: 2022-01-02T19:21:00-01:00
author: "@alexdmoss"
description: "A few options, including my preference, for replacing Docker Desktop on newer M1 Macs running Apple Silicon"
banner: "/images/containers.jpg"
tags: [ "Docker", "Docker Desktop", "Minikube", "Podman", "Lima", "nerdctl", "containerd", "CLI", "Guide", "Mac", "M1", "Apple Silicon", "WSL" ]
---

{{< figure src="/images/containers.jpg?width=1000px&classes=shadow" attr="Photo by Ian Taylor on Unsplash" attrlink="https://unsplash.com/photos/jOqJbvo1P9g" >}}

In this blog post I'm going to talk through my recent experiences as I attempted to ditch Docker Desktop - the licensing changes that come into effect at the end of January being the primary motivator.

> Without going into any detail about it, let's just say I'm not a fan of taking something that you've made freely available previously and deciding that you now want to charge for it!

{{< figure src="/images/docker-money.png?width=400px&classes=shadow" attr="From this 'Is Docker in Trouble?' Blog Post" attrlink="https://start.jcolemorrison.com/is-docker-in-trouble/" >}}

This article will mostly focus on MacOS, although there is a brief note about Windows/WSL included for completeness too. I tried out ~~three~~ four options for Mac - landing on one as my preference as it covered both the need to run on the newer Apple Silicon and allow mounting of volumes on the host OS, which is something I do fairly frequently (mostly to shorten the feedback loop when testing changes that run on an image intended to run in CI).

**Disclaimer:** Most of the steps detailed below were found through following other fantastic blog posts I found out there :clap:. These are of course noted wherever I've used them, with a few tweaks of my own I've made on top of these excellent guides. Hopefully having these options together in one blog post is somewhat helpful in choosing between them too!

---

## TL:DR

**Updated 14/01/2022:** Rancher Desktop ([rancherdesktop.io](https://rancherdesktop.io/)) provide a package that has a similar feel to Docker Desktop - download and install it, let the wizard do its thing and you're up and running within a few minutes. This has recently added Apple Silicon support and gets my vote as it provides both this + the `docker` CLI + volume host mounting all in one easy to install package. Local Kubernetes support is also a minor plus for me.

If for whatever reason this does not work for you, then I've left the other options I've tried out previously below, as well as a little more detailed instructions and some of the thought process that went into this also.

There's also a brief nod at the start to Windows + WSL, which I use very occasionally. This is easy to setup without Docker Desktop.

---

## Docker with WSL

{{< figure src="/images/wsl-funny.jpg?width=400px&classes=shadow" attr="(Image Source)" attrlink="https://www.reddit.com/r/linuxmemes/comments/niur6p/wsl_meme/" >}}

**A brief aside** - I occasionally use Windows 10 with WSL v2 installed too :scream: _(sidebar: it actually works pretty well to be honest!)_. I was pretty confident that this worked without Docker Desktop. This turned out to be 100% true - and you can manage perfectly fine without it as long as you're running WSL version 2. I won't break down the detailed steps to set this up - although if you'd like to me to, get in touch via the [Contact](/contact) option and I'd be happy to. It boils down to:

1. Ensure you are using WSL version 2
2. Install docker as you normally would in your Linux distro (I use Ubuntu, and had no problems). See [this gist](https://gist.github.com/alexdmoss/14fd33bb91beb2283d6a8ed0a8df64a3) for how I've done it
3. You need to start the docker daemon by hand (e.g. `sudo dockerd > /tmp/dockerd.log 2>&1 &`), as WSL has its own startup routines (that Docker Desktop was handling for us)

Note that Rancher Desktop also works perfectly fine on Windows 10 too it seems - it takes care of the WSL installation for you behind the scenes if you don't want to roll your sleeves up and get into that sort of thing.

Okay, enough of that Windows nonsense - onto the MacOS stuff now ...

---

## Option 1: Docker + Hyperkit + Minikube

{{< figure src="/images/cube.jpg?width=600px&classes=shadow" attr="Photo by Rostislav Uzunov on Pixabay" attrlink="https://pixabay.com/illustrations/blue-crystal-cube-deep-futuristic-5457731/" >}}

I started here. This is the most "drop-in" replacement in the list, but **does not work on M1 Macs**. I used this on my older Macbook for a little while before replacing it with Rancher Desktop. It's fully docker compliant, if there is such a thing.

The instructions that follow are heavily based on [this excellent blog post](https://itnext.io/goodbye-docker-desktop-hello-minikube-3649f2a1c469), which has some additional advice, especially if you want to get more out of the local Kubernetes cluster:

```bash
#
# pre-req: full install of XCode needed - just the CLI isn't enough
#
brew install hyperkit                               # this fails on Apple Silicon: https://github.com/moby/hyperkit/issues/310
brew install docker                                 # don't use --cask - that's Docker Desktop!
brew install minikube

minikube start --driver=hyperkit --keep-context     # this is where it errors on Apple Silicon

eval $(minikube docker-env)                         # tells docker CLI in your *current shell* to use minikube's docker daemon
```

I have the `minikube start` command set up in a `start-docker.sh` script I can run when needed, and the `minikube docker-env` in my shell startup (`.zshrc`, in my case).

As you can see, pretty straight-forward standard brew installation stuff - plus a couple of commands to run before you try and do docker things _(I personally never had Docker running all the time on startup anyway, as it was such a battery drain)_. As it's still just the same docker CLI, the credentials helper to connect to a private registry also works fine out-the-box.

However, volume mounts from the host did not ... but thankfully the blog post I linked above has captured the solution for this. You can `minikube mount` to spin up a process to mount your local directory into the minikube VM:

```bash
minikube mount your-local-directory/:/build >/dev/null 2>&1 &    # obvs do not send to dev/null if debugging it!

docker run -v /build:/build --rm -it eu.gcr.io/my-private-registry-project/alex-ubuntu:latest
```

In my opinion, the advantages of this option - and why I kept it as the setup on my older Macbook - are:

- You're still using the docker CLI, so very good from a compatibility point of view
- There's minimal extra config/scripts needed - almost a drop-in replacement
- It's great if you want to do Kubernetes development locally and liked that feature in Docker Desktop

Downsides:

- Doesn't currently support Apple Silicon (or at least, using the hyperkit driver does not)
- Volume mounts aren't seamless (although pretty simple tbh)

However, I also needed an option that worked with Apple Silicon. My first attempt was with **Podman** ...

---

## Option 2: Podman

{{< figure src="/images/podman.jpg?width=600px&classes=shadow" attr="Loved this image - found it here ..." attrlink="https://ios.dz/installation-podman-centos-8/" >}}

After realising that hyperkit didn't work on M1, this was the next option I tried. I'd heard good things. It mostly worked fine but, as mentioned earlier, for me the crucial issue was the **lack of ability to mount volumes from the host OS**. I use this option a lot.

That said, if that's not important to you or they fix it subsequently, I've included the steps below. these were cobbled together from the Podman installation guide itself plus this excellent [blog post](https://marcusnoble.co.uk/2021-09-01-migrating-from-docker-to-podman/) - although I didn't need most of the complexity involved here, as I'm guessing it has been fixed since.

From their install guide:

```sh
# a nice simple install + setup ...
brew install podman
podman machine init
podman machine start            # suspect only this needs to go into your docker startup script
podman info                     # just to confirm things are ok
```

Your docker equivalents should then work as intended:

```bash
podman build -t podman-test -f Dockerfile .         # builds a Dockerfile containing a basic nginx image
podman run -d -p 8080:80 podman-test                # run it, exposing port
curl http://localhost:8080/                         # see the appropriate output from nginx
```

Other similar `docker` commands I tend to use also seem present:

```bash
podman ps
podman stop <id>
podman exec -i -t <id> /bin/bash                    # a subtlety - requires -i -t rather than allowing -it
```

However, as mentioned earlier this crucially does not work:

```bash
podman run --rm -it -p 8080:80 -v $(pwd):/build podman-test         # /build is empty :disappointed:
```

I looked around this topic a bit and there are some suggested workarounds, such as [this one](https://github.com/containers/podman/issues/8016#issuecomment-939353204) to mount the directory onto the podman VM first. But these look quite hasslesome  _(caveat: I didn't try very hard :wink:)_

I therefore backed away at this point as I had another option to try ... Enter `lima` + `nerdctl` ...

---

## Option 3: Lima + nerdctl

{{< figure src="/images/nerdy.png?width=600px&classes=shadow" attr="N-Er-Dy" attrlink="https://www.pinterest.co.uk/pin/594897432010361218/" >}}

This option **ticked all the boxes for me** and I ran with it for a little while without issue, although with more setup needed than the minikube option. I'm comfortable with that though. I like this because it: a) distances me from Docker Inc. changes to licensing in the future (and a little bit on principle, not gonna lie!), and b) puts me closer to the OCI runtime of our Production Kubernetes clusters (they're GKE, which just run containerd by default now).

I followed the great guide [in this blog post](https://medium.com/nttlabs/containerd-and-lima-39e0b64d2a59), which basically boils down to:

```bash
brew install lima
limactl start default                               # accepted the defaults to setup the VM
```

Your docker equivalents then look like this (which can of course be aliased):

```bash
lima nerdctl build -t lima-test -f Dockerfile .     # Dockerfile containing a basic nginx image
lima nerdctl run -d -p 8080:80 lima-test
curl http://localhost:8080/                         # see the appropriate output from nginx
```

Other similar `docker` commands also seem fine, just like `podman`:

```bash
lima nerdctl ps
lima nerdctl stop <id>
lima nerdctl exec -it <id> /bin/bash
```

Crucially, this worked too without any special config or setup needed:

```bash
lima nerdctl run --rm -it -p 8080:80 -v $(pwd):/build --entrypoint=/bin/bash lima-test
```

That said, there were a couple of other steps I needed to go through to deal with my other requirements.

Because it's not docker, the existing credentials helpers I had setup to connect to e.g. Google Container Registry did not automatically work. Instead, these credentials need to be readily available on the intermediary lima VM, rather than the host. To solve this, I opted to jump onto the VM and install `gcloud`, login as I normally would, then ensure those credentials were available to the `root` user.

To do this, we start with `limactl shell default` which should get you a shell prompt on your default lima VM. We then download and unpack the GCloud SDK:

```bash
# update with your preferred gcloud version
wget --no-verbose -O /tmp/google-cloud-sdk.tar.gz \
    https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-367.0.0-linux-x86_64.tar.gz && \
    sudo tar -C /opt --keep-old-files -xz -f /tmp/google-cloud-sdk.tar.gz && \
    sudo chown $(whoami):$(whoami) /opt && \
    sudo chown -R $(whoami):$(whoami) /opt/google-cloud-sdk && \
    rm -f /tmp/google-cloud-sdk.tar.gz
```

As root we then symlink to gcloud + the credential helper, and login:

```bash
sudo ln -s /opt/google-cloud-sdk/bin/gcloud /usr/bin/gcloud
sudo ln -s /opt/google-cloud-sdk/bin/docker-credential-gcloud /usr/bin/docker-credential-gcloud
gcloud auth login                # login as normal
gcloud auth configure-docker     # will warn about docker path, can ignore
```

Unfortunately, we're not quite there yet - but nearly! Back on the host machine, spinning up my Ubuntu docker image was met with an error message I've seen a few times before on Apple Silicon: `standard_init_linux.go:228: exec user process caused: exec format error`. We need to do a bit of work to give QEMU (the hypervisor behind the scenes) the option to execute non-native images.

Thankfully the [`nerdctl` docs](https://github.com/containerd/nerdctl/blob/master/docs/multi-platform.md) point you in the right direction on this one, via [this super-useful emulator](https://github.com/tonistiigi/binfmt). We therefore do the following:

```bash
limactl shell default                       # onto the VM again
sudo systemctl start containerd
sudo nerdctl run --privileged --rm tonistiigi/binfmt --install all
ls -1 /proc/sys/fs/binfmt_misc/qemu*        # this is just to check it worked - should list several extra chipsets 
```

... et voila! Our ubuntu image built on amd64 in a private container registry with a local host volume mount works without issue :tada::

```bash
lima nerdctl run -v $(pwd):/build --rm -it eu.gcr.io/my-private-gcr-project/alex-ubuntu:latest
```

> **A small note:** If you need that local directory to be writable by the container, you need to edit the file `~/.lima/default/lima.yaml` on your host. There's a `mounts:` section where you can choose to make your home directory and everything in it writable (dodgy if you run untrusted containers!), or add a block to a separate mount path, similar to the one already there for /tmp (the option I took!).

All that's left is to add the `limactl start default` to your startup script and alias `lima nerdctl` to something - you can even alias this to `docker` if you wish (although I personally prefer to be more explicit - I chose to alias it to `nerd` :metal:).

In my opinion, the advantages of this option are:

- It works on Apple Silicon! :tada:
- It handles volume mounts seamlessly
- It distances you from docker itself and any further licensing fun and games down the line
- For me, it's closer to my Production container stack

The downsides:

- It's not actually docker - so there's a risk of hitting compatibility issues in comparison to e.g. how things are behaving in CI, perhaps
- It's a little bit of a faff to get working with a private container registry in particular
- It _seems_ to take a bit longer to startup than Minikube / Docker Desktop. This is pretty anecdotal though, and not particularly impactful to me day-to-day

---

## **Recommended Option** - Rancher Desktop

{{< figure src="/images/rancher.jpg?width=800px&classes=shadow" attr="Photo by Kendall Ruth on Unsplash" attrlink="https://unsplash.com/photos/8CTz62aHidw" >}}

**Updated 14/01/2022:** So we're done right? Lima + nerdctl does the trick? As mentioned back at the top - whilst that was my preferred option for a while, I've recently discovered Rancher Desktop, and it's latest version (0.7+) introduced Apple Silicon support.
This is the option I'm running with now on both my MacOS machines.

Why this option over the rest?

- Its installation process is very simple. _I love a good bit of command line fiddling, but if it's not necessary and the software keeps itself patched, then it doesn't have to be._
- It works with all my device types - M1 Mac, Intel Mac, Windows 10.
- It's fully compatible with the `docker` CLI. This is good news for existing build scripts (e.g. if you need to run these or your tests locally, and they depend on being able to call docker directly - pretty common, and avoids surprises if you aliased `nerdctl`).
- Volume mounts to the host OS work without any special configuration or adjustments to scripts / compose files.
- It has a local Kubernetes environment, for the odd occasion I find that handy to have.

The installation process is extremely simple with only a few choices to make (it does most of it behind the scenes - and makes use of lima to do it on Mac). You will be prompted to elevate your access a few times, which is understandable given what you're setting up here.

I opted to use the recommended stable Kubernetes release and the `dockerd/moby` engine - but I really like that it offered me the choice of `containerd (nerdctl)` in case that becomes handy in the future (and I understand from the docs that they can coexist too).

I did hit a couple of small issues that were local to my device, and probably a result of my various experiments on this topic! Here's a couple of quick bits of **troubleshooting advice** if you need it:

- As a fairly new product, its user experience when things don't work is not stellar. Logs can be found on Mac here: `~/Library/Logs/rancher-desktop/` and provided me with the clues I needed.
- It is particular about permissions. For example, the `/opt` directory on my Mac was not owned by `root` - the logs above clued me in to this need, as I was just faced with an `Error Starting Kubernetes - limactl exited with code 1` error through the UI.
- You must run it from the `/Applications` on your Mac. Copy it over if you didn't pick the .dmg install option.
- You must use versions >0.7.0 for it to work on Apple Silicon.
- _Aside: On Windows it needs to reboot a few times. Just roll with it_ 😏

## Summary

So there you have it - **Rancher Desktop** ([rancherdesktop.io](https://rancherdesktop.io/)) was my preferred option for replacing Docker Desktop on any MacOS machine. Hopefully you found this run through of the steps useful for your particular setup. As always with these things - and indeed in my own experience following the existing advice out there - it may not work flawlessly on your kit.

{{< figure src="/images/works-on-my-machine.jpg?width=400px&classes=shadow" attr="Well it worked on my machine" >}}

If you find any issues, do let me know via the [Contact](/contact) page - I'd be interested to keep this post up to date with any additional advice over time too!

---

This post was also published to [my Medium account here](https://alex-moss.medium.com/docker-desktop-alternatives-for-m1-mac-918a2dcda10).

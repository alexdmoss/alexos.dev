---
title: "Creating an up-to-date Distroless Python Container Image"
panelTitle: "Creating an up-to-date Distroless Python Container Image"
date: 2022-06-26T13:21:00-01:00
author: "@alexdmoss"
description: "How I created an up-to-date Distroless Container Image for running Python applications in Kubernetes, without relying on the out-dated one published by Google"
banner: "/images/distroless.jpg"
tags: [ "Python", "Distroless", "Containers", "Security", "Docker", "Poetry", "Pipenv", "Flask", "Gunicorn", "Fast API", "Pandas" ]
draft: true
---

// TODO: image for distroless

{{< figure src="/images/???.jpg?width=1000px&classes=shadow" attr="Photo by ??? on Unsplash" attrlink="https://unsplash.com/photos/???" >}}

In this blog post I walk through the creation of a Python Docker image based on the Distroless container base, but with an up-to-date version of Python and operating system updates - unlike the experimental version [published by Google](https://github.com/GoogleContainerTools/distroless). This image still has the same security and operational benefits - no shell or unnecessary OS libraries to reduce the security attack surface, and also a tiny image size.

Don't need the background and just want to see how I built it - [jump to that section](#building-an-alternative).

---

## What's Distroless?

// TODO: image, less is more?

A little bit like the term "serverless", the term Distroless (in my opinion anyway!) is a trendy misnomer. The Linux distribution is still there - what we really mean here is a container image that contains as little of an Operating System as possible - just enough to run your application. In particular, there's no shell.

Why do we care about this? Well, it's more secure. Or more specifically:

1. There are fewer OS libraries and tools available to be exploited by a baddie to gain access to your runtime - a smaller attack surface.
2. In the event that they did (e.g. through a code vulnerability instead), they are much more limited in the harm they can do, as there's no shell to execute their remote commands against, and fewer exploitable tools on the host.

Whilst distributions like Alpine Linux are excellent at helping with this too, you still have a shell to (ab)use. If you're interested in more detail on this, [here's a great article by a former colleague on exact this!](https://www.equalexperts.com/blog/tech-focus/docker-security-battle-of-the-base-image/).

---

## The Python Predicament

// TODO: a python

This choice is made even more complex with Python. Whilst common languages like [Java](https://github.com/GoogleContainerTools/distroless/blob/main/java/README.md) and [Node](https://github.com/GoogleContainerTools/distroless/blob/main/nodejs/README.md) have well-established Distroless variants that are released frequently, the [Python one](https://github.com/GoogleContainerTools/distroless/tree/main/experimental/python3) continues to be marked as experimental and changes rarely. In practice this means it ships with whatever the Debian upstream version of Python and its dependencies are, leading to [issues like this one](https://github.com/GoogleContainerTools/distroless/issues/1003). We can see this by running a vulnerability scan using a tool like trivy, showing that the CVE mentioned there is still present 3+ months later:

```sh
> trivy image -s=HIGH,CRITICAL gcr.io/distroless/python3:latest

# [...]
2022-06-26T13:49:19.271+0100	INFO	Number of language-specific files: 0

gcr.io/distroless/python3:latest (debian 11.3)

Total: 10 (HIGH: 7, CRITICAL: 3)
```

You might think then that perhaps with Python specifically it'd be better to use Alpine. Indeed that's something I've done a fair few times myself. But in practice, whilst this does help in many ways with security challenge, this starts to become problematic for other reasons. My understanding here is that this all stems from Alpine's use of `musl` rather than `glibc` as its standard C library. The most noticable impact here is creating **horrendous** build times when certain (common) packages are involved, as there are no pip wheels available (and sometimes out-right failing to build at all).

> Personally, I find gRPC (`grpcio`) is a common example of this, as someone who does a lot of work with Google Cloud

This isn't the only risk you take by using Python on Alpine. As [this excellent blog post](https://pythonspeed.com/articles/alpine-docker-python/) elaborates on, this - the TL:DR that as well as the build time problem, you can encounter sometimes very subtle bugs and performance issues as a result of the choice too.

So, warned off that, we're left with Debian variants (e.g. `python:*-slim-bullseye`). I'd imagine these are by far the most common in the wild - but they too suffer from the presence of OS library vulnerabilities in the base image (Debian 11), e.g.

```sh
> trivy image -s=HIGH,CRITICAL python:3.9-slim-bullseye

# [...]
2022-06-26T13:59:34.109+0100	INFO	Number of language-specific files: 1

python:3.9-slim-bullseye (debian 11.3)

Total: 16 (HIGH: 13, CRITICAL: 3)
```

A slightly higher - but different set - of vulnerabilities being reported. This reflects that on the one hand this image has been built more recently and picked up a number of fixes, but also has more "stuff" in it than the distroless GoogleContainerTools python image scanned above.

So, only one thing for it really! Lets see if we can build our own version that is a happy story from a vulnerability perspective AND still works ok.

---

## Building an Alternative

// TODO: arms rolled up image

---

## Testing The Thing

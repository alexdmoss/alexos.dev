---
title: "Creating an up-to-date Distroless Python Image"
panelTitle: "Creating an up-to-date Distroless Python Image"
date: 2022-07-08T13:21:00-01:00
author: "@alexdmoss"
description: "How I created a Distroless Container Image for Python, without relying on the out-dated Google version"
banner: "/images/distroless.jpg"
tags: [ "Python", "Distroless", "Containers", "Security", "Docker", "Poetry", "Pipenv", "Flask", "Gunicorn", "Fast API", "Pandas" ]
---

{{< figure src="/images/distroless.jpg?width=1000px&classes=shadow" attr="Photo by Prateek Katyal on Unsplash" attrlink="https://unsplash.com/photos/MGRv3qZfyTs" >}}

In this blog post I'll walk through the creation of a Python Docker image based on the Distroless container [published by Google](https://github.com/GoogleContainerTools/distroless), but with an up-to-date version of Python and operating system updates - unlike their [experimental](https://github.com/GoogleContainerTools/distroless/tree/main/experimental/python3) (and unsupported) version. This image still has the same security and operational benefits - such as no shell or unnecessary OS libraries to reduce the security attack surface, and preserving the tiny image size.

Don't need the background and just want to see how I built it? [Jump to that section](#building-an-alternative). The source code for all of what follows [can be found on Github too](https://github.com/alexdmoss/distroless-python).

---

## What's Distroless?

{{< figure src="/images/less.jpg?width=1000px&classes=shadow" attr="Photo by Etienne Girardet on Unsplash" attrlink="https://unsplash.com/photos/fti002hQCCA" >}}

A little bit like the term "serverless", the term Distroless (in my opinion anyway!) is a trendy misnomer. The Linux distribution is still there - what we really mean here is a container image that contains as little of an Operating System as possible - just enough to run your application. In particular, there's no shell.

Why do we care about this? Well, it's "more secure". Or, more specifically:

1. There are fewer OS libraries and tools available to be exploited by a baddie to gain access to your runtime - it has a smaller attack surface.
2. In the event that they do (e.g. through a code vulnerability instead), they are much more limited in the harm they can do, as there's no shell to execute more dangerous commands through, and fewer exploitable tools on the host.

Whilst distributions like Alpine Linux are excellent at helping with this too, you still have a shell to (ab)use. If you're interested in more detail, [here's a great article by one of my former colleagues on this](https://www.equalexperts.com/blog/tech-focus/docker-security-battle-of-the-base-image/).

---

## The Python Predicament

{{< figure src="/images/python.jpg?width=1000px&classes=shadow" attr="Photo by David Clode on Unsplash" attrlink="https://unsplash.com/photos/vb-3qEe3rg8" >}}

This base image choice is made even more complex with Python. Whilst common languages like [Java](https://github.com/GoogleContainerTools/distroless/blob/main/java/README.md) and [Node](https://github.com/GoogleContainerTools/distroless/blob/main/nodejs/README.md) have well-established Distroless variants that are updated frequently, the [Python one](https://github.com/GoogleContainerTools/distroless/tree/main/experimental/python3) continues to be marked as experimental and appears to change more rarely. In practice this means it ships with whatever the Debian upstream version of Python and its dependencies are, leading to [issues like this one](https://github.com/GoogleContainerTools/distroless/issues/1003). We can see this by running a vulnerability scan using a tool like [trivy](https://github.com/aquasecurity/trivy), showing that the CVE mentioned there is still present 3+ months later:

```sh
> trivy image -s=HIGH,CRITICAL gcr.io/distroless/python3:latest

# [...]
2022-06-26T13:49:19.271+0100	INFO	Detecting Debian vulnerabilities...

gcr.io/distroless/python3:latest (debian 11.3)

Total: 10 (HIGH: 7, CRITICAL: 3)
```

You might think then that perhaps with Python specifically it'd be better to use Alpine. Indeed that's something I've done a fair few times myself. But in practice, whilst this does help in many ways with the security challenge, this starts to become problematic for other reasons. My understanding here is that this all stems from Alpine's use of `musl` rather than `glibc` as its standard C library. The most noticable impact is creating **horrendous** build times when certain (common) packages are involved, as there are no pip wheels available (and sometimes out-right failing to build at all).

> Personally, as someone who does a lot of work with Google Cloud, gRPC (`grpcio`) is a particular victim of this. If memory serves it's true for `pandas` too

This isn't the only issue you have by using Python on Alpine. As [this excellent blog post](https://pythonspeed.com/articles/alpine-docker-python/) elaborates on, this - the TL:DR is that as well as the build time problem, you can encounter sometimes very subtle bugs and performance issues as a result of the choice too.

So, warned off that, we're left with Debian variants (e.g. `python:*-slim-bullseye`). I'd imagine these are by far the most common in the wild - but they too suffer from the presence of OS library vulnerabilities in the base image (Debian 11), e.g.

```sh
> trivy image -s=HIGH,CRITICAL python:3.9-slim-bullseye

# [...]
2022-06-26T13:59:34.109+0100	INFO	Detecting Debian vulnerabilities...

python:3.9-slim-bullseye (debian 11.3)

Total: 16 (HIGH: 13, CRITICAL: 3)
```

A slightly higher - but different set - of vulnerabilities being reported. This reflects that on the one hand this image has been built more recently and picked up a number of fixes, but also has more "stuff" in it than the distroless Google Python image scanned above.

> In my opinion, the Debian team do a good job in general justifying why they are or aren't patching certain CVEs - but regardless of your view on that, simply having to act at all - either to suppress in your vulnerability management tool (after reviewing it carefully) or address the issue - across many images on a regular basis is extremely [toilsome work](https://sre.google/sre-book/eliminating-toil/).

So, only one thing for it really! Let's see if we can build our own version that is a happy story from a vulnerability perspective AND still works ok.

---

## Building an Alternative

{{< figure src="/images/welding.jpg?width=1000px&classes=shadow" attr="Photo by Nazarii Yurkov on Unsplash" attrlink="https://unsplash.com/photos/m-VhHYQ4yFg" >}}

I initially began as probably most people would - taking a copy of the [distroless repo itself](https://github.com/GoogleContainerTools/distroless) and hacking about with it. As you'd expect, the images are built using Google's build tool (which they open sourced), [Bazel](https://bazel.build/). I'd not used this before and to be honest I found the learning curve rather steep. In any event, after a little while getting to grips with it, somewhere along the way I found myself thinking _"so how do I actually install my own stuff in here?"_. I came to the realisation that what I was really looking for was the equivalent of the Dockerfile to modify. Cue a pivot ...

> With both the Python upstream and Distroless both being based on Debian, I figured I could do exactly that with a [multi-layer build](https://docs.docker.com/develop/develop-images/multistage-build/), and also be happy that I was using tech that I understand well too.

I therefore ended up with a Dockerfile that looks a bit like this structure:

```dockerfile
FROM A_FULLY_FEATURED_PYTHON_UPSTREAM_IMAGE as python-base-image

### do nothing, if you like. Or do something. Up to you. So casual

FROM GOOGLE_DISTROLESS_IMAGE

# copy the files we need from the python-base-image into distroless
COPY --from=python-base-image /path/on/python-base-image /path/on/distroless

### do anything else you find useful here too. If you want to. Don't have to

# run python
ENTRYPOINT ["/usr/local/bin/python"]
```

The `COPY --from` is key - and in reality a _little_ more complicated than I've made it look. We need to copy Python and its dependencies into distroless, as well as any useful compiled libraries that we'll need for other Python packages. This involves a bit of trial and error _(unless you know the internals of Python and where it installs things better than I do, of course!)_ 🤓 . I initially started by copying huge chunks of `/lib` and `/usr` to get it working, then wittled that down to something that delivered on one of Distroless' goals of a small image (I ended up with an image of 57.7Mb. The Google Python version is 50.2Mb. One of my early versions clocked in at 300Mb+!)

In addition, for the distroless image I use as the final base layer, I use the [C version](https://github.com/GoogleContainerTools/distroless/tree/main/cc) rather than `base`, as so much of Python depends on it anyway. I pinched this idea from the Python variant that Google build.

The [final version of my Dockerfile](https://github.com/alexdmoss/distroless-python/blob/main/distroless.Dockerfile) looks a bit more complicated than the above - I'll talk through this in a bit more detail below:

{{< gist alexdmoss 3f32b62358a6677d513d735b04e93912 >}}

Breaking this down a bit:

1. I've set up some build-args to help with [the CI](https://github.com/alexdmoss/distroless-python/blob/main/.gitlab-ci.yml) - this allows me to build a Python 3.9 and 3.10 version out of the same repo, as well as give me the option of building (locally) for Apple M1 Silicon too if I wish. It's also important to build yourself a set of ["debug" images](https://github.com/GoogleContainerTools/distroless#debug-images) - distroless images but with a busybox shell in them. This is very helpful for remote access to debug problems (but _not_ for production - that defeats the point! 😏).
2. The `python-base` (PYTHON_BUILDER_IMAGE) builder image is not complex - [see Dockerfile](https://github.com/alexdmoss/distroless-python/blob/main/builder.Dockerfile). It uses Python Slim for the version we want, so we don't need to worry about installing Python's dependencies and Python itself. Tools like `pipenv` and `poetry` are installed as a convenience when using this image for other purposes.
3. In addition to copying in Python itself, I also copy in some other C libraries that many common packages I use depend on. You may not need these, but I use them frequently enough for it to be worth putting in the image. One of the examples I'll highlight below shows how to add these in yourself without needing to update the distroless image, if you need to (the need to do this usually reveals itself as an `ImportError` from Python - try it with `pandas` if you want to see an example).
4. I set up a non-root user, which involves briefly adding a shell then removing it to allow this. Don't run containers as root unless you really have to. Even though containers are ~~contained~~ 🥁 limited by default in what they can access on the host (unless you're even worse and run them with the `--privileged` flag - really really don't do that).
5. I set up some sane Python defaults. I think I originally learnt this from [this blog post](https://sourcery.ai/blog/python-docker/).

This build process if I'm sure far less deterministic than the Bazel approach if knowing **exactly** what is in your container image is crucial to you - but for my needs, it works well and is something I can understand and adapt more easily too.

Want to use the images?

- [mosstech/python-distroless](https://hub.docker.com/r/mosstech/python-distroless/tags) - e.g. `docker pull mosstech/python-distroless:3.9-debian11`
- [mosstech/python-builder](https://hub.docker.com/r/mosstech/python-builder/tags) - e.g. `docker pull mosstech/python-builder:3.10-debian11`

---

## Testing The Thing

{{< figure src="/images/test.jpg?width=1000px&classes=shadow" attr="Photo by Nicolas Thomas on Unsplash" attrlink="https://unsplash.com/photos/3GZi6OpSDcY" >}}

That's all well and good hacking about with Dockerfiles like this, but does it actually work? I cobbled together a few tests to satisfy myself that the approach was viable - you can see this [in my github repo](https://github.com/alexdmoss/distroless-python/tree/main/tests) too. They also serve as handy [examples](https://github.com/alexdmoss/distroless-python/blob/main/EXAMPLES.md) to reference later on how to use the image.

Stepping through them briefly:

1. [Hello World](https://github.com/alexdmoss/distroless-python/tree/main/tests/hello-world/). It's a classic for a reason. Nothing fancy here - just a simple print to console to assure the basic image (and its debug variant) work.
2. [Flask / Gunicorn](https://github.com/alexdmoss/distroless-python/tree/main/tests/gunicorn/). I'm involved in a fair few Flask-based apps - usually for doing "something useful" in Kubernetes. This simple example was the first that caused me to refactor the distroless Dockerfile significantly, as it highlighted both:
   - The challenge of running things with a different entrypoint. You can't just run `/gunicorn` or `pipenv` and expect it to work - see how this is solved in [`run.py`](https://github.com/alexdmoss/distroless-python/blob/main/tests/gunicorn/run.py) instead. [This blog post](https://blog.krybot.com/a?ID=01750-c0a1cb24-e375-4f31-86ee-e3280ed5302f) provided significant help for the run.py solution, as you'll see!
   - The dependency on additional OS libraries that I hadn't copied over originally (such as `libz.so.1`).
3. [Fast API](https://github.com/alexdmoss/distroless-python/tree/main/tests/fastapi/). This was followed up with a [FastAPI](https://fastapi.tiangolo.com/) example as something we're using more frequently where I work. This had very similar needs to the Flask/Gunicorn example above as you'd expect, with [uvicorn's docs](https://www.uvicorn.org/deployment/#running-programmatically) helping me out with what's needed in `run.py` here. Note also that in this Dockerfile I used my own `mosstech/python-builder` docker image (the thing we used to layer up the Distroless image itself) as my base layer. This offers a minor advantage in that it has `pipenv` or `poetry` pre-installed into it.
4. [Pandas](https://github.com/alexdmoss/distroless-python/tree/main/tests/pandas/). I rarely do things with Pandas currently, but have in the past and this served as a good illustration of additional dependencies that might be required. Note [this line](https://github.com/alexdmoss/distroless-python/blob/main/tests/pandas/Dockerfile#L19) where `libbz2.so.1.0` needs to be copied in from the builder image as well, as this isn't in the Distroless base I'd created (but `numpy` needs it). Without this, you'll see a (reasonably easy to reason about) `ImportError` when running the app. If you encounter anything similar to this, I've found the easiest way to deal with it is to spin up the builder image locally and just `find` the mentioned library on the filesystem to figure out what needs copying over. Might involve a bit of whack-a-mole if there's a few of them ...
5. [Google Cloud](https://github.com/alexdmoss/distroless-python/tree/main/tests/google-cloud/). I [work a lot with Google Cloud](https://alexos.dev/2019/02/23/a-year-in-google-cloud/), and GCP services + Alpine Python is a grim experience as mentioned earlier, so this was an important test for me personally. We're using a GCP PubSub Emulator as a fake here, as this exercises the Google Cloud PubSub Python Client Libraries, which depend on `grpcio`, which takes about 10 minutes to build in an alpine container on my reasonably modern system. That's a really cr@ppy cycle time - why that's bad is well-established now I think, but [here's a good article about why](https://steven-lemon182.medium.com/a-guide-to-reducing-development-wait-time-part-1-why-9dcbbfdc1224).
6. [Kubernetes](https://github.com/alexdmoss/distroless-python/tree/main/tests/kubernetes/). As so much of my Pythony things revolve around Kubernetes, checking that those client libraries work okay also was an important one to include, even though in truth it doesn't add a huge amount to the other examples above. The test is a bit slow unfortunately so I may drop it or move it to a manual smoke test after release perhaps - spinning up a [`kind`](https://kind.sigs.k8s.io/) cluster to minimise our dependency on something else to connect to takes time.

---

## Summary

So there we have it - a Distroless Python image that uses Google's distroless as a base, but layers in an up-to-date version of Python and its dependencies that are under your control, to tailor to your needs. Whilst this is something I've only recently put together, I'm now using it in a few places and it seems to be working well - whilst still delivering on the tiny image size (faster container startup matters!) and in particular leaner attack surface and less toil managing OS vulnerabilities (real or benign).

I hope you find some inspiration from this post to use it yourself or create your own equivalent.

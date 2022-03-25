---
title: "Migrating a JVM Application From GKE To Cloud Run"
panelTitle: "Migrating a JVM Application to Cloud Run"
date: 2019-12-14T18:00:00-01:00
author: "@alexdmoss"
draft: true
description: "Migrating a Kotlin JVM-based containerised application from a Google Kubernetes Engine cluster to Google Cloud Run"
banner: "/images/???.jpg"
tags: [ "GKE", "GCP", "Cloud Run", "Kubernetes", "Kotlin", "Berglas", "JIB" ]
---

{{< figure src="/images/squeeze-1.jpg?width=600px&classes=shadow" attr="Photo by ??? on Unsplash" attrlink="https://unsplash.com/photos/???" >}}

// summary

---

## Context

Today I thought I'd write about some recent experiences trying out Cloud Run for the first time in anger. The scenario in question is this:

- At work, we run a reasonably large multi-tenant Google Kubernetes Engine cluster. This includes a "tenant" for folks to learn about the stack - a training namespace
- One of our engineers built a useful tool and ran it there (it happens to be for helping with Change Management fun and games. That's another story). That namespace is really intended for more temporary, potentially it is a little more brittle than other services we run, and so on

So, the obvious choice might be to just move the useful tool to another namespace. Unless we were wanted to change the team that owns the service, we'd need to spin up a whole separate namespace for it. Doable, but perhaps a bit overkill? The workload is already dockerised and ready to roll on Kubernetes, but when we create a new namespace, what we really mean is a new Service, which comes with a whole load of other stuff that is useful stuff for a team to want, but pretty heavyweight for a utility tool run by (at least for now) a single developer.

I thought this sounded like quite a good use-case for Google Cloud Run - Google's new "serverless meets containers" (ugh, buzzword alert!) Knative-based runtime. It's recently gone GA too.

What follows are the steps I went through to migrate it, and some of the things I learnt about along the way.

---

## Getting Your GCP Project Up and Running

I had a GCP project ready to go for this bit of work already - it had a handful of other similar tools running there (mostly in AppEngine), but I needed to switch on the Cloud Run functionality. The project in question isn't managed through terraform so we do this imperatively (that'll be a bit of a theme through the instructions too - as far as I know everything that follows *can* be declared through Terraform HCL no worries, if you prefer) - `gcloud services enable run.googleapis.com`.

> For the purposes of this post, lets assume the GCP Project is `GCP_PROJECT_ID=my-amazing-project`. It's not. Perhaps it should be ...

Our Google Container Registry also lives in a different project so we need to handle some access requirements for Cloud Run running in Project `my-amazing-project` - it needs to be able to pull the image.

/// Add instructions here

## Deploying The Service For The First Time

So far, so good - we can just deploy it now right? 

```sh
gcloud run deploy my-awesome-tool \
  --image=eu.gcr.io/${GCR_GCP_PROJECT}/${IMAGE_NAME}:latest \
  --platform=managed \
  --region=europe-west1 \
  --project=${GCP_PROJECT_ID} \
  --no-allow-unauthenticated
```

> Don't use `:latest` by the way. I am being lazy here!

/// what happens here - should run but service doesn't work

Some immediate observations:

- How do I actually get passed the authentication and use it?
- I ran this manually, from my command line. Where's the pipeline bro?
- the URL is pretty nasty. It'd be nice to fix that
- The service doesn't work, as it actually needs some environment variables set
- We know that this tool can take a while to do its thing. Should we add some protections for the amount of resources it can consume?


---

## Invoking the Service

```sh
export GCP_PROJECT_ID=jl-engineering
export SERVICE_NAME=change-admin
export URL=$(gcloud run services list --filter=SERVICE:${SERVICE_NAME} --format='value(URL)' --platform=managed --project=${GCP_PROJECT_ID})
curl -H "Authorization: Bearer $(gcloud auth print-identity-token)" "${URL}/gitlab-status/change-admin?groupdId=1233"
```

Note: groupId=1233 = the `digital/services/catalogue` group - they usually have some commits to view!

Or, if you're me, and you have stuff aliased: `gcurl ${URL}/gitlab-status/change-admin?groupId=1233`

Existing endpoint, for comparison: https://stubs-api.training.flex.jl-digital.net/gitlab-status/change-admin?groupId=1233

---

## Deploying From Local

```sh
export GCP_PROJECT_ID=jl-engineering
export CI_COMMIT_SHA=9b0b508b
export BUCKET_NAME=jl-change-admin-store

gcloud run deploy change-admin \
  --image=eu.gcr.io/jl-container-images/alexmoss/gitlab-status:latest \
  --platform=managed \
  --region=europe-west1 \
  --project=${GCP_PROJECT_ID} \
  --timeout=60 \
  --no-allow-unauthenticated
```

We should also consider setting one or more of the following, to set some limits on the resources it can consume:

- `--concurrency=1`
- `--max-instances=1`
- `--memory=128Mi`

---

## Setup

### Service Account for CI

```sh
gcloud iam service-accounts create change-admin-ci --description 'CI account to build & deploy JL Change Admin tool'
gcloud projects add-iam-policy-binding jl-engineering --member="serviceAccount:change-admin-ci@jl-engineering.iam.gserviceaccount.com" --role="roles/run.admin"
gcloud projects add-iam-policy-binding jl-engineering --member="serviceAccount:change-admin-ci@jl-engineering.iam.gserviceaccount.com" --role="roles/iam.serviceAccountUser"
gcloud projects add-iam-policy-binding jl-container-images --member="serviceAccount:change-admin-ci@jl-engineering.iam.gserviceaccount.com" --role="projects/jl-container-images/roles/ContainerRegistryTenant"
```

### Granting Access to GCR

For the above to work, we need to give Cloud Run permission to read from the JLDP Container Registry:

```sh
gsutil iam ch serviceAccount:service-286924115525@serverless-robot-prod.iam.gserviceaccount.com:objectViewer gs://eu.artifacts.jl-container-images.appspot.com
```

Where the GCR in question is in the project `jl-container-images` and the serviceAccount name is taken from the IAM screen in `jl-engineering`, for the Cloud Run Service Agent (generated by Google).

### Secret Handling

We also need to deal with the fact that the GITLAB_TOKEN the service uses is secret. This is more complicated on Cloud Run than you'd like - you can't use Kubernetes Secrets there, and I was keen to avoid modifying the Java code. In hindsight, that option may have been easier!

Instead, we use one of Google's recently open-sourced tools - [Berglas](https://github.com/GoogleCloudPlatform/berglas) - which is designed to help with Secrets Management in serverless environments.

Setting up Berglas is straight-forward, but getting it working with `jib` is messier

#### Berglas Setup

```sh
# enable KMS & GCS
gcloud services enable --project ${GCP_PROJECT_ID} cloudkms.googleapis.com storage-api.googleapis.com storage-component.googleapis.com
# creates KMS keyring and bucket
berglas bootstrap --project $GCP_PROJECT_ID --bucket $BUCKET_NAME --bucket-location eu
# create the secret, encrypted using the keyring
berglas create ${BUCKET_NAME}/gitlab $(cat ~/.secret/gitlab-status) --key=projects/${GCP_PROJECT_ID}/locations/global/keyRings/berglas/cryptoKeys/berglas-key
# grant the Default Compute SA access to read the file, and then decrypt the secret
export SA_EMAIL=$(gcloud projects describe ${GCP_PROJECT_ID} --format 'value(projectNumber)')-compute@developer.gserviceaccount.com
gcloud projects add-iam-policy-binding ${GCP_PROJECT_ID} --member serviceAccount:${SA_EMAIL} --role roles/run.viewer
berglas grant ${BUCKET_NAME}/gitlab --member serviceAccount:${SA_EMAIL}
```

#### Using Berglas at Runtime

For this to work, the runtime needs to "understand" berglas - howto guides point to doing this sort of thing in your Dockerfile:

```Dockerfile
COPY --from=gcr.io/berglas/berglas:latest /bin/berglas /bin/berglas
ENTRYPOINT exec /bin/berglas exec -- <your-process>
```

This is not easy with jib. After some experimentation, I modified `build.gradle` to create an endpoint bash script (annoyingly, this blocks using a Google distroless image), and copying the `berglas` [executable](https://storage.googleapis.com/berglas/master/linux_amd64/berglas) into the container (ick). It would perhaps be better to bake a JLDP GCR-hosted base image for this instead - but it does work.

The entrypoint script is as follows:

```sh
export GITLAB_TOKEN=$(/berglas access jl-change-admin-store/gitlab)
java -cp /app/resources:/app/classes:/app/libs/* http4k.ApplicationKt
```

Where `jl-change-admin-store` is `${BUCKET_NAME}` from the above, and the `java` execution is taken from the `docker inspect` of a jib-created docker image. Bit hacky!

This can be tested with:

- `./gradlew clean jibDockerBuild`
- `docker run -p 8080:8080 -e GOOGLE_APPLICATION_CREDENTIALS=/secret/gitlab-status.json -v ${PATH_TO_GOOGLE_CREDS_JSON}:/secret/gitlab-status.json:ro ${IMAGE_NAME}:latest`

### Custom Domain

I'm not a big fan of an auto-generated DNS entry - while Google say it won't change, what if we tear down the service and rebuild it? I'm not sure we can guarantee it is retained, and then any client would need to be notified, update their config, etc.

Ideally, we'd going to register a custom domain under `*.jl-engineering.net`. This is a good option because: a) it maps to this GCP Project's name, and b) it is already verified in GCP (`gcloud domains list-user-verified`). For the experiment though (and because I don't have access to jl-engineering.net's DNS), the following was used:

```sh
gcloud beta run domain-mappings create --service change-admin --domain change-admin.experiments.jl-digital.net --platform=managed --region=europe-west1
```

This prompts you to create a DNS entry for `change-amdin.experiments.jl-digital.net --> CNAME to --> ghs.googlehosted.com`. I did this manually because we don't have any terraform for the `experiments` sub-domain.

The domain's status can be verified with `gcloud beta run domain-mappings list`. It looks like the Cloud Console page for Cloud Run does not update with these values (it is beta, I guess ...).

At this point hitting `https://change-admin.experiments.jl-digital.net/gitlab-status` should work.

---

## Granting Developer Access

You can specify an individual's access as follows:

```sh
gcloud run services add-iam-policy-binding change-admin --member='user:alex.moss@johnlewis.co.uk' --role='roles/run.invoker' --project=${GCP_PROJECT_ID} --region=europe-west1 --platform=managed
```

It also appears to be possible to give the whole of JL access:

```sh
gcloud run services add-iam-policy-binding change-admin --member='domain:johnlewis.co.uk' --role='roles/run.invoker' --project=${GCP_PROJECT_ID} --region=europe-west1 --platform=managed
```

Note that, as you might guess, `--member='serviceAccount:${CLIENT_SA_EMAIL}'` would also work - allowing for programmatic access.

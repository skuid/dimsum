# 🍜 dimsum

_docker image manifest summaries_

`dimsum` is a server that connectes to multiple (private) Docker registries and serves image manifests. It also serves the `History.V1Compatibility` object for direct, detailed access.

## Disclaimer
tl;dr: Don't put `dimsum` on the public internet. 

`dimsum` should only be hosted within the Kuberentes cluster (or host) that Spinnaker is running in. Firstly, this service serves HTTP, not HTTPS. Also, it is unauthenciated which means that you would be serving metadata about potentially private images to the world. 




## Purpose

Why `dimsum`? In order to utilize image metadata within pipelines, Spinnaker needs a way of accessing this information. If you use a Docker Trigger, you don't get all of the image metadata like you would with a Git trigger. You can use `LABEL` to annotate your Docker image with useful information and then use the Webhook stage to obtain this information for use within your pipeline.


Since you should be running `dimsum` alonside Spinnaker, you should be able to configure your Webhook stage to query the API with a hostname that Orca can access (since that's where the SPEL is evaluated).

Example:

If you're using a Docker Registry trigger and running `dimsum` on Kubernetes, your Webhook stage may call it like this:

```
http://dimsim.spinnaker.svc.cluster.local:8080/${trigger['account']}/${trigger['repository']}/${trigger['tag']}/history?level=0
```

The trigger information will be injected into the request and `dimsum` query the registry API that that image came from. The metadata about your image will be returned, similar to if you were to do a `docker inspect` on the same image.

It can then be accessed using SPEL in subesequent pipeline stages:

For instance, if you have a `revision` label,
```
${#stage('Webhook')['context']['buildInfo']['config']['Labels']['revision']}
```

## Configuration
Configuration is done via YAML. As a matter of fact, you can reuse your Clouddriver configuration for `dockerRegistries`. If you're running Spinnaker in Kuberentes, you can run `dimsum` in the same Pod as Clouddriver and utilize the same configuration.

For supported configuration, read the Spinnaker docs on [configuring a Docker Registry](http://www.spinnaker.io/v1.0/docs/target-deployment-configuration#section-docker-registry)

The default config path is `config.yaml` in the same directory as the executable. The `--config` flag can be used to override this behavior.

## TODO
* Docker Image
* Handle IndexOutOfBounds possibility
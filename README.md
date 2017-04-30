# dimsum

_docker image manifest summaries_

`dimsum` is a server that connectes to multiple (private) Docker registries and serves image manifests. It also serves the `History.V1Compatibility` object specifically for one reason.

## Disclaimer
tl;dr: Don't put `dimsum` on the public internet. 

`dimsum` should only be hosted within the Kuberentes cluster (or host) that Spinnaker is running in. Firstly, this service serves HTTP, not HTTPS. Also, it is unauthenciated which means that you would be serving metadata about potentially private images to the world. 




## Purpose

Why `dimsum`? In order to support annotating Kubernetes Deployments with a specific Git revision, Spinnaker needs a way of accessing this information. If you use a Docker Trigger, you don't get all of the Git information about a build so you need to find another way. The idea is to add a `LABEL` to your Docker image for the Git revision of an image. Then, the Spinnaker pipeline can call the `dimsum` endpoint with `jsonFromUrl()` and tease out this information.

A Spinnaker pipeline expression to do this might look like this:

```
${readJSON(jsonFromUrl('http://dimsum.spinnaker/dockerhub/library/nginx/latest/history?level=0'))['container_config']['Labels']['GIT_REVISION']}
``` 

_Yes, I am aware that that is covoluted_

## Configuration
Configuration is done via YAML. As a matter of fact, you can reuse your Clouddriver configuration for `dockerRegistries`. If you're running Spinnaker in Kuberentes, you can run `dimsum` in the same Pod as Clouddriver and utilize the same configuration.

For supported configuration, read the Spinnaker docs on [configuring a Docker Registry](http://www.spinnaker.io/v1.0/docs/target-deployment-configuration#section-docker-registry)

The default config path is `config.yaml` in the same directory as the executable. The `--config` flag can be used to override this behavior.

## TODO
* Docker Image
* Support Registries that have a `passwordFile` configuration
* Handle IndexOutOfBounds possibility
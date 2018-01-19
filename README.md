# 🍜 dimsum

_docker image manifest summaries_

`dimsum` is a server that connectes to multiple (private) Docker registries and serves image manifests. It also serves the `History.V1Compatibility` object for direct, detailed access.

## Disclaimer
tl;dr: Don't put `dimsum` on the public internet.

`dimsum` should only be hosted within the Kubernetes cluster (or host) that Spinnaker is running in. Firstly, this service serves HTTP, not HTTPS. Also, it is unauthenciated which means that you would be serving metadata about potentially private images to the world.


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

## Docker Image

```
quay.io/skuid/dimsum:v1.0.0
quay.io/skuid/dimsim:latest
```

## Configuration
Configuration is done via YAML. As a matter of fact, you can reuse your Clouddriver configuration for `dockerRegistries`. If you're running Spinnaker in Kuberentes, you can run `dimsum` in the same Pod as Clouddriver and utilize the same configuration.

For supported configuration, read the Spinnaker docs on [configuring a Docker Registry](http://www.spinnaker.io/v1.0/docs/target-deployment-configuration#section-docker-registry)

The default config path is `config.yaml` in the same directory as the executable. The `--config` flag can be used to override this behavior.

## Spinnaker Preconfigured Webhook

Instead of using a generic Webhook stage to obtain this information, it may be easier to use a Preconfigured Webhook. You can read more about it [here](https://medium.com/@e_frogers/custom-spinnaker-stages-with-preconfigured-webhooks-84c5b5dae861). If you'd like to make a custom stage for Dimsum, add the following snippet to `orca-local.yml`:

```yaml
webhook:
 preconfigured:
   - label: Dimsum
     description: Retrieve trigger image metadata
     type: dimsum
     url: http://<dimsum-host>/${parameterValues['account']}/${parameterValues['repository']}/${parameterValues['tag']}/history?level=0
     method: GET
     payload: {}
     parameters:
        - name: account
          label: Account
          defaultValue: ${trigger['account']}
        - name: repository
          label: Repository
          defaultValue: ${trigger['repository']}
        - name: tag
          label: Tag
          defaultValue: ${trigger['tag']}
     waitForCompletion: false
```

## TODO
* Handle IndexOutOfBounds possibility

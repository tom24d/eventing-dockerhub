# Knative Eventing DockerHub Source

Knative Eventing `DockerHubSource` defines an event source that transforms webhook events
from hub.docker.com into CloudEvents and deliver to the specified sink in the configuration yaml.

To learn more about Knative, please visit
[Knative docs](https://github.com/knative/docs) repository.

If you are interested in contributing, see [CONTRIBUTING.md](./CONTRIBUTING.md)
and [DEVELOPMENT.md](./DEVELOPMENT.md).

This project is inspired by [the idea of JBoss community](https://docs.jboss.org/display/GSOC/Google+Summer+of+Code+2020+ideas#GoogleSummerofCode2020ideas-Knative-Eventsourcesforcontainerregistries,pipelinesandbuilds).


## Before you begin

`DockerHubSource` installation requires two knative component on your kubernetes cluster.  
Plus, you need a build tool `ko`.

- Knative Serving core  
- Knative Eventing core  
- [ko](https://github.com/google/ko)  


## Installation  
Install DockerHubSource from the source:

```bash
ko apply -f config
```

## DockerHubSource usage example

Applying example(autoCallback enabled)

```bash
kubectl apply -f ./example/normal-display.yaml
kubectl apply -f ./example/source.yaml
```

The examples have also `callback-display.yaml` to try `disableAutoCallback=true` mode.  
Note that `callback-display.yaml` needs `ko` to apply whereas appropriate sink is in the `source.yaml`.  

You can see the resource is created via: `kubectl get dockerhubsource`.  
```
% kubectl get dockerhubsources.sources.knative.dev
NAME               READY   REASON   URL                                          AGE
dockerhub-source   True             http://<your-endpoint-for-DockerHubSource>   17s

```  

Copy `http://<your-endpoint-for-DockerHubSource>` to configure dockerhub webhook. See [DockerHub Reference](https://docs.docker.com/docker-hub/webhooks/).  

## API Reference

### v1alpha1

|Field|Description|
|:---|:---|
|`apiVersion`  <br>string| `sources.knative.dev/v1alpha1`|
|`kind` <br> string| `DockerHubSource`|
|`metadata` <br> [Kubernetes<br>meta/v1.metadata](https://github.com/knative/docs/blob/master/docs/reference/eventing/eventing-contrib.md#duck.knative.dev/v1.CloudEventOverrides)| Refer to the Kubernetes API documentation for the fields of the `metadata` field.|
|`Spec` <br> DockerHubSourceSpec| <table> <tr> <td><code> disableAutoCallback </code> <br> bool</td> <td> (Optional) <br> DisableAutoCallback configures whether the adapter works with automatic callback feature. Docker Hub webhook needs validation callback to receive continually its chain. If the field is false, the adapter automatically sends a corresponding callback. When the event gets delivered successfully, callback status is `success`. Otherwise the status is `failure`.  If unspecified, this will default to false.</td> </tr> <tr> <td><code> SourceSpec </code> <br> <a href="https://github.com/knative/docs/blob/master/docs/reference/eventing/eventing-contrib.md#duck.knative.dev/v1.SourceSpec">SourceSpec</a></td> <td>(Members of SourceSpec are embedded into this type.)</td> </tr> </table>|
|`Status` <br> DockerHubSourceStatus| <table>  <tr> <td><code> SourceStatus </code> <br> <a href="https://github.com/knative/docs/blob/master/docs/reference/eventing/eventing-contrib.md#duck.knative.dev/v1.SourceStatus">SourceStatus</a></td> <td>(Members of SourceStatus are embedded into this type.)</td> </tr> <tr> <td><code> AutoCallbackDisabled </code> <br> bool</td> <td>  AutoCallbackDisabled is the status whether automatic callback is disabled.</td> </tr> <tr> <td><code> URL </code> <br> knative.dev/pkg/apis.URL</td> <td> (Optional) <br> URL is the current active allocated URL that has been configured for the Source endpoint.</td> </tr> <tr> <td><code> FirstServiceName </code> <br> string</td> <td> (Optional) <br> FirstServiceName holds the information of knative service name to recreate service when accidentally deleted. </td> </tr> </table>|

---
### Example CloudEvent payload

```
Validation: valid
Context Attributes,
  specversion: 1.0
  type: dev.knative.source.dockerhub.push
  source: https://hub.docker.com/r/tom24d/postwebhook
  subject: tom24d
  id: fd4c1670-b126-4289-b025-61579ceeee3d
  time: 2020-06-19T05:32:48Z
  datacontenttype: application/json
Extensions,
  tag: latest
Data,
  {
    "callback_url": "https://registry.hub.docker.com/u/tom24d/postwebhook/hook/2g51gffa2h1b4eh3ed355je0bc4g0e2h/",
    "push_data": {
      "images": [],
      "pushed_at": 1592544800,
      "pusher": "tom24d",
      "tag": "latest"
    },
    "repository": {
      "comment_count": 0,
      "date_created": 1591154200,
      "description": "[test] dockerhub website description.",
      "dockerfile": "",
      "full_description": "",
      "is_official": false,
      "is_private": false,
      "is_trusted": false,
      "name": "postwebhook",
      "namespace": "tom24d",
      "owner": "tom24d",
      "repo_name": "tom24d/postwebhook",
      "repo_url": "https://hub.docker.com/r/tom24d/postwebhook",
      "star_count": 0,
      "status": "Active"
    }
  }
```

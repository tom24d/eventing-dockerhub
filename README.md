# [WIP]Knative Eventing Dockerhub Source

Knative Eventing `dockerhub-source` defines an event source that transforms webhook events
from hub.docker.com into CloudEvents and deliver to the specified sink in the configuration yaml.

To learn more about Knative, please visit
[Knative docs](https://github.com/knative/docs) repository.

If you are interested in contributing, see [CONTRIBUTING.md](./CONTRIBUTING.md)
and [DEVELOPMENT.md](./DEVELOPMENT.md).

This project is inspired by [the idea of JBoss community](https://docs.jboss.org/display/GSOC/Google+Summer+of+Code+2020+ideas#GoogleSummerofCode2020ideas-Knative-Eventsourcesforcontainerregistries,pipelinesandbuilds).

# DockerHubSource usage example

Make sure you have `ko`. If you don't, see [link](https://github.com/google/ko).

1. Install DockerHubSource

```bash
ko apply -f config
```

2. apply example(autoCallback enabled)

```bash
kubectl apply -f ./example/normal-display.yaml
kubectl apply -f ./example/source.yaml
```

The examples have also `callback-display.yaml` to try autoCallback disabled mode.
Note that you have to apply `callback-display.yaml` with `ko`.

<!-- TODO write with better style -->

3. You can see the resource is created via: `kubectl get dockerhubsource`.  
The endpoint should be available via: `kubectl get ksvc`. The output should be like:
```bash
 % k get ksvc
NAME                     URL                                                          LATESTCREATED                  LATESTREADY                   READY     REASON
dockerhub-source-jd6rq   http://<your-domain-for-DockerHubSource>   dockerhub-source-jd6rq-8kqrt                                 Unknown   RevisionMissing
event-display            http://<your-domain-for-event-display>            event-display-w58d4            event-display-w58d4           True      

```  

4. Copy `http://<your-domain-for-DockerHubSource>` to use to configure hub.docker.com webhook. See [Link](https://docs.docker.com/docker-hub/webhooks/).  

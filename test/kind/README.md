These scripts help development with [KinD](https://kind.sigs.k8s.io/) cluster.

```shell script
./test/kind/bootstrap.sh
```
creates KinD cluster with [local Docker Registry](https://kind.sigs.k8s.io/docs/user/local-registry/).  


```shell script
KO_DOCKER_REPO=localhost:port ./test/kind/run-tests.sh
```
runs integration test on KinD.

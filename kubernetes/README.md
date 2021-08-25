# [cassandra-web](https://orzhaha.github.io/cassandra-web/) 
## Installing the Chart - In this git repo
1.  To install the chart with the release name `cassandra-web` in namespace $NAMESPACE you defined above:

```bash
$ helm install cassandra-web --name cassandra-web --namespace $NAMESPACE
```

## Installing the Chart - Use our official helm chart repo

```bash
helm repo add orzhaha https://orzhaha.github.io/charts/orzhaha
helm search repo orzhaha
helm install my-release orzhaha/<chart>
```

## Configuration

Set the db host and port in the kubernetes/helm/cassandra-web/values.yaml

## Uninstall the Chart

To uninstall the `cassandra-web` release:

```bash
$ helm uninstall cassandra-web
```

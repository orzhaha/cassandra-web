# [cassandra-web](https://orzhaha.github.io/cassandra-web/) 
## Installing the Chart
1.  To install the chart with the release name `cassandra-web` in namespace $NAMESPACE you defined above:

```bash
$ helm install cassandra-web --name cassandra-web --namespace $NAMESPACE
```

## Configuration

Set the db host and port in the kubernetes/helm/cassandra-web/values.yaml

## Uninstalling the Chart

To uninstall/delete the `cassandra-web` release but continue to track the release:

```bash
$ helm delete cassandra-web
```

To uninstall/delete the `cassandra-web` release completely and make its name free for later use:

```bash
$ helm del --purge cassandra-web
```
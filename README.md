# oc-clusteroperator
OpenShift CLI plugin to change the state of ClusterOperators from managed to unmanaged and back again. Inspired by 
the [Kubernetes sample-cli-plugin](https://github.com/kubernetes/sample-cli-plugin).

## Testing
`make test`

## Building
`make plugin`

## Installing
`make install` will put the resulting `oc-clusteroperator` binary in `/usr/local/bin`. You can put it anywhere in 
your `$PATH` to access the plugin's functionality.
## gpupod
gpupod is simple tool to list and watch NVIDIA GPU pod in kubernetes cluster.

### Install
#### Linux
```bash
GOBIN=/usr/local/bin/ go install github.com/ide-rea/gpupod
```

### Usage
```
gpupod is a tool to list/watch NVIDIA GPU pod

Usage:
gpupod [flags]

Flags:
-t, --createdTime         show pod created time(default without created time)
-h, --help                help for gpupod
-k, --kubeconfig string   kubernetes config path (default "/path/to/home/.kube/config")
-r, --reason              show status reason(default without status reason)
-s, --success-pod         list pod include success pod(default not list success pod)
-w, --watch               watch gpu pod(default only list pod)
```


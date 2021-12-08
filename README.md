## gpupod
gpupod is simple tool to list and watch NVIDIA GPU pod in kubernetes cluster.

### Install
#### Linux
```bash
GOBIN=/usr/local/bin/ go install github.com/ide-rea/gpupod
```

### Usage
```
Usage:
gpupod [flags]

Flags:
-t, --createdTime         with pod created time
-h, --help                help for gpupod
-k, --kubeconfig string   kubernetes config path (default "/Users/zhangxiaoyu15/.kube/config")
-r, --reason              with pod created time
-s, --success-pod         with success pod
-w, --watch               watch gpu pod
```


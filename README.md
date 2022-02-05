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
--add_dir_header                   If true, adds the file directory to the header
--alsologtostderr                  log to standard error as well as files
-t, --createdTime                      show pod created time(default without created time)
-h, --help                             help for gpupod
-k, --kubeconfig string                kubernetes config path (default "/root/.kube/config")
--log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
--log_dir string                   If non-empty, write log files in this directory
--log_file string                  If non-empty, use this log file
--log_file_max_size uint           Defines the maximum size a log file can grow to. Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
--logtostderr                      log to standard error instead of files (default true)
-r, --reason                           show status reason(default without status reason)
--request-timeout int              list watch request timeout seconds(default zero)
--skip_headers                     If true, avoid header prefixes in the log messages
--skip_log_headers                 If true, avoid headers when opening log files
--stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
-s, --success-pod                      list pod include success pod(default not list success pod)
-v, --v Level                          number for the log level verbosity
--vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
-w, --watch                            watch gpu pod(default only list pod)
```


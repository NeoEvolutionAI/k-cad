# k-cad
Kubernetes cadvisor collector. 

## About 
k-cad is a simple collector for kubernetes `cadvisor`. `cadvisor` already expose its
metrics in a prometheus format, and the only way to get its metrics you need 
to do a static config for `k8s-api/nodes/<node-name>/proxy/metrics/cadvisor` of course
assuming you have `cadvisor` running in your cluster!

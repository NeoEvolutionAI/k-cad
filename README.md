# k-cad

Kubernetes cadvisor to prometheus exporter. 

![License](https://img.shields.io/github/license/NeoEvolutionAI/k-cad?color=orange&logoColor=red)
![Version](https://img.shields.io/badge/version-0.1.0-green)

## About 

k-cad is a simple exporter for Kubernetes `cadvisor`. `cadvisor` already expose its
metrics in a Prometheus format, and the only way to get its metrics you need 
to do a static config for `<host>/api/v1/nodes/<node-name>/proxy/metrics/cadvisor` of course
assuming you have `cadvisor` running in your cluster!

## Assumptions

* Authentication happening from within a pod.
* Using Kubernetes DNS server `kube-dns.`
* Running k8s v1.15+

## How to use it

k-cad should run as a Kubernetes pod in your cluster. Typically, you will need 
a single pod for your entire cluster.

To run k-cad successfully, you will need these three k8s resources:

- Deployment resource
- Service resource 
- Service account resource

You will also need to configure Prometheus configs to point to k8s service 
for k-cad. 

---

### Deployment resource example: 

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: k-cad-deployment
  namespace: monitoring
  labels:
    app: k-cad
spec:
  selector:
    matchLabels:
      app: k-cad
  template:
    metadata:
      labels:
        app: k-cad
    spec:
      serviceAccountName: admin-user
      containers:
      - name: k-cad
        image: "yourdockerregistry/k-cad:v0.1.0"
        env:
        - name: PORT
          value: "41001"
```

---

### Service account resource example

You will need to have a service account that will allow your pod to query Kubernetes API from within the pod. 

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-user
  namespace: monitoring
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: admin-user
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: admin-user
  namespace: monitoring

```

---

### Service resource example

You will also need to create a Kubernetes service so that you can point your Prometheus static configs too.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: k-cad-service
  namespace: monitoring
  labels:
    app: k-cad
    belongTo: monitoring
spec:
  ports:
  - name: http
    port: 41001
    targetPort: 41001
  selector:
    app: k-cad
```

---

Here is an example of your Prometheus configs to query k-cad

```yaml
- job_name: c-advisor
      static_configs:
      - targets:
        - k-cad-service.monitoring.svc.cluster.local:41001
```
package main

import (
	"fmt"
	"strings"
)

func manifestsAppend(old, new string) string {
	return fmt.Sprintf("%s\n%s", old, new)
}

func manifestRender(m manifestValues) string {
	var replacer = strings.NewReplacer(
		"{namespace}", m.namespace,
		"{serviceName}", m.serviceName,
		"{deploymentName}", m.deploymentName,
		"{deploymentLabelKey}", m.deploymentLabelKey,
		"{deploymentLabelValue}", m.deploymentLabelValue,
		"{ingressName}", m.ingressName,
		"{ingressHost}", m.ingressHost,
		"{ingressClass}", m.ingressClass,
		"{databaseURL}", m.databaseURL,
		"{image}", m.image,
	)
	str := replacer.Replace(manifestTemplate)
	return str
}

type manifestValues struct {
	namespace   string
	serviceName string

	deploymentName       string
	deploymentLabelKey   string
	deploymentLabelValue string

	ingressName  string
	ingressHost  string
	ingressClass string

	databaseURL string

	image string
}

var manifestTemplate = `
---
apiVersion: v1
kind: Namespace
metadata:
  name: {namespace}

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: grafana
  namespace: {namespace}

---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {ingressName}
  namespace: {namespace}
  annotations:
    kubernetes.io/ingress.class: "{ingressClass}"
spec:
  rules:
  - host: {ingressHost} # team.grafana.example.com
    http:
      paths:
      - backend:
          serviceName: {serviceName}
          servicePort: 3000
        path: /

---
apiVersion: v1
kind: Service
metadata:
  name: {serviceName}
  namespace: {namespace}
spec:
  ports:
  - name: http
    port: 3000
    targetPort: http
  selector:
    {deploymentLabelKey}: {deploymentLabelValue}

---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
    {deploymentLabelKey}: {deploymentLabelValue}
  name: {deploymentName}
  namespace: {namespace}
spec:
  replicas: 1
  selector:
    matchLabels:
      {deploymentLabelKey}: {deploymentLabelValue}
  template:
    metadata:
      labels:
        {deploymentLabelKey}: {deploymentLabelValue}
    spec:
      containers:
      - image: {image} # grafana/grafana:5.1.0
        name: grafana
        env:
        - name: GF_DATABASE_URL
          value: "{databaseURL}"
        ports:
        - containerPort: 3000
          name: http
        resources:
          limits:
            cpu: 400m
            memory: 400Mi
          requests:
            cpu: 100m
            memory: 100Mi
        volumeMounts:
        - mountPath: /var/lib/grafana
          name: grafana-storage
          readOnly: false
      securityContext:
        runAsNonRoot: true
        runAsUser: 65534
      serviceAccountName: grafana
      volumes:
      - emptyDir: {}
        name: grafana-storage
`

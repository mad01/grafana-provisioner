package main

import "strings"

func manifestRender(m manifestValues) string {
	var replacer = strings.NewReplacer(
		"{namespace}", m.namespace,
		"{serviceName}", m.serviceName,
		"{deploymentName}", m.deploymentName,
		"{deploymentLableKey}", m.deploymentLableKey,
		"{deploymentLableValue}", m.deploymentLableValue,
		"{ingressName}", m.ingressName,
		"{ingressHost}", m.ingressHost,
		"{ingressClass}", m.ingressClass,
	)
	str := replacer.Replace(manifestTemplate)
	return str
}

type manifestValues struct {
	namespace   string
	serviceName string

	deploymentName       string
	deploymentLableKey   string
	deploymentLableValue string

	ingressName  string
	ingressHost  string
	ingressClass string
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
	{deploymentLableKey}: {deploymentLableValue}

---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  labels:
	{deploymentLableKey}: {deploymentLableValue}
  name: {deploymentName}
  namespace: {namespace}
spec:
  replicas: 1
  selector:
    matchLabels:
	  {deploymentLableKey}: {deploymentLableValue}
  template:
    metadata:
      labels:
		{deploymentLableKey}: {deploymentLableValue}
    spec:
      containers:
      - image: grafana/grafana:5.1.0
        name: grafana
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

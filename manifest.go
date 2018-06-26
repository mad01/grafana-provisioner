package main

import (
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func manifestsAppend(old, new string) string {
	return fmt.Sprintf("%s\n%s", old, new)
}

func manifestRender(m Manifest) string {
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

type Manifest struct {
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

func (m *Manifest) Render() string {
	data := ""
	join := func(old, new string) string {
		return fmt.Sprintf("%s\n---\n\n%s", old, new)
	}

	data = join(data, m.RenderNamespace())
	data = join(data, m.RenderServiceAccount())
	data = join(data, m.RenderIngress())
	data = join(data, m.RenderService())
	data = join(data, m.RenderDeployment())

	return data
}
func (m *Manifest) RenderNamespace() string {
	o := v1.Namespace{}
	o.Name = m.namespace

	y, err := yaml.Marshal(o)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return ""
	}

	return string(y)

}
func (m *Manifest) RenderServiceAccount() string {
	o := v1.ServiceAccount{}
	o.Name = m.namespace
	o.Namespace = m.namespace

	y, err := yaml.Marshal(o)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return ""
	}

	return string(y)
}
func (m *Manifest) RenderIngress() string {
	o := extensions.Ingress{}
	o.Name = m.namespace
	o.Namespace = m.namespace
	o.Annotations = map[string]string{"kubernetes.io/ingress.class": m.ingressClass}
	o.Spec = extensions.IngressSpec{
		Rules: []extensions.IngressRule{
			{
				Host: m.ingressHost,
				IngressRuleValue: extensions.IngressRuleValue{
					HTTP: &extensions.HTTPIngressRuleValue{
						Paths: []extensions.HTTPIngressPath{
							{
								Path: "/",
								Backend: extensions.IngressBackend{
									ServiceName: m.serviceName,
									ServicePort: intstr.FromInt(3000),
								},
							},
						},
					},
				},
			},
		},
	}

	y, err := yaml.Marshal(o)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return ""
	}
	return string(y)
}
func (m *Manifest) RenderService() string {
	o := v1.Service{}
	o.Name = m.namespace
	o.Namespace = m.namespace
	o.Spec.Selector = map[string]string{m.deploymentLabelKey: m.deploymentLabelValue}
	o.Spec.Ports = append(
		o.Spec.Ports,
		v1.ServicePort{Name: "http", Port: 80, TargetPort: intstr.FromString("http")},
	)

	y, err := yaml.Marshal(o)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return ""
	}
	return string(y)
}
func (m *Manifest) RenderDeployment() string {
	runAsNonRoot := true
	runAsUser := int64(65534)
	selector := map[string]string{m.deploymentLabelKey: m.deploymentLabelValue}
	o := extensions.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.deploymentName,
			Namespace: m.namespace,
			Labels:    selector,
		},
		Spec: extensions.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels:      selector,
				MatchExpressions: []metav1.LabelSelectorRequirement{},
			},
			Replicas: int32(1),
			Template: api.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: selector,
				},
				Spec: api.PodSpec{
					Containers: []api.Container{{
						Name:  "grafana",
						Image: m.image,
						Env: []api.EnvVar{
							{
								Name:  "GF_DATABASE_URL",
								Value: m.databaseURL,
							},
						},
						Resources: api.ResourceRequirements{
							Requests: api.ResourceList{
								api.ResourceName(v1.ResourceCPU):    resource.MustParse("100m"),
								api.ResourceName(v1.ResourceMemory): resource.MustParse("100Mi"),
							},
							Limits: api.ResourceList{
								api.ResourceName(v1.ResourceCPU):    resource.MustParse("400m"),
								api.ResourceName(v1.ResourceMemory): resource.MustParse("400Mi"),
							},
						},
						VolumeMounts: []api.VolumeMount{
							{
								MountPath: "/var/lib/grafana",
								Name:      "grafana-storage",
								ReadOnly:  false,
							},
						},
					}},
					SecurityContext: &api.PodSecurityContext{
						RunAsNonRoot: &runAsNonRoot,
						RunAsUser:    &runAsUser,
					},
					ServiceAccountName: "grafana",
					Volumes: []api.Volume{
						{
							Name:         "grafana-storage",
							VolumeSource: api.VolumeSource{EmptyDir: &api.EmptyDirVolumeSource{}},
						},
					},
				},
			},
		},
	}

	y, err := yaml.Marshal(o)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return ""
	}
	return string(y)
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
    port: 80
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

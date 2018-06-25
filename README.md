# grafana-provisioner
The project is currently pre-alpha and it is expected that breaking changes to the API will be made in the upcoming releases.

### Overview 
Small comandline tool to manage N grafana deployments backed by a mysql database as a dashboard store for grafana. When running the tool will. 
* create a database for every team to be used by grafana
* provision a namespace, serviceaccount, ingress, service, deployemnt in k8s

if there is changes in the manifests the changes will be applied in the cluster to all teams in the config, or the single team passed by flag. 


### Problem statement
having N deployments of grafana managed will result is image version drift, having a solution for git sync like the current solution is for teams to manage deployments is not a good user experience. and maintaining using helm or equivalent will need to much tooling.


### Usage
```
alexander.brandstedt $ ./grafana-provisioner provision --help

Usage:
  grafana-provisioner provision [flags]

Flags:
  -d, --db.dns string           mysql database dns (default "localhost")
  -p, --db.pass string          mysql database password
  -P, --db.port int             mysql database port (default 3306)
  -u, --db.user string          mysql database username
  -D, --dry-run                 only output data
  -h, --help                    help for provision
  -i, --image string            grafana official container image (default "grafana/grafana:5.1.0")
  -I, --ingress.prefix string   dna prefix template %s.prefix (default "grafana.example.com")
  -k, --kube.config string      outside cluster path to kube config
  -t, --team string             team name
```

example command 
```
$ grafana-provisioner provision \
    --kube.config ~/.kube/config \
    --db.dns grafana.abcd.eu-west-1.rds.amazonaws.com \
    --db.pass abc123 \
    --db.user root \
    --image="grafana/grafana:5.1.4"

```

example config file
```
teams:
  - foo
  - bar
  - baz
```

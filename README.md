# grafana-provisioner


### Problem statement
having N deployments of grafana managed will result is image version drift, having a
solution for git sync like the current solution is for teams to manage deployments is
not a good user experience. and maintaining using helm or equivalent will need to much tooling.


### Solution
having a custom grafana provisioner that creates a database for every team and updating all deployment
at one time.


### Usage

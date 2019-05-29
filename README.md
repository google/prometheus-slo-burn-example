# Prometheus SLO Burn

This is home to example code for exposing SLIs using open source code in prometheus.

## Build Images

-   `$ gcloud builds submit --project $GOOGLE_PROJECT` in the root directory.
-   These images are currently published and publicly available from the project
    `cre-prometheus-slo-alerting`.

## Terraform setup

 - Install Terraform
 - Set up terraform env (assumes you have a working gcloud install and a google project):

```
$ [[ $CLOUD_SHELL ]] || gcloud auth application-default login
$ export GOOGLE_PROJECT=$(gcloud config get-value project)
$ export REGION=europe-west2
```

-   `$ cd terraform`
-   `$ terraform init` - installs terraform deps
-   `$ terraform apply -var "gcp_region=$REGION"` - Will ask you before it does
    anything. Will take ~10m to actually run. You can also run `terraform plan`
    to just get a dry run output.
-   `$ gcloud container clusters get-credentials example --region $REGION
    --project $GOOGLE_PROJECT` - Configures `kubectl` to work with the cluster
    you just created.
-   `$ kubectl create clusterrolebinding $USER-cluster-admin-binding
    --clusterrole=cluster-admin --user=$(gcloud config get-value account
    --project $GOOGLE_PROJECT)` - Gives your user permissions to create cluster
    role bindings that prometheus needs.
-   `$ kubectl apply -f ./k8s`

## Teardown

-   `$ cd terraform; terraform destroy -var "gcp_region=$REGION"`

## Running Locally

-   Start kubernetes (see
    https://kubernetes.io/docs/setup/pick-right-solution/#local-machine-solutions
    ).
-   Run `$ kubectl config current-context` to make sure you are in the correct
    context.
-   `$ cd terraform`
-   `$ kubectl apply -f ./k8s`
-   `$ kubectl get services --namespace=monitoring` you will see something like:

```
$ kubectl get services --namespace=monitoring
NAME            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
cloudprober     NodePort    10.104.187.119   <none>        8080:31589/TCP   21m
grafana         NodePort    10.104.206.150   <none>        8080:30431/TCP   21m
node-exporter   ClusterIP   None             <none>        9100/TCP         21m
prometheus      NodePort    10.101.58.210    <none>        9090:31517/TCP   21m
server          NodePort    10.111.115.243   <none>        8080:31796/TCP   21m
```

This means that now you can visit http://localhost:30431 and see the grafana
dashboard.

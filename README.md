# Service Catalog Templates

Cluster operators use Service Catalog Templates to define a default class, plan and parameters when provisioning a 
service of a particular type, such as “mysqldb”. This enables an application to specify a dependency on a type of 
service without requiring upfront knowledge of which broker it will be provisioned on, in addition to all the 
broker-specific parameters for a particular class and plan.

> “My application requires a mysql database.”
>
> “I just want to try out this application using the recommended defaults.”
> 
> “I need to encourage everyone to use the cheapest plans possible in our integration environment”

[Read the full proposal](https://docs.google.com/document/d/1vUxiOCKdnl47RKzeRgJ43_6eV2g95T6yM_m50Msrm3c)

This repository is a proof of concept for the above proposal. Our goal is to validate assumptions and technical 
decisions, and then contribute it upstream to Kubernetes.

# QuickStart

1. Create a cluster (v1.9+) with a service broker installed. The [Open Service Broker for Azure QuickStart on Minikube](https://github.com/Azure/open-service-broker-azure/blob/master/docs/quickstart-minikube.md)
    is a great guide to get up and running quickly. Currently supporting `v0.9.0-alpha` of the Open Service Broker for Azure.

1. Clone this repository and change to its directory:

    ```
    git clone https://github.com/Azure/service-catalog-templates.git
    cd service-catalog-templates
    ```

1. Install the Service Catalog Templates CLI, svcatt:

    ```
    go install ./cmd/svcatt
    ```

1. Install Service Catalog Templates on your cluster:

    ```
    helm install --name svcatt-crd --namespace svcatt charts/svcatt-crd --wait
    helm install --name svcatt-osba --namespace svcatt charts/svcatt-osba --wait
    helm install --name svcatt --namespace svcatt charts/svcatt --wait
    ```

1. Install the example Wordpress chart that takes advantage of Service Catalog Templates
    to provision a provider agnostic mysql database.
    
    ```
    helm install --name wordpress charts/wordpress --namespace svcatt
    ```
    
While the database provisions (it can take a while!), let's take a look at what 
the Wordpress chart needed in order to provision a database.
Rather than provisioning a ServiceInstance directly from the Service Catalog, the chart
defines a TemplatedInstance, requesting a service of type `mysqldb`. The chart
also defines a TemplatedBinding:

```yaml
apiVersion: templates.servicecatalog.k8s.io/experimental
kind: TemplatedInstance
metadata:
  name: wordpress-mysql-instance
spec:
  serviceType: mysqldb
---
apiVersion: templates.servicecatalog.k8s.io/experimental
kind: TemplatedBinding
metadata:
  name: wordpress-mysql-binding
spec:
  instanceRef:
    name: wordpress-mysql-instance
  secretName: wordpress-mysql-secret
```

From that the Templates controller, using the templates provided by the OSBA broker,
resolved a ServiceClass, ServicePlan and default parameters:

```yaml
apiVersion: templates.servicecatalog.k8s.io/experimental
kind: BrokerInstanceTemplate
metadata:
  name: default-mysqldb
  labels:
    serviceType: mysqldb
spec:
  serviceType: mysqldb
  clusterServiceClassExternalName: azure-mysql
  clusterServicePlanExternalName: basic50
  parameters:
    location: eastus
    resourceGroup: default
    sslEnforcement: disabled
    firewallRules:
    - startIPAddress: "0.0.0.0"
      endIPAddress: "255.255.255.255"
      name: "AllowAll"
---
apiVersion: templates.servicecatalog.k8s.io/experimental
kind: BrokerBindingTemplate
metadata:
  name: default-mysqldb
  labels:
    serviceType: mysqldb
  spec:
    serviceType: mysqldb
  # OSBA returns standard keys for the connection data so
  # this is pretty boring, if it didn't, the broker can
  # provide a mapping:
  # secretKeys:
  #   database-name: database
  #   server: host
  #   passwd: password
```

Using the OSBA broker template, the Templates controller created a corresponding ServiceInstance:

```console
$ svcatt get templated-instances -n svcatt
                 NAME                  NAMESPACE   SERVICE TYPE      CLASS       PLAN     STATUS
+------------------------------------+-----------+--------------+-------------+---------+--------+
  wordpress-wordpress-mysql-instance   svcatt      mysqldb        azure-mysql   basic50

$ svcatt get instances -n svcatt
                 NAME                  NAMESPACE      CLASS       PLAN        STATUS
+------------------------------------+-----------+-------------+---------+--------------+
  wordpress-wordpress-mysql-instance   svcatt      azure-mysql   basic50   Provisioning
```

\* NOTE: The status field on Templated resources is not implented yet.

Before we can proceed, wait until the instance is `Ready`:

```console
$ watch svcatt get instances -n svcatt
                 NAME                  NAMESPACE      CLASS       PLAN     STATUS
+------------------------------------+-----------+-------------+---------+--------+
  wordpress-wordpress-mysql-instance   svcatt      azure-mysql   basic50   Ready
```

After the instance is provisioned, Service Catalog creates a secret named "wordpress-mysql-secret-template"
containing the values returned from the broker:
 
```console
$ kubectl describe secret wordpress-wordpress-mysql-secret-template -n svcatt
Name:         wordpress-wordpress-mysql-secret-template
Namespace:    svcatt
Labels:       <none>
Annotations:  <none>

fype:  Opaque

Data
====
tags:         9 bytes
uri:          180 bytes
username:     47 bytes
database:     10 bytes
host:         61 bytes
password:     16 bytes
port:         4 bytes
sslRequired:  5 bytes
```
 
The Templates controller then applied the TemplatedBinding from the Wordpress chart,
creating the final secret for Wordpress to bind to:

```console
$ kubectl describe secret wordpress-wordpress-mysql-secret -n svcatt
Name:         wordpress-wordpress-mysql-secret
Namespace:    svcatt
Labels:       <none>
Annotations:  <none>

Type:  Opaque

Data
====
tags:         9 bytes
uri:          180 bytes
username:     47 bytes
database:     10 bytes
host:         61 bytes
password:     16 bytes
port:         4 bytes
sslRequired:  5 bytes
```

In this case, the keys are the same for both secrets. However if the broker returns non-standard keys,
for example "DatabaseName" instead of "database", the Templates controller would 
use the BrokerBindingTemplate to remap the keys to match the standard keys and save that
in final secret so that the Wordpress chart can rely upon a standard set of keys in the secret.

To prove that it all works, run the following command to open a web browser and view the Wordpress
site:

```
open http://$(minikube ip):$(kubectl get service wordpress-wordpress -n svcatt -o jsonpath={.spec.ports[?\(@.name==\"http\"\)].nodePort})
```

MAGIC! 🎩✨

Now let's clean up the resources created by this quickstart:

```console
$ svcatt unbind wordpress-wordpress-mysql-instance -n svcatt
deleted wordpress-wordpress-mysql-binding

$ svcatt get templated-bindings -n svcatt
  NAME   NAMESPACE   INSTANCE   STATUS
+------+-----------+----------+--------+

$ svcatt deprovision wordpress-wordpress-mysql-instance -n svcatt

$ svcatt get templated-instances -n svcatt
  NAME   NAMESPACE   SERVICE TYPE   CLASS   PLAN   STATUS
+------+-----------+--------------+-------+------+--------+

$ svcatt get instances -n svcatt
               NAME                  NAMESPACE      CLASS       PLAN         STATUS
+------------------------------------+-----------+-------------+---------+----------------+
wordpress-wordpress-mysql-instance   svcatt      azure-mysql   basic50   Deprovisioning
```

\* NOTE: The templated instance is deleted immediately because a finalizer has not yet been implemented.

# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.microsoft.com.

When you submit a pull request, a CLA-bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., label, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

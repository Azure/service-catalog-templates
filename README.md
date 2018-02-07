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

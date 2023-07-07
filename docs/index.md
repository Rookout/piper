# Introduction
![alt text](https://www.rookout.com/wp-content/uploads/2022/10/ArgoPipeline_1.0_Hero.png.webp?raw=true)

Welcome to Piper! Piper is open source project that aimed at providing multibranch pipeline functionality to Argo Workflows, allows users to create distinct Workflows based on Git branches.

## General explanation

Piper configures a webhook in git provider and listens to the webhooks sends. It will create a Workflow CRD out of branches that contains `.workflows` folder. This folder should contain declarations of the templates and main DAG that will be running. Finally, it will submit the Workflow as a K8s resource in the cluster. To access more detailed explanations, please navigate to the [Usage](docs/usage.md).



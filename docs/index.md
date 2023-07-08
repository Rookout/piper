# Introduction
![alt text](https://www.rookout.com/wp-content/uploads/2022/10/ArgoPipeline_1.0_Hero.png.webp?raw=true)

Welcome to Piper! 

Piper is an open source project that aimed at providing multibranch pipeline functionality to Argo Workflows, allows users to create distinct Workflows based on Git branches.

## General explanation

![alt text](https://raw.githubusercontent.com/Rookout/piper/main/docs/img/flow.svg)

To achieve multibranch pipeline functionality Piper will do the hard works for us.
At initialization, it will load all configuration and create a webhook in repository or organization scope.
Then each branch that have `.workflows` folder will create a Workflow CRD out of the files in this folder.

To learn more go to, [Use piper](usage.md).
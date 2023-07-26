# Introduction

<p align="center">
  <img src="https://www.rookout.com/wp-content/uploads/2022/10/ArgoPipeline_1.0_Hero.png.webp?raw=true" />
</p>

Welcome to Piper! 

Piper is an open source project that aimed at providing multibranch pipeline functionality to Argo Workflows, allows users to create distinct Workflows based on Git branches.

## General explanation

<p align="center">
  <img src="https://raw.githubusercontent.com/Rookout/piper/main/docs/img/flow.svg" />
</p>

To achieve multibranch pipeline functionality Piper will do the hard works for us.
At initialization, it will load all configuration and create a webhook in repository or organization scope.
Then each branch that have `.workflows` folder will create a Workflow CRD out of the files in this folder.

![type:video](./img/piper-demo-1080.mp4)
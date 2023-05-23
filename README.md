# Piper
Welcome to Piper! Piper is open source project that aimed at providing multibranch pipeline functionality to Argo Workflows, allows users to create distinct Workflows based on Git branches.

## Table of Contents

- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Reporting Issues](#reporting-issues)
- [Pull Requests](#pull-requests)
- [Coding Guidelines](#coding-guidelines)
- [License](#license)

## Concept

Piper configures a weebhook in git provider and listens to the webhooks sends. It will create a Workflow CRD out of branches that contains `.workflows` folder. This folder should contain delclerations of the templates and main DAG that will be running. Finally it will submit the Workflow as a K8s resource in the cluster.
## Usage

Piper should be deployed in the cluster with Argo Workflows. 

```
helm repo add piper https://piper.rookout.com
helm install piper piper/piper --namespace workflows
```

## Roadmap
1. Create Github provider handler.
2. Create webhook server. 
3. Create helm chart.
3. Implmentation of Workflow creation by the example.
4. debug pause injection implmentation feature (will fail the pipeline).


## How to Contribute

If you're interested in contributing to this project, please feel free to submit a pull request. We welcome all contributions and feedback.
Please checkout our [Contribution guidelines for this project](docs/CONTRIBUTING.md)

## License

This project is licensed under the Apache License. Please see the [LICENSE](LICENSE) file for details.




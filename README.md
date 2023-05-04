# ArgoWorkflows Multibranch Pipeline
This project aimed at providing Multibranch pipeline functionality within the Argo stack of CI/CD tools. This project allows users to create distinct pipelines based on Git branches, making workflow organization and maintenance more efficient.

## Concept



## Usage

The seeder is executed via the command line. The following parameters are required:

```python seeder.py <path to template> <path to .workflows directory> <event type> <branch>``

<path to template>: Path to the seeder-workflow-template.yaml file.
<path to .workflows directory>: Path to the directory containing the .yaml DAG files.
<event type>: The type of the event triggering the workflow (e.g., push, pull_request, merge).
<branch>: The branch destintion branch to merge into.


## Examples

This repository includes an examples folder with a `seeder-workflow-template.yaml` file, a `workflow.yaml` file, and a `.workflows` directory with various DAG files.

To generate a new workflow based on the provided example, navigate to the root of the repository and run:

```python seeder.py examples/seeder-workflow-template.yaml examples/.workflows push main```

This will generate a new workflow YAML file (`workflow.yaml`) in the root of the repository, based on the `seeder-workflow-template.yaml` file and the `.yaml` DAG files in the `.workflows` directory.

## Roadmap
1. Create a microservice that listen to external webhooks
2. Ability to submit directly to ArgoWorkflows server

## Contributing

If you're interested in contributing to this project, please feel free to submit a pull request. We welcome all contributions and feedback.

## License

This project is licensed under the Apache License. Please see the LICENSE file for details.




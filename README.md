# Argo Piper
This project aimed at providing Multibranch pipeline functionality to Argo Workflows. This project allows users to create distinct pipelines based on Git branches, making workflow organization and maintenance more efficient.

## Concept

![alt text](https://github.com/Rookout/argo-workflows-multibranch-pipeline/blob/main/docs/seeder-pipeline.png?raw=true)

This script will be running within Seeder Workflow, after genrate inside Workflow by ArgoEvents Sensor. You can template the Seeder Workflow with metadata from the webhook and then pass it to the Multi Brnach Workflow.
## Usage

The seeder is executed via the command line. The following parameters are required:

```
python seeder.py <path to template> <path to .workflows directory> <event type> <branch>
```

<path to template>: Path to the seeder-workflow-template.yaml file.
<path to .workflows directory>: Path to the directory containing the .yaml DAG files.
<event type>: The type of the event triggering the workflow (e.g., push, pull_request, merge).
<branch>: The branch destintion branch to merge into.


## Examples

This repository includes an examples folder with a `seeder-workflow-template.yaml` file, a `workflow.yaml` file, and a `.workflows` directory with various DAG files.

To generate a new workflow based on the provided example, navigate to the root of the repository and run:

```
python seeder.py examples/seeder-workflow-template.yaml examples/.workflows push main
```

This will generate a new workflow YAML file (`workflow.yaml`) in the root of the repository, based on the `seeder-workflow-template.yaml` file and the `.yaml` DAG files in the `.workflows` directory.

`main.yaml` is the default pipeline to run if '<event-type>-<branch>.yaml' don't exists.

## Roadmap
1. Debug interface
2. Multiple template.yaml files
3. Create a microservice that listen to external webhooks
4. Ability to submit directly to ArgoWorkflows server

## Contributing

If you're interested in contributing to this project, please feel free to submit a pull request. We welcome all contributions and feedback.

## License

This project is licensed under the Apache License. Please see the LICENSE file for details.




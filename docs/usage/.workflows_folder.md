## .workflows Folder

Piper will look in each of the target branches for a `.workflows` folder. [example](https://github.com/Rookout/piper/tree/main/examples/.workflows).
We will explain each of the files that should be included in the `.workflows` folder.

### triggers.yaml (convention name)

This file holds a list of triggers that will be executed `onStart` by `events` from specific `branches`. 
Piper will execute each of matching triggers, so configure it wisely.
```yaml
- events:
    - push
    - pull_request.synchronize
  branches: ["main"]
  onStart: ["main.yaml"]
  onExit: ["exit.yaml"]
  templates: ["templates.yaml"]
  config: "default"
```
Can be found [here](https://github.com/Rookout/piper/tree/main/examples/.workflows/triggers.yaml).

In this example `main.yaml` will be executed as DAG when `push` or `pull_request.synchronize` events will be applied in `main` branch.
`onExit` will be executed `exit.yaml` when finished the workflow as exit handler.


`onExit` can overwrite the default `onExit` configuration from by reference existing DAG tasks as in the [example](https://github.com/Rookout/piper/tree/main/examples/.workflows/exit.yaml).

`config` field used for workflow configuration selection. the default value is `default` configuration.

#### events
Events field used to terminate when the trigger will be executed. name of the event depends on the git provider. 

For instance, GitHub pull_request event have few action, one of them is synchronize.

#### branches
For which branch that trigger will be executed.

####  onStart
This [file](https://github.com/Rookout/piper/tree/main/examples/.workflows/main.yaml) can be named as you wish and will be referenced in `triggers.yaml` file. It will define an entrypoint DAG that the Workflow will execute.

As a best practice, this file should contain the dependencies logic and parametrization of each of referenced templates. It should not implement new templates, for this, use template.yaml file.

####  onExit
This field used to pass verbose exitHandler to the triggered workflow.
It will override the default onExit from the provided `config` or the default `config`.

In the provided `exit.yaml` describes a DAG that will overwrite the default `onExit` configuration.
[Example](https://github.com/Rookout/piper/tree/main/examples/.workflows/exit.yaml)

####  templates
This field will have additional templates that will be injected to the workflows. 
The purpose of this field is to create repository scope templates that can be referenced from the DAGs templates at `onStart` or `onExit`.
[Example](https://github.com/Rookout/piper/tree/main/examples/.workflows/templates.yaml)

As a best practice, use this field for template implementation and reference them from executed.
[Example](https://github.com/Rookout/piper/tree/main/examples/.workflows/main.yaml).

### config
configured by `piper-workflows-config` [configMap](docs/usage/workflws_config.md). 
Can be passed explicitly, or will use `deafault` configuration.

###  parameters.yaml (convention name)
Will hold a list of global parameters of the Workflow. 
can be referenced from any template with `{{ workflow.parameters.___ }}.

[Example](https://github.com/Rookout/piper/tree/main/examples/.workflows/parameters.yaml)
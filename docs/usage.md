## Files and Directory Convention

Piper will look in each of the target branches a `.workflows` [folder](https://github.com/Rookout/piper/tree/main/examples/.workflows). This folder should contain the following files to create a Workflow out of branch configuration:

### triggers.yaml (convention name)

This file holds a list of triggers that will be executed `onStart` by `events` from specific `branch`. In this [example](https://github.com/Rookout/piper/tree/main/examples/.workflows/triggers.yaml), `main.yaml` will be executed as DAG when `push` or `pull_request.synchronize` events will be applied in `main` branch. 

`onExit` can overwrite the default `onExit` configuration from `piper-workflows-config` by reference existing DAG tasks as [exit.yaml](https://github.com/Rookout/piper/tree/main/examples/.workflows/exit.yaml).

`config` field used for workflow configuration selection. the default value is `default` configuration.

###  main.yaml or others 

This [file](https://github.com/Rookout/piper/tree/main/examples/.workflows/main.yaml) can be named as you wish and will be referenced in `triggers.yaml` file. It will define an entrypoint DAG that the Workflow will execute.

As a best practice, this file should contain the dependencies logic and parametrization of each of referenced templates. It should not implement new templates, for this, use template.yaml file.

###  exit.yaml or others

This [file](https://github.com/Rookout/piper/tree/main/examples/.workflows/exit.yaml) describes a DAG that will overwrite the default `onExit` configuration from `piper-workflows-config`.

###  templates.yaml or others

This [file](https://github.com/Rookout/piper/tree/main/examples/.workflows/templates.yaml) will have additional templates that will be injected to the workflows. The purpose of this file is to implemented repository scope templates that can be referenced in the executed file.

As a best practice, use this file as implementation of template and reference them from executed [file](https://github.com/Rookout/piper/tree/main/examples/.workflows/main.yaml).

###  parameters.yaml (convention name)

This [file](https://github.com/Rookout/piper/tree/main/examples/.workflows/parameters.yaml) will hold a list of global parameters of the Workflow. can be referenced from any template with `{{ workflow.parameters.___ }}
Also piper provided global parameters as followed:

1. {{ workflow.parameters.event }} the event that triggered the workflow.

2. {{ workflow.parameters.action }} }} the action that triggered the workflow.

3. {{ workflow.parameters.dest_branch }} the destination branch for pull request.

4. {{ workflow.parameters.commit }} the commit that triggered the workflow.

5. {{ workflow.parameters.repo }} repository name that triggered the workflow.

6. {{ workflow.parameters.user }} the username that triggered the workflow.

7. {{ workflow.parameters.user_email }} the user's email that triggered the workflow.

8. {{ workflow.parameters.pull_request_url }} the url of the pull request that triggered the workflow.

9. {{workflow.parameters.pull_request_title }} the tile of the pull request that triggered the workflow.


## Workflow Configuration (Spec)

Best to configure using helm chart in `piper.workflowsConfig` parameter.
To support Workflow configuration (defining Workflow spec field) as presented in the [examples](https://github.com/Rookout/piper/tree/main/examples/config.yaml), Piper consumes a configMap named `piper-workflows-config`. This config map can have `default` Workflow spec, that will be used for any Workflow created or, create other configuration sets that have to be explicitly called on each of the [triggers](https://github.com/Rookout/piper/tree/main/examples/.workflows/triggers.yaml) (`config` field). Please notice that the fields `onStart` and `onExit` should not exist. Instead, `onStart` is a managed field, and `onExit` can configure a default DAG to execute when the workflow finishes.  

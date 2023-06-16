## Files and Directory Convention

Piper will look in each of the target branches a `.workflows` [folder](../examples/.workflows). This folder should contain the following files to create a Workflow out of branch configuration:

### triggers.yaml

This file holds a list of triggers to that will be executed `onStart` by `event` from specific `branch`. In this [example](../examples/.workflows/triggers.yaml), `main.yaml` will be executed as DAG when `push` or `merge` events will be applied in `main` branch. 

`onExit` can overwrite the default `onExit` configuration from `piper-workflows-config` by reference existing DAG tasks as [exit.yaml](../examples/.workflows/exit.yaml).

###  main.yaml or others 

This [file](../examples/.workflows/main.yaml) can be named as you wish and will be referenced in `triggers.yaml` file. It will define an entrypoint DAG that the Workflow will execute.

As a best practice, this file should contain the dependencies logic and parametrization of each of referenced templates. It should not implement new templates, for this, use template.yaml file.

###  exit.yaml or others

This [file](../examples/.workflows/exit.yaml) describes a DAG that will overwrite the default `onExit` configuration from `piper-workflows-config`.

###  templates.yaml or others

This [file](../examples/.workflows/templates.yaml) will have additional templates that will be injected to the workflows. The purpose of this file is to implemented repository scope templates that can be referenced in the executed file.

As a best practice, use this file as implementation of template and reference them from executed [file](../examples/.workflows/main.yaml).

###  parameters.yaml

This [file](../examples/.workflows/parameters.yaml) will hold a list of global parameters of the Workflow. can be referenced from any template with `{{ workflow.parameters.___ }}


## Workflow Configuration (Spec)

Best to configure using helm chart in `piper.workflowsConfig` parameter.
To support Workflow configuration (defining Workflow spec field) as presented in the [examples](../examples/config.yaml), Piper consumes a configMap named `piper-workflows-config`. This config map can have `default` Workflow spec, that will be used for any Workflow created or, create other configuration sets that have to be explicitly called on each of the [triggers](../examples/.workflows/triggers.yaml) (`config` field). Please notice that the fields `entrypoint` and `onExit` should not exist. Instead, `entrypoint` is a managed field, and `onExit` can configure a default DAG to execute when the workflow finishes.  



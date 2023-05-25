## Files and Directory Convention

Piper will look in each of the target branches a `.workflows` [folder](../examples/.workflows). This folder should contain the following files to create a Workflow out of branch configuration:

### triggers.yaml

This files will hold a list of triggers to `execute` by `event` from specifc `branch`. In this [example](../examples/.workflows/triggers.yaml), `main.yaml` will be executed as DAG when `push` or `merge` events will be applied in `main` branch. 

`onExit` can overwrite the default `onExit` configuration from `piper-workflows-config` by refrence existing DAG instuctions as [exit.yaml](../examples/.workflows/exit.yaml).

###  exit.yaml

This [file](../examples/.workflows/exit.yaml) describes a DAG that will overwrite the default `onExit` configuration from `piper-workflows-config`.

###  paramters.yaml

This [file](../examples/.workflows/paramters.yaml) will hold a list of global parameters of the Workflow. can be refrenced from any template with `{{ workflow.paramters.___ }}

###  templates.yaml

This [file](../examples/.workflows/templates.yaml) will have additional templates that will be injected to the workflows. The purpose of this file is to implmented repository scope templates that can be referenced in the executed file. 

As a best practice, use this file as implmentation of template and reference them from executed [file](../examples/.workflows/main.yaml).

###  main.yaml

This [file](../examples/.workflows/main.yaml) can be named as you wish and will be referenced in `triggers.yaml` file. It will define a entrypoint DAG that the Workflow will execute. 

As a best practice, this file should contain the depencies logic and parametrizations of each of referenced templates. It should not implment new templates, for this, use template.yaml file.

## Workflow Configuration (Spec)

To support Workflow configuration (defining Workflow spec field) as prestened in the [examples](../examples/config.yaml), Piper consumes a configMap named `piper-workflows-config`. This config map can have `default` Workflow spec, that will be used for any Workflow created or, create other configuration sets that have to be explicity called on each of the [triggers](../examples/.workflows/triggers.yaml) (`config` field). Plesase notice that the fields `entrypoint` and `onExit` should not exists. Instead, `entrypoint` is a managed field, and `onExit` can configure a default DAG to execute when the workflow finish.  



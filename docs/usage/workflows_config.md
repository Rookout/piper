## Workflow Configuration

Piper can inject configuration for Workflows that Piper creates.

`default` config used as a convention for all Workflows that piper will create, even if not explicitly mentioned in triggers.yaml file.

### ConfigMap
Piper will mount a configMap when helm used.
`piper.workflowsConfig` variable in helm chart, will create a configMap that hold set of configuration for Piper.
Here is an [examples](https://github.com/Rookout/piper/tree/main/examples/config.yaml) of such configuration.

### Spec
This will be injected to Workflow spec field. can hold all configuration of the Workflow. 
> :warning: Please notice that the fields `entrypoint` and `onExit` should not exist in spec. both of them are managed fields.

### onExit
This is the exit handler for each of the Workflows create by piper. 
It configures a DAG that will be executed when the workflow ends. 
You can provide the templates to it us in the following [Examples](https://github.com/Rookout/piper/tree/main/examples/config.yaml).
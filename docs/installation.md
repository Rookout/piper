## Instalation

Piper should be deployed in the cluster with Argo Workflows. Piper will create a CRD that Argo Workflows will pick, so install or configure Piper to create those CRDs in the right namespace. 

Please check out values.yaml file of the helm chart.

To add piper helm repo run:
```bash
helm repo add piper https://piper.rookout.com
```

After configuring Piper value.yaml, run the following command for installation:
```bash
helm install piper piper/piper
```

## Configuration

### Ingress

Piper should listen to webhooks from your git provider. Checkout value.yaml for `ingress`
### Git Token

The git token should be passed as secret in the helm chart at `gitProvider.github.token`

### Webhook creation

Piper will create a webhook configuration for you, for the whole orgnization or for each repo you configure.
First configure which git provider you are using `gitProvider.name` (Now only supports GitHub)

For organization level configuration provide the following value `gitProvider.organization.name` and `gitProvider.webhook.orgLevel` to `true`.

For granular repo configuration provide `gitProvider.organization.name` and `gitProvider.webhook.repoList`. 

When Piper will be deleted the finalizer should delete the configured webhooks for your git provider. (On development)

### Argo Workflow Server

Piper will use REST API to communicate with Argo Workflows server for linting or for creation of workflows (ARGO_WORKFLOWS_CREATE_CRD). Please follow this [configuration](https://argoproj.github.io/argo-workflows/rest-api/).

To lint the workflow before submitting it, please configure the internal address of Argo Workflows server (for example, `argo-server.workflows.svc.cluster.local`) in the field: `argoWorkflows.server.address`. Argo will need a [token](https://argoproj.github.io/argo-workflows/access-token/) to authenticate. please provide the secret in `argoWorkflows.server.token`, Better to pass as a refrences to a secret in the field `argoWorkflows.server.token`.

### Skip CRD Creation

Piper can communicate directly to Argo Workflow using ARGO_WORKFLOWS_CREATE_CRD environment variable, if you want to skip the creation of CRD change `argoWorkflows.crdCreation` to `false`.
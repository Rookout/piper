## Instalation

Piper should be deployed in the cluster with Argo Workflows. Piper will create a CRD that Argo Workflows will pick, so install or configure Piper to create those CRDs in the right namespace. 

Please checkout values.yaml file of the helm chart.

To add piper helm repo run:
```bash
helm repo add piper https://piper.rookout.com
```

After configuring Piper, run the following command for installtion:
```bash
helm install piper piper/piper --namespace workflows
```

## Configuration

### Ingress

Piper should listen to webhooks from your git provider. Checkout value.yaml for `ingress`
### Git Token

The git token should be passed as secret in the helm chart at `gitProvider.github.token`

### Webhook creation

Piper will create a webhook configuration for you, for the whole orgnization or for each repo you configure.
For orgnization configuration provide the following value `gitProvider.organization.name` and `gitProvider.webhook.org` to `true`.

For granular repo configuration provide `gitProvider.organization.name` and `gitProvider.webhook.repoList`. 

When Piper will be deleted the finilizer should deleted the configured webhooks for your git provider.

### Argo Workflow Server

Piper will use REST API to communicate with Argo Workflows server. Please follow this [configuation](https://argoproj.github.io/argo-workflows/rest-api/).

To lint the workflow before submiting it, please configure the internal address of Argo Workflows server (for example, `argo-server.workflows.svc.cluster.local`) in the field: `argoWorkflows.server.address`. Argo will need a [token](https://argoproj.github.io/argo-workflows/access-token/) to authenticate. please provide a refrences to a secret in the field `argoWorkflows.server.token`.

### Skip CRD Creation

Piper can communicate directly to Argo Workflow, if you want to skip the creation of CRD change `argoWorkflows.crdCreation` to `false`.
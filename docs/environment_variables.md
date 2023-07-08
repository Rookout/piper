## Environment Variables

The environment variables used by Piper to configure its functionality.
The helm chart populate them using [values.yaml](https://github.com/Rookout/piper/tree/main/helm-chart/values.yaml) file

### Git

* GIT_PROVIDER
  The git provider that Piper will use, possible variables: GitHub (will support bitbucket and gitlab)

* GIT_TOKEN
  The git token that will be used.

* GIT_ORG_NAME
  The organization name.

* GIT_ORG_LEVEL_WEBHOOK
  Boolean variable, whether to config webhook in organization level. default `false`

* GIT_WEBHOOK_REPO_LIST
  List of repositories to configure webhooks to.

* GIT_WEBHOOK_URL
  URL of piper ingress, to configure webhooks.

### Argo Workflows Server
* ARGO_WORKFLOWS_TOKEN
  The token of Argo Workflows server.

* ARGO_WORKFLOWS_ADDRESS
  The address of Argo Workflows Server.
  
* ARGO_WORKFLOWS_CREATE_CRD
  Whether to directly send Workflows instruction or create a CRD in the Cluster.

* ARGO_WORKFLOWS_NAMESPACE
  The namespace of Workflows creation for Argo Workflows.

* KUBE_CONFIG
  Used to configure Argo Workflows client with local kube configurations.

### Rookout
* ROOKOUT_TOKEN
  The token used to configure Rookout agent. If not provided, will not start the agent.
* ROOKOUT_LABELS 
  The labels to label instances at Rookout, default to "service:piper"
* ROOKOUT_REMOTE_ORIGIN
  The repo URL for source code fetching, default:"https://github.com/Rookout/piper.git".
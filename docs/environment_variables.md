## Environment Variables

* GIT_PROVIDER
  The git provider that Piper will use, possible variables: github (will support bitbucket and gitlab)

* GIT_ORG_NAME
  The organization name.

* GIT_ORG_LEVEL_WEBHOOK
  Boolean variable, whether to config webhook in organization level. default `false`

* GIT_WEBHOOK_REPO_LIST
  Comma separated list of repositories to configure webhooks to.

* ARGO_WORKFLOWS_TOKEN
  The token of Argo Workflows server.

* ARGO_WORKFLOWS_ADDRESS
  The address of Argo Workflows Server.
  
* ARGO_WORKFLOWS_CREATE_CRD
  Whether to directly send Workflows instruction or create a CRD in the Cluster.

* ARGO_WORKFLOWS_NAMESPACE
  The namespace of Workflows creation for Argo Workflows.

* KUBE_CONFIG
  Path to local kubernetes configuration
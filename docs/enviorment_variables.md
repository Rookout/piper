## Enviorment Vraiables

* GIT_PROVIDER
  The git provider that Piper will use, posiable variables: github (will support bitbucket and gitlab)

* GIT_ORG_NAME
  The organization name.

* GIT_ORG_LEVEL_WEBHOOK
  Boolean vriable, wheter to config webhook in organization level. default `false`

* WEBHOOK_REPO_LIST
  Comma seperate list of repositories to configure webhooks to.

* ARGO_WORKFLOWS_TOKEN
  The token of Argo Workflows server.

* ARGO_WORKFLOWS_ADDRESS
  The address of Argo Workflows Server.
  
* ARGO_WORKFLOWS_CREATE_CRD
  Wheter to directly send Workflows instruction or create a CRD in the Clutser.
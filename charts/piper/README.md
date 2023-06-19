# piper

![Version: 1.0.1](https://img.shields.io/badge/Version-1.0.1-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

A Helm chart for Piper

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Assign custom [affinity] rules to the deployment |
| autoscaling.enabled | bool | `false` | Wheter to enable auto-scaling of piper. |
| autoscaling.maxReplicas | int | `5` | Maximum reoplicas of Piper. |
| autoscaling.minReplicas | int | `1` | Minimum reoplicas of Piper. |
| autoscaling.targetCPUUtilizationPercentage | int | `85` | CPU utilization percentage threshold. |
| autoscaling.targetMemoryUtilizationPercentage | int | `85` | Memory utilization percentage threshold. |
| env | list | `[]` | Additional environment variables for Piper. A list of name/value maps. |
| extraLabels | object | `{}` | Deployment and pods extra labels |
| fullnameOverride | string | `""` | String to fully override "piper.fullname" template |
| image.name | string | `"piper"` | Piper image name |
| image.pullPolicy | string | `"IfNotPresent"` | Piper image pull policy |
| image.repository | string | `"rookout"` | Piper public dockerhub repo |
| image.tag | string | `"latest"` | Piper image tag |
| imagePullSecrets | list | `[]` | secret to use for image pulling |
| ingress.annotations | object | `{}` | Piper ingress annotations |
| ingress.className | string | `""` | Piper ingress class name |
| ingress.enabled | bool | `false` | Enable Piper ingress support |
| ingress.hosts | list | `[{"host":"piper.example.local","paths":[{"path":"/","pathType":"ImplementationSpecific"}]}]` | Piper ingress hosts # Hostnames must be provided if Ingress is enabled. |
| ingress.tls | list | `[]` | Controller ingress tls |
| lifecycle | object | `{}` | Specify postStart and preStop lifecycle hooks for Piper container |
| nameOverride | string | `""` | String to partially override "piper.fullname" template |
| nodeSelector | object | `{}` | [Node selector] |
| piper.argoWorkflows.crdCreation | bool | `true` | Whether create Workflow CRD or send direct commands to Argo Workflows server. |
| piper.argoWorkflows.server.address | string | `""` | The DNS address of Argo Workflow server that Piper can address. |
| piper.argoWorkflows.server.existingSecret | string | `nil` |  |
| piper.argoWorkflows.server.namespace | string | `""` | The namespace in which the Workflow CRD will be created. |
| piper.argoWorkflows.server.token | string | `""` | This will create a secret named <RELEASE_NAME>-token and with the key 'token' |
| piper.gitProvider.existingSecret | string | `nil` |  |
| piper.gitProvider.name | string | `"github"` | Name of your git provider (github/gitlab/bitbucket). for now, only github supported. |
| piper.gitProvider.organization.name | string | `""` | Name of your Git Organization |
| piper.gitProvider.token | string | `nil` | This will create a secret named <RELEASE_NAME>-git-token and with the key 'token' |
| piper.gitProvider.webhook.existingSecret | string | `nil` |  |
| piper.gitProvider.webhook.orgLevel | bool | `false` | Whether config webhook on org level |
| piper.gitProvider.webhook.repoList | list | `[]` | Used of orgLevel=false, to configure webhook for each of the repos provided. |
| piper.gitProvider.webhook.secret | string | `""` | This will create a secret named <RELEASE_NAME>-webhook-secret and with the key 'secret' |
| piper.gitProvider.webhook.url | string | `""` | The url in which piper listens for webhook, the path should be /webhook |
| piper.workflowsConfig | object | `{}` |  |
| podAnnotations | object | `{}` | Annotations to be added to the Piper pods |
| podSecurityContext | object | `{"fsGroup":1001,"runAsGroup":1001,"runAsUser":1001}` | Security Context to set on the pod level |
| replicaCount | int | `1` | Piper number of replicas |
| resources | object | `{"requests":{"cpu":"200m","memory":"512Mi"}}` | Resource limits and requests for the pods. |
| rookout.existingSecret | string | `""` |  |
| rookout.token | string | `""` | Rookout token for agent configuration and enablement. |
| securityContext | object | `{"capabilities":{"drop":["ALL"]},"readOnlyRootFilesystem":true,"runAsNonRoot":true,"runAsUser":1001}` | Security Context to set on the container level |
| service.annotations | object | `{}` | Piper service extra annotations |
| service.labels | object | `{}` | Piper service extra labels |
| service.port | int | `80` | Service port For TLS mode change the port to 443 |
| service.type | string | `"ClusterIP"` | Sets the type of the Service |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| tolerations | list | `[]` | [Tolerations] for use with node taints |
| volumeMounts | list | `[]` | Volumes to mount to Piper container. |
| volumes | list | `[]` | Volumes of Piper Pod. |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)

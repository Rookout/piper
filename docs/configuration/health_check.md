## Health Check

Health check executed every 1 minute as configured in the helm chart under `livenessProbe`, and triggered by `/healthz` endpoint:
```yaml
livenessProbe:
  httpGet:
    path: /healthz
    port: 8080
    scheme: HTTP
  initialDelaySeconds: 10
  timeoutSeconds: 10
  periodSeconds: 60
  successThreshold: 1
  failureThreshold: 4
```

The mechanism for checking the health of Piper is:

1. Piper set health status of all webhooks to not-healthy.

2. Piper requests ping from all the webhooks configured. 

3. Git Provider send ping to `/webhook` endpoint, this will set the health status to `healthy` with timeout of 5 seconds.

4. Piper check the status of all webhooks configured.

Therefore, the criteria for health checking are:
1. The registered webhook exists. 
2. The webhook send a ping in 5 seconds.



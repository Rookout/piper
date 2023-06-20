## Roadmap

## Piper
- [x] Create webhook server. 
- [x] Github provider handler.
- [x] Implementation of Workflow creation by the example.
- [x] onExit overwrite in triggers.yaml.
- [x] HPA support.
- [x] Run as non-root user.
- [ ] Argo Workflows Server linter and submitter.
- [ ] Workflow status management service.
- [ ] Graceful shutdown.
- [ ] Logger.
- [ ] debug pause injection implementation feature (will fail the pipeline).
- [ ] triggers.yaml config selection.
- [ ] Label management.
- [ ] Add tag event to github webhook handler.

## General
- [x] gh pages for docs.
- [x] gh pages for chart.
- [ ] Test suite - In progress.
- [ ] e2e tests - In progress.

## CI
- [x] Dockerfile .
- [x] Local deployment using kind.
- [x] Multi arch build.
- [x] Helm chart - lint, package, version, publish to gh pages.
- [x] Application - test, build, version, publish.
- [x] Administration - DOC, PR title linter, branch permissions, change log, PR template.
- [ ] Administration - code coverage
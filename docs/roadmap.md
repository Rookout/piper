## Roadmap

## Piper
- [x] Create webhook server. 
- [x] Github provider handler.
- [x] Implementation of Workflow creation by the example.
- [ ] Argo Workflows Server linter and submitter.
- [ ] Workflow status management service.
- [ ] Graceful shutdown.
- [ ] Logger.
- [ ] debug pause injection implementation feature (will fail the pipeline).
- [x] onExit overwrite in triggers.yaml.
- [ ] triggers.yaml config selection.
- [x] HPA support.
- [x] Run as non-root user.
- [ ] Label management.

## General
- [ ] Test suite.
- [ ] e2e tests.
- [ ] gh pages for docs and chart.

## CI
- [x] Dockerfile .
- [x] Local deployment using kind.
- [x] Multi arch build.
- [x] Helm chart - lint, package, version, publish to gh pages.
- [x] Application - test, build, version, publish.
- [x] Administration - DOC, PR title linter, branch permissions, change log, PR template.
- [ ] Administration - code coverage
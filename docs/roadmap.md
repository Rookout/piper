## Roadmap

## Piper
- [x] Create webhook server. 
- [x] Github provider handler.
- [ ] Implementation of Workflow creation by the example.
- [ ] Argo Workflows Server linter and submitter.
- [ ] Workflow status management service.

### Future Piper Features
- [ ] debug pause injection implementation feature (will fail the pipeline).
- [x] onExit overwrite in triggers.yaml
- [ ] triggers.yaml config selection.
- [x] HPA support
- [ ] Run as non-root user.

## General
- [ ] Test suite
- [ ] e2e tests
- [ ] gh pages for docs.

## CI
- [ ] Dockerfile 
- [ ] Local deployment using kind
- [ ] Multi arch build
- [ ] Helm chart - lint, package, version, publish to artifact hub.
- [ ] Application - test, build, version, publish.
- [ ] Administration - DOC, PR title linter, branch permissions, change log, code coverage, PR template.
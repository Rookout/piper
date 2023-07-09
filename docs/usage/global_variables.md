## Global variables

Piper will automatically add Workflow scope parameters that can be referenced from any template.
The parameters taken from webhook metadata, and will be populated respectively to GitProvider and event that triggered the workflow.

1. `{{ workflow.parameters.event }}` the event that triggered the workflow.

2. `{{ workflow.parameters.action }}` the action that triggered the workflow.

3. `{{ workflow.parameters.dest_branch }}` the destination branch for pull request.

4. `{{ workflow.parameters.commit }}` the commit that triggered the workflow.

5. `{{ workflow.parameters.repo }}` repository name that triggered the workflow.

6. `{{ workflow.parameters.user }}` the username that triggered the workflow.

7. `{{ workflow.parameters.user_email }}` the user's email that triggered the workflow.

8. `{{ workflow.parameters.pull_request_url }}` the url of the pull request that triggered the workflow.

9. `{{workflow.parameters.pull_request_title }}` the tile of the pull request that triggered the workflow.

10. `{{workflow.parameters.pull_request_labels }}` comma seperated labels of the pull request that triggered the workflow.

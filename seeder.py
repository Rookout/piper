from sys import argv
import yaml

class Workflow():
  def __init__(self, wf_template_path):
    with open(wf_template_path) as f:
      self.wf = yaml.safe_load(f)

  def inject_templates(self, templates_path):
    try:
      with open(templates_path) as f:
        additional_templates = yaml.safe_load(f)
      self.wf["spec"]["templates"] += additional_templates
    except FileNotFoundError:
      print("didn't find any template.yaml file")

  def inject_dag(self, dag_path):
    with open(dag_path) as f:
      dag_flow = yaml.safe_load(f)
    entrypoint = self.wf["spec"]["entrypoint"]
    dag_flow = [{"name": entrypoint, "dag": {"tasks": dag_flow }}]
    self.wf["spec"]["templates"] += dag_flow

  def inject_name(self, name):
    self.wf["metadata"]["generateName"] = name + "-"
  
  def inject_parameters(self, param_path):
    with open(param_path) as f:
      parameters = yaml.safe_load(f)
    self.wf["spec"]["arguments"]["parameters"] = parameters

  def dump(self):
    with open("workflow.yaml", "w") as f:
      yaml.safe_dump(self.wf, f)    


def seed(event_type, branch, seeder_workflow_template_path, workflows_dir):
  workflow = Workflow(seeder_workflow_template_path)
  workflow.inject_templates(workflows_dir + "/templates.yaml")
  try:
    workflow.inject_dag(workflows_dir + "/main.yaml")
  except FileNotFoundError:
    workflow.inject_dag(workflows_dir + "/" + event_type + branch + ".yaml")
  workflow.name()
  
  return workflow

if __name__ == "__main__":
  '''
  python seeder.py merge main ./seeder-workflow-template.yaml ./.workflows 
  python seeder.py pr main ./seeder-workflow-template.yaml ./.workflows 
  '''
  event_type = argv[1]
  branch = argv[2]
  seeder_workflow_template_path = argv[3]
  workflows_dir = argv[4]
  seed(event_type, seeder_workflow_template_path, workflows_dir).dump()



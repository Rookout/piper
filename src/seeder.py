from sys import argv
from os import getenv
from yaml import safe_load, safe_dump

class Workflow():
  def __init__(self, wf_path):
    with open(wf_path) as f:
      self.wf = safe_load(f)

  def inject_templates(self, templates_path):
    try:
      with open(templates_path) as f:
        additional_templates = safe_load(f)
      self.wf["spec"]["templates"] += additional_templates
    except FileNotFoundError:
      print("didn't find any {} file".format(templates_path))

  def inject_dag(self, dag_path):
    with open(dag_path) as f:
      dag_flow = safe_load(f)
    entrypoint = self.wf["spec"]["entrypoint"]
    dag_flow = [{"name": entrypoint, "dag": {"tasks": dag_flow }}]
    self.wf["spec"]["templates"] += dag_flow

  def inject_name(self, name):
    if name:
      self.wf["metadata"]["generateName"] = name + "-"
    else:
      self.wf["metadata"]["generateName"] = "missing-name-"
  
  def inject_parameters(self, param_path):
    try:
      with open(param_path) as f:
        parameters = safe_load(f)
      self.wf["spec"]["arguments"]["parameters"] = parameters
    except FileNotFoundError:
      print("didn't find any {} file".format(param_path))
  
  def inject_labels(self):
    labels = safe_load(getenv("LABELS"))
    labels_list = ["branch","repo","commit","user"]
    self.wf["metadata"]["labels"] = {k:v for (k,v) in labels.items() if k in labels_list}

  def dump(self):
    with open("workflow.yaml", "w") as f:
      safe_dump(self.wf, f, width=10000)    


def seed(seeder_workflow_path, dot_workflows_dir, branch=None, event_type=None):
  workflow = Workflow(seeder_workflow_path)
  workflow.inject_templates(dot_workflows_dir + "/template.yaml")
  try:
    workflow.inject_dag(dot_workflows_dir + "/" + event_type + branch + ".yaml")
  except FileNotFoundError:
    workflow.inject_dag(dot_workflows_dir + "/main.yaml")
  workflow.inject_name(branch)
  workflow.inject_parameters(dot_workflows_dir + "/" + "parameters.yaml")
  workflow.inject_labels()
  
  return workflow

if __name__ == "__main__":
  '''
  python seeder.py .../seeder-workflow-template.yaml .../.workflows merge main
  '''
  
  seeder_workflow_template_path = argv[1]
  dot_workflows_dir = argv[2]
  event_type = argv[3]
  branch = argv[4]
  
  seed(seeder_workflow_template_path, dot_workflows_dir, branch, event_type).dump()



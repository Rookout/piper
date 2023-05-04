import unittest
import io
from seeder import *
from os import getcwd,environ
from unittest.mock import patch


class TestVersioning(unittest.TestCase):

    @patch.dict(
        environ,
        { 
            "LABELS": '''{"branch": "test-branch","commit": "xxxxxxxxxxxxxx","events.argoproj.io/action-timestamp": 1676882259356,"events.argoproj.io/sensor": "seeder","events.argoproj.io/trigger": "seeder","repo": "somerepo","user": "gosharo","workflows.argoproj.io/phase": "Running","workflows.argoproj.io/resubmitted-from-workflow": "seeder-n2w5w"}'''
        }
    )
    def test_e2e(self):
        self.maxDiff = None

        self.wf = seed(getcwd() + "/../examples/seeder-workflow-template.yaml", getcwd() + "/../examples/.workflows", branch="test")


        self.wf.dump()

        with open(getcwd() + "/../examples/workflow.yaml") as f:
            self.wf_target = safe_load(f)
        self.assertDictEqual(self.wf.wf, self.wf_target)
        

if __name__ == "__main__":
    unittest.main()
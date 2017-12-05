#
# Module: legolas.py
# Dave Wu @ Nov.2017
#
import json
import re
import time

from pymongo import MongoClient, ReturnDocument
from bson.objectid import ObjectId

"""
dependencies:

1) mongo servers are running
   * mongo db: legolas
   * mongo col: actions, cases, job states

2) entities in mongo:
   * cases: required field: name, path, params(a dict/hash/js object)
   * actions: required field: case_name, case_path, name, prev_action, params(a dict/hash/js object)

3) action results are saved in mongo:

4) input context (job):
   * case_path
   * case_name
   * action_name
   * case_run_id
   * prev_action name
"""

class Legolas:
    """the context object used in Legolas action snippet"""

    def __init__(self, cstr):
        self._ctx = json.loads(cstr)
        self._mongo = MongoClient()          # TODO: using specific host and port
        self._param_in_case = {}
        self._context_of_case = {}
        self._param_in_action = {}
        self._prev_results = {}
        self._results_dict = {}
        self._tracing_dict = {}

        # shortcuts. None by default
        self.run_id = self._ctx.get("run_id")
        self.action_id = self._ctx.get("action_id")
        self.prev_action_id = self._ctx.get("prev_action_id")
        print(">>> context <<<")
        print("run_id:         " + str(self.run_id))
        print("action_id:      " + str(self.action_id))
        print("prev_action_id: " + str(self.prev_action_id))

        # init params and prev result
        self._set_param_of_action()
        self._set_param_of_case()
        self._set_context_of_case()
        self._set_prev_results()

    #----------------------
    # helpers used by init
    #----------------------
    def _set_context_of_case(self):
        print(">>> loading context dict")
        col = self._mongo.legolas.runs
        rr = col.find_one({"_id": ObjectId(self.run_id)})
        if rr:
            self._context_of_case = rr["context"]

    def _set_param_of_action(self):
        print(">>> loading params of the action")
        col = self._mongo.legolas.actions
        t = col.find_one({"_id": ObjectId(self.action_id)})
        if t:
            self._param_in_action = t["params"]

    def _set_param_of_case(self):
        print(">>> loading params of the case")
        col = self._mongo.legolas.runs
        rr = col.find_one({"_id": ObjectId(self.run_id)})
        if rr and rr["case_id"]:
            col = self._mongo.legolas.cases
            tc = col.find_one({"_id": ObjectId(rr["case_id"])})
            if tc:
                self._param_in_case = tc["params"]

    def _set_prev_results(self):
        if self.prev_action_id:
            print(">>> loading previous results")
            col = self._mongo.legolas.jobstates
            t = col.find_one({"run_id": ObjectId(self.run_id), "action_id": ObjectId(self.prev_action_id)})
            if t:
                print('>>> got previous job state. its results: ' + str(t["results"]))
                self._prev_results = t["results"]
            else:
                print('>>> found no prev jobstate of id: ' + self.prev_action_id)

    #-----------------------------
    # helpers for running process
    #-----------------------------
    def _upload_results(self):
        """uploading results to jobstate and traced data to case run"""
        print(">>> uploading context data to run")
        if self._context_of_case:
            col = self._mongo.legolas.runs
            t = col.find_one_and_update(
                {"_id": ObjectId(self.run_id)},
                {'$set': {'context': self._context_of_case}},
                return_document=ReturnDocument.AFTER)
            if t:
                print(">>> upload context data succeeded")

        print(">>> uploading results to jobstate")
        if self._results_dict:
            col = self._mongo.legolas.jobstates
            t = col.find_one_and_update(
                {"run_id": ObjectId(self.run_id), "action_id": ObjectId(self.action_id)},
                {'$set': {'results': self._results_dict}},
                return_document=ReturnDocument.AFTER)
            if t:
                print(">>> upload result succeeded")

        print(">>> uploading traced data to case run")
        if self._tracing_dict:
            col = self._mongo.legolas.runs
            t = col.find_one_and_update(
                {"_id": ObjectId(self.run_id)},
                {'$set': {'traced_data': self._tracing_dict}},
                return_document=ReturnDocument.AFTER)
            if t:
                print(">>> upload traced data succeeded")

    def _set_start_time(self):
        col = self._mongo.legolas.jobstates
        t = col.find_one_and_update(
            {"run_id": ObjectId(self.run_id), "action_id": ObjectId(self.action_id)},
            {'$set': {'started_at': time.ctime()}},
            return_document=ReturnDocument.AFTER)
        if t:
            print('>>> set start time succeeded')

    def _set_end_time(self):
        col = self._mongo.legolas.jobstates
        t = col.find_one_and_update(
            {"run_id": ObjectId(self.run_id), "action_id": ObjectId(self.action_id)},
            {'$set': {'ended_at': time.ctime()}},
            return_document=ReturnDocument.AFTER)
        if t:
            print('>>> set end time succeeded')


    #-------------------------
    # methods for users to use
    #-------------------------
    def get_param(self, name):
        """search in action param first, then if in format of '$(xxx)', search case's params"""
        if self._param_in_action.has_key(name):
            if name:
                p = self._param_in_action[name]
                print("> get param: " + str(name) + ' : ' + str(p))
                matchObj = re.match(r'^\$\((.*)\)$', str(p))
                if matchObj:
                    pname = matchObj.group(1)
                    if self._param_in_case.has_key(pname):
                        return self._param_in_case[pname]
                else:
                    return p

    def get_prev_result(self, name):
        if name:
            name = str(name)
            name = name.replace(".", "_")
            return self._prev_results[name]

    def save_result(self, name, value):
        if name:
            name = str(name)
            name = name.replace(".", "_")
            self._results_dict[name] = value

    def save_context(self, name, value):
        if name:
            name = str(name)
            name = name.replace(".", "_")
            self._context_of_case[name] = value

    def get_context(self, name):
        if name:
            name = str(name)
            name = name.replace(".", "_")
            return self._context_of_case[name]

    def trace_data(self, name, value):
        if name:
            name = str(name)
            name = name.replace(".", "_")
            self._tracing_dict[name] = value

    def set_failed(self, msg):
        """set job as failed, along with error message"""
        print("> set job as failed")
        col = self._mongo.legolas.jobstates
        col.find_one_and_update(
            {"run_id": ObjectId(self.run_id), "action_id": ObjectId(self.action_id)},
            {'$set': {'state': "failed", 'error': msg}})

if __name__ == "__main__":
    print("don't run me. just import me.")

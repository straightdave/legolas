#
# Module: legolas.py
# Dave Wu @ Nov.2017
#
import json
import re

from pymongo import MongoClient
import redis

"""
dependencies:

1) redis and mongo servers are running
   * mongo db: legolas
   * mongo col: actions, cases

2) entities in mongo:
   * cases: required field: name, path, params(a dict/hash/js object)
   * actions: required field: case_name, case_path, name, prev_action, params(a dict/hash/js object)

3) action results are saved in redis:
   * key: <case_run_id>#<action_name>, value: number/text (redis types)

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
        print("loading context string: " + cstr)
        self._ctx = json.loads(cstr)
        self._mongo = MongoClient()          # TODO: using specific host and port
        self._redis = redis.StrictRedis()    # TODO: using specific host and port
        self._param_in_case = {}
        self._param_in_action = {}
        self._result_pre_action = {}

        # shortcuts. None by default
        self.case_path   = self._ctx.get("case_path")
        self.case_name   = self._ctx.get("case_name")
        self.action_name = self._ctx.get("action_name")
        self.case_run_id = self._ctx.get("case_run_id")
        self.prev_action = self._ctx.get("prev_action")

        # init params and prev result
        self._set_param_of_action()
        self._set_param_of_case()
        self._set_prev_results()

    def _set_param_of_action(self):
        col = self._mongo.legolas.actions
        t = col.find_one({"case_path": self.case_path, "case_name": self.case_name, "name": self.action_name})
        if t: self._param_in_action = t["params"]

    def _set_param_of_case(self):
        col = self._mongo.legolas.cases
        t = col.find_one({"path": self.case_path, "name": self.case_name})
        if t: self._param_in_case = t["params"]

    def _set_prev_results(self):
        if self.prev_action:
            prev_act_res_key = "{0}#{1}".format(self.case_run_id, self.prev_action)
            t = self._redis.hgetall(prev_act_res_key)
            i = 0
            while i < len(t):
                if i % 2 == 1:
                    self._result_pre_action[t[i - 1]] = t[i]
                i += 1

    #-------------------------
    # methods for users to use
    #-------------------------
    def get_param(self, name):
        """search in action param first, then if in format of '$(xxx)', search case's params"""
        if self._param_in_action.has_key(name):
            p = self._param_in_action[name]
            print("get param:name=" + str(p))

            matchObj = re.match(r'^\$\((.*)\)$', str(p))
            if matchObj:
                pname = matchObj.group(1)
                if self._param_in_case.has_key(pname):
                    return self._param_in_case[pname]
            else:
                return p

    def save_result(self, name, value):
        reskey = "{0}#{1}".format(self.case_run_id, self.action_name)
        self._redis.hset(reskey, name, value)

if __name__ == "__main__":
    print("don't run me. just import me.")

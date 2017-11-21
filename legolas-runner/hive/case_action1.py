# sample action snippet (injected with some other statements)
import sys
import time
from legolas import Legolas

#==========================
# user defined begins
#==========================
def action_main(ctx):
    w = "context.action_name = {0}".format(ctx.action_name)
    ctx.save_result('hello', w)

#==========================
# user defined ends
#==========================

if __name__ == "__main__":
    ctx = Legolas(sys.argv[1])
    ctx.save_result('started_at', time.ctime())
    action_main(ctx)
    ctx.save_result('ended_at', time.ctime())

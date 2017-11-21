#
# --casename--
#
import sys
import time
from legolas import Legolas

#==========================
# user defined begins
#==========================

--snippet--

#==========================
# user defined ends
#==========================

if __name__ == "__main__":
    ctx = Legolas(sys.argv[1])
    ctx.save_result('started_at', time.ctime())
    action_main(ctx)
    ctx.save_result('ended_at', time.ctime())

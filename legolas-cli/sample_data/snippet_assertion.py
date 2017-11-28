#
# Template: Assertion
#
def action_main(ctx):
    expected = ctx.get_param("expected")
    actual = ctx.get_param("actual")
    fail_msg = ctx.get_param("fail_msg")

    if expected != actual:
        # set job state as failed
        ctx.set_failed("{0}: expected '{1}', actual '{2}.".format(
            fail_msg, expected, actual))

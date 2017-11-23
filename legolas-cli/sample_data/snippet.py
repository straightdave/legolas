# users can add their own functions here,
# but only 'action_main(ctx)' would be invoked directly

def action_main(ctx):
    print("hello legolas")

    # since in the template, it defines such params,
    # so any applying actions would have such params
    # in this template, it defines:
    # name : string
    # age  : number
    my_name = ctx.get_param("name")  # a string
    my_age = ctx.get_param("age")    # a number

    print("hello, {0}. you are {1}.".format(my_name, my_age))

    # save something as the result of this action
    # which can be used in another actions in the same case run
    ctx.save_result("act-1.greetingword", "hello legolas!")

    print("ended show")

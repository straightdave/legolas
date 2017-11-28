#
# Template snippet: RESTful API call
#

def action_main(ctx):
    host = ctx.get_param("host")
    api_url = ctx.get_param("api_url")
    method = ctx.get_param("method")
    body = ctx.get_param("body")

    from python_http_client import Client




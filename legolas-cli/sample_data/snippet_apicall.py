def action_main(ctx):
    import requests
    host = ctx.get_param('host')
    url  = ctx.get_param('url')

    _api = "{0}{1}".format(host, url)
    print("sending GET request to " + _api)

    r = requests.get(_api)
    ctx.save_result('response', r.text)

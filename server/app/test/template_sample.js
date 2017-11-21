//
// action template sample
//

var template = {
    "name": "add_driver",
    "path": "$template/path/",
    "desc": "add a driver's location",

    "params": [
        {
            "name": "Id",
            "desc": "driver's Id",
            "type": "text",
            "required": true,
            "default": ""
        },
        {
            "name": "Name",
            "desc": "driver's legal name",
            "type": "text",
            "required": true,
            "default": ""
        },
        {
            "name": "Location",
            "desc": "driver's location",
            "type": "geo",
            "required": true,
            "default": ""
        }
    ],

    "snippet": `
    print('you are adding drivers:' + Legolas.GetParam('Name'))
    res = Legolas.CallApi('add_driver_location', Legolas.GetParam('Id'), Legolas.GetParam('Location'))
    Legolas.SetResult(res)
    print('done')
    `
}

var action = {
    "name": "wahaha, add a driver",
    "desc": "whahah hahaha add whahah",
    "case": "$/case/path/case-1",

    "template": {
        "name": "Add a driver",
        "path": "$template/path/"
    },

    "params": [
        {
            "name": "Id",
            "value": "d001001001"
        },
        {
            "name": "Name",
            "value": "Ali Mehamedov"
        },
        {
            "name": "Location",
            "value": "1001.01,202.45"
        }
    ]
}

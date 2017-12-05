<template>
<div id="app-detail">
    <div id="detail-header">
        <p><input id="path" type="text" v-model="localInfo.path"></input></p>
        <p><input id="name" type="text" v-model="localInfo.name"></input></p>
        <p><input id="desc" type="text" v-model="localInfo.desc"></input></p>
        <div id="buttons">
            <button v-on:click="save()">Save</button>
        </div>
    </div>
    <div id="nav-in-template">
        <ul>
            <li><a v-on:click.stop.prevent="switchToTab(0)" :class="{active: activeTab == 0}">Snippet</a></li>
            <li><a v-on:click.stop.prevent="switchToTab(1)" :class="{active: activeTab == 1}">Params</a></li>
        </ul>
    </div>
    <div id="detail-panel">
        <div :class="{hidden: activeTab != 0}">
            <pre id="brace-editor">{{ localInfo.snippet }}</pre>
        </div>

        <div :class="{hidden: activeTab != 1}">
            <div id="param-list">
                <div class="param-list-item" v-for="(p, index) in paramList" :key="index">
                    <a v-on:click.stop.prevent="removeParam(p.name)">
                        <i class="fa fa-minus-circle"></i>
                    </a>
                    <input type="text" v-model="p.name" size="20" />
                    <select v-model="p.type">
                        <option value="text"    :selected="p.type == 'text'">text</option>
                        <option value="number"  :selected="p.type == 'number'">number</option>
                        <option value="boolean" :selected="p.type == 'boolean'">boolean</option>
                    </select>
                    <input :id="index + '_required'" type="checkbox" name="isRequired" v-model="p.required" />
                    <label :for="index + '_required'">Required</label>
                    <input type="text" v-model="p.default" size="50" />
                </div>
                <div id="new-param-box">
                    <a v-on:click.stop.prevent="newParam()">
                        <i class="fa fa-plus"></i> New
                    </a>
                </div>
            </div>
        </div>
    </div>
</div>
</template>

<script>
import Vue from 'vue'
import VueResource from 'vue-resource'
Vue.use(VueResource)

var ace = require('brace')
require('brace/mode/python')
require('brace/theme/solarized_light')

var AppTemplateDetail = Vue.extend({
    props: {
        templateObject: {
            type: Object,
            required: true
        }
    },
    watch: {
        templateObject(newOne) {
            console.log('user click other template: ' + JSON.stringify(newOne))
            this.localInfo = JSON.parse(JSON.stringify(newOne))
            this.isNew = newOne.isNew
            this.activeTab = 0
            this.initEditor()
            this.initParameters()
        }
    },
    mounted() {
        this.initEditor()
        this.initParameters()
    },
    data() {
        return {
            activeTab: 0,
            localInfo: JSON.parse(JSON.stringify(this.templateObject)),
            paramList: [],
            isNew: this.templateObject.isNew,

            editor: null
        }
    },
    methods: {
        switchToTab(tab) {
            this.activeTab = tab
        },
        initParameters() {
            if (!this.localInfo.hasOwnProperty('params')) {
                this.localInfo['params'] = {}
            }
            var pdict = this.localInfo.params
            this.paramList = Object.keys(pdict).map(key => ({
                name:     key,
                required: pdict[key].required,
                type:     pdict[key].type,
                default:  pdict[key].default
            }))
        },
        initEditor() {
            var content = ''
            if (!this.isNew) {
                content = this.localInfo.snippet
            }

            var editor = ace.edit('brace-editor')
            editor.$blockScrolling = Infinity
            editor.getSession().setMode('ace/mode/python')
            editor.setTheme('ace/theme/solarized_light')
            editor.setValue(content)
            this.editor = editor
        },
        newParam() {
            console.log('adding a new param')
            this.paramList.push({
                name: '',
                required: false,
                type: 'text'
            })
        },
        removeParam(key) {
            console.log('removing a param: ' + key)
            var r = confirm("remove this param?")
            if (r != true) {
                return
            }
            var index = this.paramList.findIndex(i => i.name === key)
            this.paramList.splice(index, 1)
        },
        save() {
            this.localInfo.snippet = this.editor.getValue()

            this.localInfo.params = {}
            for (var item of this.paramList) {
                console.log('item: ' + JSON.stringify(item))
                var _k = item.name.trim()
                if (_k !== 'new' && _k !== '') {
                    this.localInfo.params[_k] = {
                        'type': item.type,
                        'required': item.required,
                        'default': item.default
                    }
                }
            }

            console.log('updating template...')
            if (!this.isNew) {
                var template_id = this.templateObject._id
                this.$http.put(`/template/${encodeURI(template_id)}`, this.localInfo).then(
                    resp => {
                        console.log('save template succeeded: ' + JSON.stringify(resp))
                        this.$emit('save-succeeded')
                    },
                    resp => {
                        console.log('http put failed: ' + JSON.stringify(resp))
                    })
            }
            else {
                console.log('creating new template...')
                this.$http.post(`/templates`, this.localInfo).then(
                    resp => {
                        console.log('save template succeeded: ' + JSON.stringify(resp))
                        this.$emit('save-succeeded')
                    },
                    resp => {
                        console.log('http put failed: ' + JSON.stringify(resp))
                    })
            }
        }
    }
})
export default AppTemplateDetail
</script>

<style scoped>
div.hidden {
    display: none;
}

div#app-detail {
    float: left;
    padding: 10px;
    min-width: 1100px;
}

div#detail-header {
    height: 80px;
    overflow: hidden;
}

div#detail-header input[type="text"] {
    border: 0;
    width: 80%;
    background-color: #fff;
}

input#path, input#desc {
    font-size: 15px;
    color: gray;
}

input#name {
    font-size: 30px;
    font-weight: 300;
}

div#buttons {
    position: absolute;
    right: 10px;
    top: 10px;
}

div#detail-header div#buttons button {
    height: 30px;
    width: 60px;
    font-size: 15px;
    border: none;
    background-color: #00B140;
    color: #fff;
    cursor: pointer;
}

div#nav-in-template ul {
    display: inline-flex;
    list-style: none;
    margin-top: 15px;
    margin-bottom: 15px;
}

div#nav-in-template ul li {
    padding-right: 20px;
}

div#nav-in-template ul li a {
    text-decoration: none;
    color: gray;
    cursor: pointer;
    display: block;
}

div#nav-in-template li a:hover,
div#nav-in-template li a.active {
    border-bottom: solid 2px #00B140;
    padding-bottom: 8px;
}

pre#brace-editor {
    font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', 'source-code-pro', monospace;
    font-size: 18px;
    font-weight: normal;
    min-height: 600px;
}

div#new-param-box {
    margin-top: 10px;
}
div#new-param-box a {
    font-size: 18px;
    text-decoration: none;
    color: gray;
    cursor: pointer;
    display: block;
}

div.param-list-item input[type='text'] {
    font-weight: 200;
    font-size: 18px;
    border: 0;
    background-color: #fff;
}
div.param-list-item select {
    font-size: 18px;
    font-weight: 200;
}

div.param-list-item {
    margin-bottom: 10px;
}

div.param-list-item a {
    font-size: 18px;
    text-decoration: none;
    color: #00B140;
    cursor: pointer;
}

</style>

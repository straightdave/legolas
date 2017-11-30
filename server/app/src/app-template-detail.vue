<template>
<div id="app-detail">
    <div id="detail-header">
        <p><input id="path" type="text" v-model="localInfo.path"></input></p>
        <p><input id="name" type="text" v-model="localInfo.name"></input></p>
        <p><input id="desc" type="text" v-model="localInfo.desc"></input></p>
        <div id="buttons">
            <button v-on:click="save">Save</button>
        </div>
    </div>
    <div id="nav-in-template">
        <ul>
            <li><a v-on:click.stop.prevent="activeTab = 0" v-bind:class="{active: activeTab == 0}">Snippet</a></li>
            <li><a v-on:click.stop.prevent="activeTab = 1" v-bind:class="{active: activeTab == 1}">Params</a></li>
        </ul>
    </div>
    <div id="detail-panel">
        <div v-if="activeTab == 0">
            <pre id="brace-editor">{{ localInfo.snippet }}</pre>
        </div>

        <div v-else-if="activeTab == 1">
            <div id="param-list">
                <div v-for="(value, key) in localInfo.params">
                    {{ key }} : {{ value }}
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
        templateObject: function (newOne) {
            console.log('user click other template: ' + JSON.stringify(newOne))
            this.localInfo = JSON.parse(JSON.stringify(newOne))
            this.isNew = newOne.isNew
            this.initEditor()
        }
    },
    data() {
        return {
            activeTab: 0,
            localInfo: JSON.parse(JSON.stringify(this.templateObject)),
            isNew: this.templateObject.isNew,
            editor: null
        }
    },
    mounted() {
        this.initEditor()
    },
    updated() {
        this.initEditor()
        this.initParamList()
    },
    methods: {
        initParamList() {
            if (this.activeTab != 1) {
                return
            }

            var params = {}
            if (!this.isNew) {
                params = this.localInfo.params
            }
        },
        initEditor() {
            if (this.activeTab != 0) {
                return
            }

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
        save() {
            this.localInfo.snippet = this.editor.getValue()
            console.log('updating template...')
            if (!this.isNew) {
                var template_id = this.templateObject._id
                this.$http.put(`/template/${encodeURI(template_id)}`, this.localInfo).then(
                    resp => {
                        console.log('save template succeeded. new one: ' + JSON.stringify(resp))
                    },
                    resp => {
                        console.log('http put failed: ' + JSON.stringify(resp))
                    }
                )
            }
            else {
                console.log('creating new template...')
                this.$http.post(`/templates`, this.localInfo).then(
                    resp => {
                        console.log('save template succeeded. new one: ' + JSON.stringify(resp))
                    },
                    resp => {
                        console.log('http put failed: ' + JSON.stringify(resp))
                    }
                )
            }
        }
    }
})
export default AppTemplateDetail
</script>

<style scoped>
div#app-detail {
    float: left;
    padding: 10px;
    min-width: 1100px;
    overflow: scroll;
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
    height: 500px;
}

</style>

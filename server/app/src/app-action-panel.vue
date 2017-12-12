<!-- app action panel -->
<template>
<div id="action-panel-main">
    <div id="template-store" :class="{hidden: notToShowStore}">
        <div v-for="t in templates" class="template-box" :key="t._id" @click="selectTemplate(t._id)">
            <div class="name">{{ t.name }}</div>
            <div class="desc">{{ t.desc }}</div>
            <div class="createdat"><i class="fa fa-clock-o"></i> {{ (new Date(t.created_at)).toLocaleString() }}</div>
        </div>
    </div>

    <div id="action-panel" :class="{hidden: !notToShowStore}">
        <div class="action-inner-div">
            <div class="capital">
                <span>Basic info</span>
            </div>

            <div id="action-name-box">
                <input type="text" v-model="localActionObject.name" placeholder="Action name">
            </div>
            <div id="action-desc-box">
                <input type="text" v-model="localActionObject.desc" placeholder="Desc goes here...">
            </div>
            <div id="action-seqno-box">
                <label for="seq-no-input">Seq No.</label>
                <input id="seq-no-input" type="text" v-model="localActionObject.seq_no">
            </div>
        </div>

        <div id="tpl-info" class="action-inner-div">
            <div class="capital">
                <span>Template Info</span>
            </div>

            <div id="tpl-title">
                <!-- TODO: use route link to that template/need more specific routers -->
                <span>{{ templateInfo.name }}</span>
            </div>
            <div id="tpl-desc">
                <span>{{ templateInfo.desc }}</span>
            </div>
        </div>

        <div id="param-info" class="action-inner-div">
            <div class="capital">
                <span>Parameters</span>
            </div>

            <div id="param-list">
                <div class="param-list-item" v-for="(p, index) in paramList" :key="index">
                    <a @click.stop.prevent="removeParam(p.name)">
                        <i class="fa fa-minus-circle"></i>
                    </a>
                    <input type="text" v-model="p.name" size="20" />
                    <input type="text" v-model="p.value" size="30" />
                </div>
                <div class="new-param-box">
                    <a @click.stop.prevent="newParam()">
                        <i class="fa fa-plus"></i> New
                    </a>
                </div>
            </div>
        </div>

        <div id="mock-data" class="action-inner-div">
            <div class="capital">
                <span>Mockingbird</span>
                <span><input type="checkbox" v-model="is_mocking"></span>
            </div>

            <div v-if="is_mocking">
                <p>Mock data will be automatically saved to <strong>action's result</strong> and <strong>Run's context</strong> dictionary during runtime</p>
                <div id="mock-data-list">
                    <div class="param-list-item" v-for="(p, index) in mockDataList" :key="index">
                        <a @click.stop.prevent="removeMockData(p.name)">
                            <i class="fa fa-minus-circle"></i>
                        </a>
                        <input type="text" v-model="p.name" size="20" />
                        <input type="text" v-model="p.value" size="30" />
                    </div>
                    <div class="new-param-box">
                        <a @click.stop.prevent="newMockData()">
                            <i class="fa fa-plus"></i> New
                        </a>
                    </div>
                </div>
            </div>
        </div>

        <div class="action-inner-div">
            <button @click="saveAction">Save</button>
            <button @click="deleteAction">Delete</button>
        </div>
    </div>
</div>
</template>

<script>
import Vue from 'vue'
import VueResource from 'vue-resource'
Vue.use(VueResource)

var AppActionPanel = Vue.extend({
    props: {
        actionObject: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            notToShowStore: this.actionObject.hasOwnProperty('template_id') && this.actionObject.template_id != '',
            localActionObject: JSON.parse(JSON.stringify(this.actionObject)),
            isNew: this.actionObject.isNew,
            paramList: [],

            templates: [],
            templateInfo: {name: "", desc: ""}, // init value is important for async programming

            is_mocking: !!this.actionObject.is_mocking,
            mockDataList: []
        }
    },
    created() {
        this.initMockData()
        this.initTemplates()
        this.initTemplateInfo()
    },
    watch: {
        // caused by use selecting another action in the list
        actionObject(newActionObject) {
            this.localActionObject = JSON.parse(JSON.stringify(newActionObject))
            this.isNew = newActionObject.isNew
            console.log('changed to action: ' + JSON.stringify(this.localActionObject))

            this.initTemplateInfo()
            this.notToShowStore = newActionObject.hasOwnProperty('template_id') && newActionObject.template_id != ''
        }
    },
    methods: {
        initMockData() {
            console.log('init mock data...')
            // flatten mock data into arrays
            if (!!this.localActionObject.mock_data) {
                var pdict = this.localActionObject.mock_data
                this.mockDataList = Object.keys(pdict).map(key => ({
                    name:  key,
                    value: pdict[key]
                }))
            }
        },
        initTemplates() {
            console.log('init templates ...')
            this.$http.get('/templates').then(
                resp => {
                    this.templates = resp.body
                },
                resp => {
                    console.log('http failed: ' + JSON.stringify(resp.body))
                })
        },
        initTemplateInfo() {
            if (!this.localActionObject.hasOwnProperty('template_id')) {
                console.log('has no template id yet')
                return
            }

            console.log('getting template info of id: ' + this.localActionObject.template_id)
            var url = `/template/${encodeURI(this.localActionObject.template_id)}`
            this.$http.get(url).then(
                resp => {
                    this.templateInfo = resp.body
                    console.log('got template info: ' + JSON.stringify(this.templateInfo))
                },
                resp => {
                    console.log('http get failed: ' + JSON.stringify(resp))
                }).then(() => {
                    this.initParameters()
                })
        },
        initParameters() {
            if (!!this.isNew) {
                console.log('== init parameters for new action ...')
                if (this.localActionObject.hasOwnProperty('template_id')) {
                    console.log('looking for default param of template id: ' + this.localActionObject.template_id)

                    if (this.templateInfo.hasOwnProperty('params')) {
                        this.localActionObject.params = {}
                        var pTpl = this.templateInfo.params

                        for (var k of Object.keys(pTpl)) {
                            this.localActionObject.params[k] = pTpl[k].default // 'default' could be undefined
                        }
                        console.log('got params of default: ' + JSON.stringify(this.localActionObject.params))
                    }
                }
                console.log('== end init params for new action')
            }

            if (!this.localActionObject.hasOwnProperty('params')) {
                this.localActionObject['params'] = {}
            }

            // flatten params into arrays
            var pdict = this.localActionObject.params
            this.paramList = Object.keys(pdict).map(key => ({
                name:  key,
                value: pdict[key]
            }))
            console.log('got paramlist: ' + JSON.stringify(this.paramList))
        },
        newParam() {
            console.log('adding a param')
            this.paramList.push({
                name:  '',
                value: ''
            })
        },
        removeParam(key) {
            console.log('removing a param: ' + key)
            var r = confirm("remove this param?")
            if (!r) {
                return
            }
            var index = this.paramList.findIndex(i => i.name === key)
            this.paramList.splice(index, 1)
        },
        newMockData() {
            console.log('adding a param')
            this.mockDataList.push({
                name:  '',
                value: ''
            })
        },
        removeMockData(key) {
            console.log('removing a mock data: ' + key)
            var r = confirm("remove this data?")
            if (!r) {
                return
            }
            var index = this.mockDataList.findIndex(i => i.name === key)
            this.mockDataList.splice(index, 1)
        },
        saveAction() {
            // re-check the params
            this.localActionObject.params = {}
            for (var item of this.paramList) {
                var _k = item.name.trim()
                if (_k !== 'new' && _k !== '') {
                    this.localActionObject.params[_k] = item.value
                }
            }

            // re-check mock data
            this.localActionObject.mock_data = {}
            for (var item of this.mockDataList) {
                var _k = item.name.trim()
                if (_k !== 'new' && _k !== '') {
                    this.localActionObject.mock_data[_k] = item.value
                }
            }
            this.localActionObject.is_mocking = this.is_mocking

            if (!this.localActionObject.hasOwnProperty('case_id')) {
                console.log('saving action to a new case')
                alert('Please save the case first!')
                return
            }

            if (this.localActionObject.hasOwnProperty('seq_no')) {
                var seq = this.localActionObject.seq_no
                this.localActionObject.seq_no = parseInt(seq)
            }
            else {
                this.localActionObject.seq_no = 0
            }


            if (!!this.isNew) {
                console.log('save new action: ' + JSON.stringify(this.localActionObject))

                this.$http.post('/actions', this.localActionObject).then(
                    resp => {
                        var hasError = resp.body.hasOwnProperty("error")
                        if (hasError) {
                            alert(JSON.stringify(resp.body))
                            return
                        }
                        this.isNew = false
                        this.$emit('action-list-refresh-needed', true)
                    },
                    resp => {
                        console.log('http failed: ' + JSON.stringify(resp.body))
                    })
            }
            else {
                console.log('save existing action')
                this.$http.put(`/action/${encodeURI(this.localActionObject._id)}`, this.localActionObject).then(
                    resp => {
                        var hasError = resp.body.hasOwnProperty("error")
                        if (hasError) {
                            alert(resp.body.error)
                            return
                        }
                        console.log('http succeeded: ' + JSON.stringify(resp.body))
                        this.$emit('action-list-refresh-needed', true)
                    },
                    resp => {
                        console.log('http failed: ' + JSON.stringify(resp.body))
                    })
            }
        },
        deleteAction() {
            if (!!this.isNew) {
                console.log('delete a new (unsafed) action')
                this.$emit('action-list-refresh-needed', true)
                return
            }

            var r = confirm("delete this action (not an unsafed one)?")
            if (!r) {
                return
            }

            var url = `/action/${encodeURI(this.localActionObject._id)}`
            this.$http.delete(url).then(
                resp => {
                    console.log('succussfully deleted the action: ' + JSON.stringify(resp.body))
                    this.$emit('action-list-refresh-needed', true)
                },
                resp => {
                    console.log('failed to delete action: ' + JSON.stringify(resp.body))
                })
        },
        selectTemplate(tid) {
            console.log('select template: ' + tid)
            this.localActionObject.template_id = tid
            this.initTemplateInfo()
            this.notToShowStore = true
        }
    }
})
export default AppActionPanel
</script>

<style scoped>
div.hidden {
    display: none;
}

div#action-panel-main {
    width: 800px;
    min-width: 600px;
    float: left;
    margin-left: 10px;
}

div#template-store {
    overflow-y: scroll;
}

div.template-box {
    float: left;
    height: 200px;
    width: 200px;
    border: solid 4px #ececec;

    margin: 0 5px 5px 0;
    position: relative;
    cursor: pointer;
}

div.template-box > div.name {
    text-align: center;
    background-color: #ececec;
    font-size: 20px;
    margin-bottom: 5px;
    padding: 3px 0;
}
div.template-box > div.desc {
    padding: 5px;
}
div.template-box > div.createdat {
    position: absolute;
    left: 10px;
    bottom: 5px;
}

div#action-panel > div {
    margin-bottom: 10px;
}

div#action-name-box input[type="text"]
{
    font-size: 20px;
    padding:3px;
    width: 200px;
    border: none;
}
div#action-desc-box input[type="text"]
{
    font-size: 15px;
    padding:3px;
    width: 200px;
    border: none;
}

div#action-panel button {
    height: 35px;
    font-size: 15px;
    width: 80px;
    border: none;
    background-color: #ececec;
    cursor: pointer;
}

div.action-inner-div {
    margin-bottom: 10px;
}

div.capital {
    background-color: #ececec;
    padding: 5px;
    margin-bottom: 5px;
}

div#tpl-title {
    font-size: 18px;
    font-weight: 400;
    margin-bottom: 3px;
}

div.new-param-box {
    margin-top: 10px;
}
div.new-param-box a {
    font-size: 18px;
    text-decoration: none;
    color: gray;
    cursor: pointer;
    display: block;
}

div.param-list-item {
    margin-bottom: 5px;
}
div.param-list-item input[type='text'] {
    font-weight: 200;
    font-size: 18px;
    background-color: #fff;
}
div.param-list-item a {
    font-size: 18px;
    text-decoration: none;
    color: #00B140;
    cursor: pointer;
}
</style>

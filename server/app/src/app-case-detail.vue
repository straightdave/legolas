<!-- app detail -->
<template>
<div id="app-detail">
    <div id="detail-header">
        <p><input id="path" type="text" v-model="localCaseInfo.path"></input></p>
        <p><input id="name" type="text" v-model="localCaseInfo.name"></input></p>
        <p><input id="desc" type="text" v-model="localCaseInfo.desc"></input></p>
        <div id="buttons">
            <button @click="saveCase">Save</button>
            <button @click="runCase">Run</button>
        </div>
    </div>
    <div id="nav-in-case">
        <ul>
            <li><a @click.stop.prevent="cliTab(0)" :class="{active: activeTab == 0}">Actions</a></li>
            <li><a @click.stop.prevent="cliTab(1)" :class="{active: activeTab == 1}">Variables</a></li>
            <li><a @click.stop.prevent="cliTab(2)" :class="{active: activeTab == 2}">Runs</a></li>
            <li><a @click.stop.prevent="cliTab(3)" :class="{active: activeTab == 3}">Tracing</a></li>
        </ul>
    </div>
    <div id="detail-panel">
        <div :class="{hidden: activeTab != 0}">
            <div id="innerlist">
                <AppAction
                    v-if="hasActions"
                    v-for="a in actions"
                    :key="a.name"
                    :action="a"
                    @actionClicked="setCurrentAction(a)" />
                <div id="newaction" @click="addNewAction">
                    <i class="fa fa-plus-square"></i> new
                </div>
            </div>
            <AppActionPanel
                v-if="hasCurrentAction"
                :action-object="currentAction"
                @action-list-refresh-needed="refreshActionList" />
        </div>

        <div :class="{hidden: activeTab != 1}">
            <div id="param-list">
                <div class="param-list-item" v-for="(p, index) in paramList" :key="index">
                    <a @click.stop.prevent="removeParam(p.name)">
                        <i class="fa fa-minus-circle"></i>
                    </a>
                    <input type="text" v-model="p.name" size="20" /> :
                    <input type="text" v-model="p.value" size="50" />
                </div>
                <div id="new-param-box">
                    <a @click.stop.prevent="newParam()">
                        <i class="fa fa-plus"></i> New
                    </a>
                </div>
            </div>
        </div>

        <div :class="{hidden: activeTab != 2}">
            <div id="innerlist">
                <AppRun
                    v-if="hasRuns"
                    v-for="run in runs"
                    :key="run._id"
                    :run="run"
                    @run-clicked="setCurrentRun(run)" />
            </div>
            <AppRunPanel v-if="hasCurrentRun" :run-object="currentRun" />
        </div>

        <div :class="{hidden: activeTab != 3}">
            Data tracing
        </div>
    </div>
</div>
</template>

<script>
import Vue from 'vue'
import VueResource from 'vue-resource'
Vue.use(VueResource)

import AppAction from './app-action.vue'
import AppActionPanel from './app-action-panel.vue'
import AppRun from './app-run.vue'
import AppRunPanel from './app-run-panel.vue'

var AppCaseDetail = Vue.extend({
    components: {AppAction, AppActionPanel, AppRun, AppRunPanel},
    props: {
        caseInfo: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            localCaseInfo: JSON.parse(JSON.stringify(this.caseInfo)),
            isNew: this.caseInfo.isNew,

            activeTab: 0,

            actions: [],
            runs: [],
            paramList: [],

            currentAction: null,
            currentRun: null
        }
    },
    watch: {
        caseInfo(newCaseInfo) {
            console.log('user click the case: ' + newCaseInfo.name)
            this.localCaseInfo = JSON.parse(JSON.stringify(newCaseInfo))
            this.isNew = newCaseInfo.isNew
            this.refreshActionList(true)
            this.initRuns()
            this.initParameters()
        }
    },
    computed: {
        hasActions() {
            return this.actions && this.actions.length > 0
        },
        hasRuns() {
            return this.runs && this.runs.length > 0
        },
        hasCurrentAction() {
            return this.currentAction !== null
        },
        hasCurrentRun() {
            return this.currentRun !== null
        }
    },
    mounted() {
        this.refreshActionList(true)
        this.initRuns()
        this.initParameters()
    },
    methods: {
        initParameters() {
            if (!this.localCaseInfo.hasOwnProperty('params')) {
                this.localCaseInfo['params'] = {}
            }
            var pdict = this.localCaseInfo.params
            this.paramList = Object.keys(pdict).map(key => ({
                name:  key,
                value: pdict[key]
            }))
        },
        initRuns() {
            this.currentRun = null
            console.log(`init runs for case: ${this.localCaseInfo.name}, isNew: ${!!this.isNew}`)

            if (!!this.isNew) {
                console.log('new case, no need to fetch runs')
                this.runs = []
                return
            }

            var url = `/case/${encodeURI(this.caseInfo._id)}/runs`
            this.$http.get(url).then(
                resp => {
                    console.log('get runs succeeded')
                    this.runs = resp.body
                },
                resp => {
                    console.log('http put failed: ' + resp.body)
                })
        },
        cliTab(item) {
            this.activeTab = item
        },
        setCurrentAction(act) {
            console.log('set current action to: ' + JSON.stringify(act))
            this.currentAction = act
        },
        setCurrentRun(run) {
            console.log('set current run to: ' + JSON.stringify(run))
            this.currentRun = run
        },
        addNewAction() {
            // WARNING: it could be new case (no _id yet)
            var newAction = {
                case_id: this.localCaseInfo._id,
                name: "action-new",
                desc: "this is a new action.",
                isNew: true
            }
            console.log('add action at list: ' + JSON.stringify(newAction))
            if (!this.actions) {
                this.actions = []
            }
            this.actions.push(newAction)
            this.currentAction = newAction
        },
        refreshActionList(toCloseActionPanel) {
            if (!!this.isNew) {
                console.log('new case, no fetching actions')
                this.actions = []
                this.currentAction = null
                return
            }

            console.log('refreshing action list')
            var url = `/case/${encodeURI(this.caseInfo._id)}/actions`

            this.$http.get(url).then(
                resp => {
                    console.log('get actions succeeded')
                    this.actions = resp.body
                    if (toCloseActionPanel) {
                        this.currentAction = null
                    }
                },
                resp => {
                    console.log('failed to get actions')
                })
        },
        newParam() {
            console.log('adding a new param')
            this.paramList.push({
                name: '',
                value: ''
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
        saveCase() {
            // get new param list
            this.localCaseInfo.params = {}
            for (var item of this.paramList) {
                var _k = item.name.trim()
                if (_k !== 'new' && _k !== '') {
                    this.localCaseInfo.params[_k] = item.value
                }
            }

            if (!!this.isNew) {
                console.log('saving new case...')
                this.$http.post('/cases', this.localCaseInfo).then(
                    resp => {
                        console.log('new case saved: ' + JSON.stringify(resp.body))
                        this.$emit('refresh-sidebar-list')
                    },
                    resp => {
                        console.log('http failed: ' + JSON.stringify(resp.body))
                    })
            }
            else {
                console.log('update existing one')
                var oldCaseUrl = `/case/${encodeURI(this.caseInfo._id)}`
                this.$http.put(oldCaseUrl, this.localCaseInfo).then(
                    resp => {
                        console.log('case updated: ' + JSON.stringify(resp.body))
                        this.$emit('refresh-sidebar-list')
                    },
                    resp => {
                        console.log('http failed: ' + JSON.stringify(resp.body))
                    })
            }
        },
        runCase() {
            if (!!this.isNew) {
                alert('cannot run an unsaved case')
                console.log('cannot run an unsaved case')
                return
            }

            console.log('run case: ' + this.caseInfo.name)
            var url = `/case/${encodeURI(this.caseInfo._id)}/runs`

            this.$http.post(url).then(
                resp => {
                    console.log('http success: ' + JSON.stringify(resp.body))
                    this.initRuns()
                },
                resp => {
                    console.log('http failed: ' + JSON.stringify(resp.body))
                })
        }
    }
})
export default AppCaseDetail
</script>

<style scoped>
div.hidden {
    display: none;
}
div#app-detail {
    float: left;
    padding: 10px;
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

div#detail-header div#buttons button{
    height: 30px;
    width: 60px;
    font-size: 15px;
    border: none;
    background-color: #00B140;
    color: #fff;
    cursor: pointer;
}

div#nav-in-case ul {
    display: inline-flex;
    list-style: none;
    margin-top: 15px;
    margin-bottom: 15px;
}

div#nav-in-case ul li {
    padding-right: 20px;
}

div#nav-in-case ul li a {
    text-decoration: none;
    color: gray;
    cursor: pointer;
    display: block;
}

div#nav-in-case li a:hover,
div#nav-in-case li a.active {
    border-bottom: solid 2px #00B140;
    padding-bottom: 8px;
}

div#innerlist {
    width: 300px;
    float: left;
}

div#newaction {
    width: 296px;
    text-align: center;
    font-size: 30px;
    cursor: pointer;
    color: #ececec;
    border: solid 2px #ececec;
}

/* for params */
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
    background-color: #fff;
    border: solid 1px gray;
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

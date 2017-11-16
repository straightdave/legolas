<!-- app detail -->
<template>
<div id="app-detail">
    <div id="detail-header">
        <p><input id="path" type="text" v-model="localCaseInfo.path"></input></p>
        <p><input id="name" type="text" v-model="localCaseInfo.name"></input></p>
        <p><input id="desc" type="text" v-model="localCaseInfo.desc"></input></p>
        <div id="buttons">
            <button v-on:click="saveCase">Save</button>
            <button v-on:click="runCase">Run</button>
        </div>
    </div>
    <div id="nav-in-case">
        <ul>
            <li><a v-on:click.stop.prevent="cliTab(0)" v-bind:class="{active: activeTab == 0}">Actions</a></li>
            <li><a v-on:click.stop.prevent="cliTab(1)" v-bind:class="{active: activeTab == 1}">Variables</a></li>
            <li><a v-on:click.stop.prevent="cliTab(2)" v-bind:class="{active: activeTab == 2}">Runs</a></li>
            <li><a v-on:click.stop.prevent="cliTab(3)" v-bind:class="{active: activeTab == 3}">Tracing</a></li>
        </ul>
    </div>
    <div id="detail-panel">
        <div v-if="activeTab == 0">
            <div id="actionlist">
                <AppAction
                    v-if="hasActions"
                    v-for="a in actions"
                    :key="a.name"
                    :action="a"
                    @actionClicked="setCurrentAction(a)"
                />
                <div id="newaction" v-on:click="addNewAction">
                    <i class="fa fa-plus-square"></i>&nbsp;new
                </div>
            </div>
            <AppActionPanel
                v-if="hasCurrentAction"
                :action-object="currentAction"
                @action-list-refresh-needed="refreshActionList"
            />
        </div>
        <div v-else-if="activeTab == 1">
            viriables
        </div>
        <div v-else-if="activeTab == 2">
            <AppRunInfo
                :case-info="this.caseInfo"
            />
        </div>
        <div v-else-if="activeTab == 3">
            Data tracing
        </div>
    </div>
</div>
</template>

<script>
import AppAction from './app-action.vue'
import AppActionPanel from './app-action-panel.vue'
import AppRunInfo from './app-run-info.vue'
import $ from 'jquery'

var AppDetail = Vue.extend({
    components: {AppAction, AppActionPanel, AppRunInfo},
    props: {
        caseInfo: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            // data that will change in this page
            localCaseInfo: {
                path: this.caseInfo.path,
                name: this.caseInfo.name,
                desc: this.caseInfo.desc
            },
            isNew: this.caseInfo.isNew,
            activeTab: 0,
            actions: [],
            currentAction: {},
        }
    },
    watch: {
        // watch prop:caseInfo
        // since we use local copy of prop here,
        // this change is only caused by user selecting another case in the list
        caseInfo: function (newCaseInfo) {
            console.log('prop:caseInfo changed. new name: ' + newCaseInfo.name)

            // update local variables. will discard unsaved local changes
            this.localCaseInfo.path = newCaseInfo.path
            this.localCaseInfo.name = newCaseInfo.name
            this.localCaseInfo.desc = newCaseInfo.desc
            this.isNew = newCaseInfo.isNew
            this.refreshActionList(true)
        }
    },
    computed: {
        hasActions: function () {
            return this.actions !== null && this.actions.length > 0
        },
        hasCurrentAction: function () {
            if (!this.hasActions) {
                this.currentAction = {}
            }
            return this.currentAction.name !== undefined
        }
    },
    mounted() {
        this.refreshActionList(true)
    },
    methods: {
        cliTab(item) {
            this.activeTab = item
        },
        setCurrentAction(act) {
            console.log('set current action to: ' + JSON.stringify(act))
            this.currentAction = act
        },
        addNewAction() {
            var newAction = {
                case_path: this.caseInfo.path,
                case_name: this.caseInfo.name,
                name: "action-new",
                desc: "this is a new action.",
                snippet: "",
                isNew: true
            }
            console.log('add action at list: ' + JSON.stringify(newAction))
            this.actions.push(newAction)
        },
        refreshActionList(toCloseActionPanel) {
            console.log('refreshing action list')
            // retrieve the actions again
            var self = this
            var url = `/case/${encodeURI(self.caseInfo.path)}/${encodeURI(self.caseInfo.name)}/actions`
            $.get(url, function (data) {
                if (data && data.length > 0) {
                    self.actions = data
                }
                else {
                    self.actions = []
                }
            })

            if (toCloseActionPanel) {
                self.currentAction = {}
            }
        },
        saveCase() {
            if (this.isNew) {
                console.log('create new case')
                $.post("/cases", JSON.stringify(this.localCaseInfo), function (data) {
                    console.log(JSON.stringify(data))
                    this.isNew = false
                }, "json")
            }
            else {
                console.log('update existing one')
                var oldCaseUrl = `/case/${encodeURI(this.caseInfo.path)}/${encodeURI(this.caseInfo.name)}`
                $.ajax({
                    url: oldCaseUrl,
                    type: 'PUT',
                    data: JSON.stringify(this.localCaseInfo),
                    success: function (resp) {
                        console.log(JSON.stringify(resp))
                    },
                    error: function (resp) {
                        console.log(JSON.stringify(resp))
                    }
                })
            }
        },
        runCase() {
            if (this.isNew) {
                alert('cannot run an unsaved case')
                console.log('cannot run an unsaved case')
                return
            }

            console.log('to run case: ' + this.caseInfo.name)
            var url = `/case/${encodeURI(this.caseInfo.path)}/${encodeURI(this.caseInfo.name)}/runs`
            var self = this
            $.ajax({
                url: url,
                type: 'POST',
                success: function (resp) {
                    console.log('success: add case into run-queue: ' + JSON.stringify(resp))
                }
            })
        }
    }
})
export default AppDetail
</script>

<style scoped>
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
    font-weight: 200;
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
    font-weight: 300;
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

div#detail-panel {
    width: 1000px;
    min-width: 1000px;
}

div#actionlist {
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
</style>

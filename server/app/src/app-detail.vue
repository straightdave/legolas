<!-- app detail -->
<template>
<div id="app-detail">
    <div id="detail-header">
        <p><input id="path" type="text" v-model="caseInfo.path"></input></p>
        <p><input id="name" type="text" v-model="caseInfo.name"></input></p>
        <p><input id="desc" type="text" v-model="caseInfo.desc"></input></p>
        <button v-on:click="saveCase">Save</button>
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
                    @actionClicked="viewAction(a.name)"
                />
                <div id="newaction" v-on:click="addNewAction">
                    <i class="fa fa-plus-square"></i>&nbsp;new
                </div>
            </div>
            <AppActionPanel v-if="hasCurrentAction" :action-object="currentAction"></AppActionPanel>
        </div>
        <div v-else-if="activeTab == 1">
            viriables
        </div>
        <div v-else-if="activeTab == 2">
            Runs
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
import $ from 'jquery'

var AppDetail = Vue.extend({
    components: {AppAction, AppActionPanel},
    props: {
        caseInfo: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            activeTab: 0,
            actions: [],
            currentAction: {},

            caseInfoCopy: {}
        }
    },
    watch: {
        // watch prop:caseInfo
        // probably caused by user selecting another case in the list
        caseInfo: function (newCaseInfo) {
            console.log('prop:caseInfo changed')
            var self = this
            var url = `/case/${encodeURI(self.caseInfo.path)}/${encodeURI(self.caseInfo.name)}/actions`
            $.get(url, function (data) {
                self.actions = data
            })
        }
    },
    mounted() {
        this.caseInfoCopy = JSON.parse(JSON.stringify(this.caseInfo))
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
    methods: {
        cliTab: function (item) {
            this.activeTab = item
        },
        viewAction: function (name) {
            this.currentAction = { name: name }
        },
        addNewAction: function () {
            this.actions.push({
                name: "newaction"
            })
        },
        saveCase() {
            alert('save')
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

div#detail-header button {
    height: 30px;
    width: 60px;
    font-size: 15px;
    font-weight: 200;
    border: none;
    background-color: #00B140;
    color: #fff;
    position: absolute;
    right: 10px;
    top: 5px;
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

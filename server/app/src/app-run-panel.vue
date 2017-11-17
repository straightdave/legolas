<!--
    run detail panel
    input: runObject, JSON of caseRun {
        uid: run id
        path: case path
        name: case name
        stati: [] of {
            name: state name string
            at: begin time
            message: message words
        }
        action_names: [] of string. names of actions
    }
-->
<template>
<div id="run-panel" v-if="runObject">
    <div id="run-title-box">
        <div id="title" v-on:click="toggleStati()">{{ runObject.uid }}</div>
        <div id="stati-box" v-if="toShowStati">
            <ul class="blocklist">
                <li v-for="s in runObject.stati">{{ s.at }} - {{ s.name }} {{ s.message }}</li>
            </ul>
        </div>
    </div>

    <div id="actions-log-box" v-if="hasActions">
        <div id="action-log" v-for="a in latestActions">
            <div id="action-name">
                <span><i class="fa fa-file-code-o"></i></span>
                <span>{{ a.name }}</span>
            </div>

            <div id="action-output">
                {{ a.output }}
            </div>
        </div>
    </div>
</div>
</template>

<script>
import VueResource from 'vue-resource'
Vue.use(VueResource)

var AppRunPanel = Vue.extend({
    props: {
        runObject: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            // flag indicates to show or not show stati box
            toShowStati: false,

            // array to store the lastest data of all actions, each element per one action
            latestActions: this.getAllActionState()
        }
    },
    computed: {
        hasActions() {
            return this.runObject.action_names && this.runObject.action_names.length > 0
        }
    },
    watch: {
        runObject(newRunObject) {
            // user clicks other runs
            this.toShowStati = false,
            this.latestActions = this.getAllActionState()
        }
    },
    methods: {
        toggleStati() {
            this.toShowStati = !this.toShowStati
        },
        getAllActionState() {
            // call API to get all action state of one case run
            // this happens rarely to reduce network calls
            return [
                {
                    name: "fake-action-1",
                    output: "Fugiat duis sunt et aliquip id exercitation occaecat in ut sint duis in laborum sunt do in voluptate cillum voluptate eu proident in anim ullamco sit dolor ut."

                },
                {
                    name: "fake-action-2",
                    output: "Lorem ipsum dolor in eiusmod mollit sed aute id ut officia officia deserunt ut dolore ad anim nulla aliqua id ad sed exercitation."
                }
            ]
        }
    }
})
export default AppRunPanel
</script>

<style scoped>
div#run-panel {
    width: 600px;
    float: left;
    margin-left: 10px;
}

div#run-title-box {
    margin-bottom: 20px;
}

div#run-title-box div#title {
    font-size: 20px;
    margin-bottom: 5px;
    cursor: pointer;
}

div#run-title-box div#stati-box {
    padding: 5px;
    background-color: #ececec;
}

ul.blocklist {
    list-style: none;
}

ul.blocklist li {
    display: block;
}

div#action-log {
    margin-bottom: 10px;
    padding: 5px;
}
</style>

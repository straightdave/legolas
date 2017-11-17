<!--
    case runs info: whole panel including run list, when switch to 'runs' tab
    input: case info
-->
<template>
<div id="app-info-panel">
    <div id="run-list">
        <div v-if="hasRuns()">
            <a v-for="r in runs" :key="r.uid" v-on:click.stop.prevent="showRunDetail(r)">
                {{ r.uid }}
            </a>
        </div>
        <div v-else>
            No Runs Yet
        </div>
    </div>
    <AppRunPanel :run-object="currentRun" />
</div>
</template>

<script>
import AppRunPanel from './app-run-panel.vue'
import VueResource from 'vue-resource'
Vue.use(VueResource)

var AppRunInfo = Vue.extend({
    components: {AppRunPanel},
    props: {
        caseInfo: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            path: this.caseInfo.path,
            name: this.caseInfo.name,
            isNew: this.caseInfo.isNew,
            runs: [],
            currentRun: {}
        }
    },
    mounted() {
        this.getRuns()
        this.currentRun = {}
    },
    watch: {
        caseInfo(newCaseInfo) {
            // behavior after user clicks other case items in the list
            this.path = newCaseInfo.path
            this.name = newCaseInfo.name
            this.isNew = newCaseInfo.isNew
            this.getRuns()
            this.currentRun = {}
        }
    },
    methods: {
        hasRuns() {
            return this.runs && this.runs.length > 0
        },

        showRunDetail(run) {
            this.currentRun = run
        },

        getRuns() {
            if (this.isNew) {
                console.log('no run info for new case')
                return
            }

            var url = `/case/${encodeURI(this.path)}/${encodeURI(this.name)}/runs`
            this.$http.get(url).then(
                resp => {
                    var data = resp.body
                    console.log('got runs: ' + JSON.stringify(data))
                    if (data && data.length > 0) {
                        this.runs = data
                    }
                    else {
                        this.runs = []
                    }
                },
                resp => {
                    console.log('http get failed, resp: ' + resp)
                    this.runs = []
                }
            )
        }
    }
})
export default AppRunInfo
</script>

<style scoped>
.green {
    color: #00B140;
}
.aborted {
    color: red;
}
.inprogress {
    color: #0366d6;
}

div#run-list {
    width: 300px;
    float: left;
}

div#run-list a {
    padding: 0;
    margin: 0;
    display: block;
    margin-bottom: 10px;
    cursor: pointer;
}



</style>

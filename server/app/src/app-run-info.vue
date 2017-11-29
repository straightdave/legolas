<template>
<div id="app-info-panel">
    <div id="run-list">
        <div v-if="hasRuns">
            <a v-for="r in runs"
               :key="r.uid"
               v-on:click.stop.prevent="showRunDetail(r)">{{ r._id }}</a>
        </div>
        <div v-else>
            No run yet
        </div>
    </div>
    <AppRunDetail v-if="currentRun" :run-object="currentRun" />
</div>
</template>

<script>
import Vue from 'vue'
import AppRunDetail from './app-run-detail.vue'
import VueResource from 'vue-resource'
Vue.use(VueResource)

var AppRunInfo = Vue.extend({
    components: {AppRunDetail},
    props: {
        caseInfo: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            isNew: this.caseInfo.isNew,
            runs: [],
            currentRun: null
        }
    },
    mounted() {
        if (!this.isNew) {
            this.runs = this.getRuns(this.caseInfo._id)
        }
    },
    watch: {
        caseInfo(newCaseInfo) {
            this.isNew = newCaseInfo.isNew
            this.currentRun = null
            this.runs = []

            if (!this.isNew) {
                this.runs = this.getRuns(this.caseInfo._id)
            }
        }
    },
    computed: {
        hasRuns() {
            return this.runs && this.runs.length > 0
        }
    },
    methods: {
        showRunDetail(runObject) {
            console.log(`click to show details of run:${runObject._id}`)
            this.currentRun = runObject
        },
        getRuns(case_id) {
            console.log(`getting runs of case:${case_id}`)
            this.$http.get(`/case/${encodeURI(case_id)}/runs`).then(
                resp => {
                    var data = resp.body
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

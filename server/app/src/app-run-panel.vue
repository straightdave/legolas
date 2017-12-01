<template>
<div id="run-panel">
    <div id="job-state-list">
        <div id="action-name">

        </div>



    </div>
</div>
</template>

<script>
import Vue from 'vue'
import VueResource from 'vue-resource'
Vue.use(VueResource)

var AppRunDetail = Vue.extend({
    props: {
        runObject: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            jobStates: []
        }
    },
    mounted() {
        console.log('run detail panel clicked, runObj = ' + JSON.stringify(this.runObject))

    },
    watch: {
        runObject(newRunObject) {
            console.log('use click run: ' + JSON.stringify(newRunObject))
            this.refreshRunDetails(newRunObject._id)
        }
    },
    methods: {
        refreshRunDetails(runId) {
            console.log(`getting job states of run:${runId}`)
            this.$http.get(`/run/${encodeURI(runId)}/jobstates`).then(
                resp => {
                    var data = resp.body
                    console.log('got jobStates: ' + JSON.stringify(data))
                    if (data && data.length > 0) {
                        this.jobStates = data
                    }
                    else {
                        this.jobStates = []
                    }
                },
                resp => {
                    console.log('http get failed, resp: ' + resp)
                    this.jobStates = []
                }
            )
        },

    }
})
export default AppRunDetail
</script>

<style scoped>
div#run-panel {
    width: 600px;
    float: left;
    margin-left: 10px;
}

</style>

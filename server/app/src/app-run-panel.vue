<template>
<div id="run-panel">
    <div id="run-title">
        <span>{{ runObject._id }}</span>
        <span id="refresh-btn"><a @click.stop.prevent="refreshRunDetails"><i class="fa fa-refresh"></i></a></span>
    </div>
    <div id="job-state-list">
        <div id="job-item" v-for="job in jobStates" :key="job.action_id">
            <div id="job-title">
                {{ job.action_name }} - {{ job.state }} - {{ job.error }}
            </div>
            <div>
                <pre>{{ job.output }}</pre>
            </div>
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
    watch: {
        runObject(newRunObject) {
            this.refreshRunDetails()
        }
    },
    mounted() {
        this.refreshRunDetails()
    },
    methods: {
        refreshRunDetails() {
            console.log(`getting job states of run: ${this.runObject._id}`)
            this.$http.get(`/run/${encodeURI(this.runObject._id)}/jobstates`).then(
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
                    console.log('getting jobstates failed, resp: ' + resp.body)
                }
            )
        },
    }
})
export default AppRunDetail
</script>

<style scoped>
div#run-panel {
    width: 800px;
    float: left;
    margin-left: 10px;
}

div#job-state-list {
    overflow: scroll;
}

div#job-title {
    color: #0366d6;
    font-weight: 300;
}

div#run-title {
    background-color: #ececec;
    font-family: 'Consolas', 'source-code-pro', monospace;
    font-size: 18px;
    margin-bottom: 5px;
    padding: 3px;
    position: relative;
}

span#refresh-btn {
    display: inline-block;
    position: absolute;
    right: 5px;

    font-size: 18px;
    font-weight: 200;
    cursor: pointer;
}

div#job-item > div {
    padding: 3px;
}

</style>

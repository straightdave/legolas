<!-- case runs info -->
<template>
<div id="app-info-panel">
    <div id="run-list">
        <div v-if="hasRuns()">
            <a v-for="r in runs" :key="r.uid" v-on:click.stop.prevent="showRunDetail(r.uid)">
                <span v-if="isStatus(r.uid, 'done')" class="green">
                    <i class="fa fa-check-circle"></i>&nbsp;
                </span>
                <span v-else-if="isStatus(r.uid, 'aborted')" class="aborted">
                    <i class="fa fa-stop-circle"></i>
                </span>
                <span v-else class="inprogress">
                    <i class="fa fa-play-circle"></i>
                </span>
                {{ r.uid }}
            </a>
        </div>
        <div v-else>
            No Runs Yet
        </div>
    </div>
    <AppRunPanel
        v-if="hasCurrentRun"
        :run-detail="currentRun"
    />
</div>
</template>

<script>
import $ from 'jquery'

var AppRunInfo = Vue.extend({
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
    },
    watch: {
        caseInfo: function (newCaseInfo) {
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
        isStatus(uid, stateStr) {
            var r
            for (var run in this.runs) {
                if (run.uid === uid) {
                    r = run
                    break
                }
            }

            console.log('got r to check status: ' + JSON.stringify(r))

            if (r !== undefined) {
                for (var s in r.stati) {
                    console.log('state: ' + JSON.stringify())
                    if (s.name === stateStr) {
                        return true
                    }
                }
            }

            return false
        },
        showRunDetail(uid) {
            console.log('show run details: ' + uid)
            this.currentRun.uid = uid
        },
        getRuns() {
            if (this.isNew) {
                console.log('no run info for new case')
                return
            }

            var url = `/case/${encodeURI(this.path)}/${encodeURI(this.name)}/runs`
            var self = this
            $.get(url, function (resp) {
                if (resp && resp.length > 0) {
                    console.log('success: got runs: ' + JSON.stringify(resp))
                    self.runs = resp
                }
                else {
                    console.log('got no run info')
                    self.runs = []
                }
            })
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

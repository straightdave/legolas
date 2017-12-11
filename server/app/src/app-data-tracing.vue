<template>
<div v-if="hasTracedData">
    <div id="chart-list" v-for="td in tracedData">
        <vue-chart
            type="line" :data="td"
            :options="{responsive: false}"
            :width="800" />
    </div>
</div>
</template>

<script>
import Vue from 'vue'
import VueChart from 'vue-chart-js'
Vue.use(VueChart)

var AppDataTracing = Vue.extend({
    props: {
        runObjects: {
            type: Array,
            required: true
        }
    },
    data() {
        return {
            tracedData: [],
            hasTracedData: false
        }
    },
    watch: {
        runObjects(newRuns) {
            console.log('use clicked another case (so new runs)')
            this.initTracedData()
            this.hasTracedData = this.tracedData.length > 0
        }
    },
    created() {
        console.log('chart created')
        this.initTracedData()
        this.hasTracedData = this.tracedData.length > 0
    },
    methods: {
        initTracedData() {
            this.tracedData.length = 0
            console.log('init traced data: ' + JSON.stringify(this.tracedData))

            var arrayRunLabels = []
            var arrayNames = {} // dict of keys
            for (var run of this.runObjects) {
                arrayRunLabels.push( (new Date(run.started_at)).toLocaleString() )

                if (run.traced_data) {
                    for (var k of Object.keys(run.traced_data)) {
                        arrayNames[k] = 0
                    }
                }
            }
            arrayNames = Object.keys(arrayNames)

            if (arrayNames.length === 0) {
                console.log('no data being traced. ingored')
                return
            }

            for (var name of arrayNames) {
                var data_array = []
                for (var run of this.runObjects) {
                    if (run.traced_data[name]) {
                        data_array.push(run.traced_data[name])
                    }
                    else {
                        data_array.push(0) // if no value, use 0 instead for now
                    }
                }
                console.log(`got data array for name:${name} -> ${JSON.stringify(data_array)}`)
                var _t = {
                    labels: arrayRunLabels,
                    datasets: []
                }
                _t.datasets.push({
                    label: name,
                    data: data_array
                })

                this.tracedData.push(_t)
            }
            console.log('new data: ' + JSON.stringify(this.tracedData))
        }
    }
})
export default AppDataTracing
</script>

<template>
<div>
    <div id="one-traced" v-for="td in tracedData">
        <vue-chart type="line" :data="td"></vue-chart>
    </div>
</div>
</template>

<script>
import Vue from 'vue'
import VueResource from 'vue-resource'
import VueChart from 'vue-chart-js'
Vue.use(VueResource)
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
            tracedData: []
        }
    },
    watch: {
        runObjects(newRuns) {
            console.log('use clicked another case (so new runs)')
            this.initTracedData()
        }
    },
    mounted() {
        console.log('chart mounted')
        this.initTracedData()
    },
    methods: {
        initTracedData() {
            this.tracedData = []

            var arrayRunIds = []
            var arrayNames = {} // dict of key - 0
            for (var run of this.runObjects) {
                arrayRunIds.push(run._id)

                if (run.traced_data) {
                    for (var k of Object.keys(run.traced_data)) {
                        arrayNames[k] = 0
                    }
                }
            }
            arrayNames = Object.keys(arrayNames)

            for (var name of arrayNames) {
                var data_array = []
                for (var run of this.runObjects) {
                    if (run.traced_data[name]) {
                        data_array.push(run.traced_data[name])
                    }
                    else {
                        data_array.push(0)
                    }
                }
                console.log(`got data array for name:${name} -> ${JSON.stringify(data_array)}`)
                var _t = {
                    labels: arrayRunIds,
                    datasets: []
                }
                console.log("_t= " + JSON.stringify(_t))
                _t.datasets.push({
                        label: name,
                        data: data_array
                })
                console.log("_t= " + JSON.stringify(_t))

                this.tracedData.push(_t)
            }

            console.log('data: ' + JSON.stringify(this.tracedData))
        }
    }
})
export default AppDataTracing
</script>


<style>

</style>

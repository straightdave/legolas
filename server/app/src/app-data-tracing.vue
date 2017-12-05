<template>
<div>
    <div class="namebox">
        Some Data 1
    </div>
    <vue-chart type="line" :data="myData"></vue-chart>

    <div class="namebox">
        Some Data 2
    </div>
    <vue-chart type="line" :data="myData"></vue-chart>

    <div class="namebox">
        Some Data 3
    </div>
    <vue-chart type="line" :data="myData"></vue-chart>
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
        caseObject: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            myData: {
                labels: ['run 1', 'run 2', 'run 3', 'run 4', 'run 5'],
                datasets: [
                    {
                        label: 'Some Data',
                        data: [10, 20, 0, 30, 40]
                    }
                ]
            }
        }
    },
    watch: {
        caseObject(newCaseObject) {
            console.log('use clicked another case')
        }
    },
    mounted() {
        console.log('chart mounted')
    },
    methods: {

    }
})
export default AppDataTracing
</script>


<style>

</style>

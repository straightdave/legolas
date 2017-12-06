<!-- run listed item -->
<template>
<div id="run" v-on:click="$emit('run-clicked', run)" :class="{red: color === 'red', green: color === 'green'}">
    <div id="runtitlebar">
        <span id="name">{{ run._id }}</span>
        <span id="desc"><i class="fa fa-clock-o"></i> {{ (new Date(run.started_at)).toLocaleString() }}</span>
    </div>
</div>
</template>

<script>
import Vue from 'vue'

var AppRun = Vue.extend({
    props: {
        run: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            color: "nocolor"
        }
    },
    mounted() {
        if (this.run.output === 'failed' || this.run.output === 'aborted') {
            this.color = 'red'
        }
        else if (this.run.output === 'done') {
            this.color = 'green'
        }
    }
})
export default AppRun
</script>

<style scoped>
.green {
    border-left: solid 10px #00B140 !important;
}
.red {
    border-left: solid 10px red !important;
}

div#run {
    margin-bottom: 5px;
    padding:10px;
    cursor: pointer;
    background-color: #ececec;
    overflow: hidden;
    overflow: hidden;
    border-left: solid 10px gray;
}

div#run span {
    display: block;
}

div#runtitlebar span#name {
    font-family: 'Consolas', 'source-code-pro', monospace;
    font-size: 18px;
}
div#runtitlebar span#desc {
    font-size: small;
    color: gray;
}
</style>

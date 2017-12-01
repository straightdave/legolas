<!-- run listed item -->
<template>
<div id="run" v-on:click="$emit('run-clicked', run)">
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
    }
})
export default AppRun
</script>

<style>
div#run {
    margin-bottom: 5px;
    padding:10px;
    cursor: pointer;
    background-color: #ececec;
    overflow: hidden;
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

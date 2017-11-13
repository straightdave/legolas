<!-- case list -->
<template>
<div id="app-sidebar">
    <app-search></app-search>
    <div class="titlebar">
        <span>My Cases</span>
        <a id="create" v-on:click.stop.prevent="create"><i class="fa fa-plus"></i></a>
    </div>

    <div id="list" v-if="cases.length">
        <app-case
            v-for="c in cases"
            :key="c.id"
            :case="c"
            @view-case="$emit('case-clicked', c.id)"
        />
    </div>
</div>
</template>

<script>
import AppSearch from './app-search.vue'
import AppCase from './app-case.vue'

let nextItemId = 1

export default {
    components: {AppSearch, AppCase},
    data() {
        return {
            cases: [
                {
                    id: nextItemId++,
                    name: 'hello',
                    lastResult: 'passed',
                    group: "group-a"
                },
                {
                    id: nextItemId++,
                    name: 'goodbye',
                    lastResult: 'failed',
                    group: "group-a"
                },
                {
                    id: nextItemId++,
                    name: 'go back home',
                    lastResult: 'unknown',
                    group: "group-b"
                }
            ]
        }
    },
    methods: {
        create: function () {
            this.cases.append({
                id: nextItemId++,
                name: 'new one'
            })
        }
    }
}
</script>


<style scoped>
div#app-sidebar {
    width: 340px;
    border-right: solid 4px #ececec;
    float: left;
}

div.titlebar {
    width: 300px;
    margin: 0 auto;
}

div.titlebar span {
    display: inline-block;
    margin-right: 10px;
    font-weight: 400;
    padding-bottom: 8px;
}

a#create {
    display: inline-block;
    text-decoration: none;
    color: gray;
    cursor: pointer;
    font-size: 0.7em;
    font-weight: 200px;
}

div#list {
    background-color: #ececec;
    padding: 4px 0;
}
</style>

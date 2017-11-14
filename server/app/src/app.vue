<!-- app -->
<template>
<div id="app">
    <app-nav></app-nav>
    <div>
        <app-sidebar @case-clicked="showCaseDetail"></app-sidebar>
        <app-detail v-if="hasCInfo" :case-info="cinfo"></app-detail>
    </div>
</div>
</template>

<script>
import AppNav from './app-nav.vue'
import AppSidebar from './app-sidebar.vue'
import AppDetail from './app-detail.vue'

var App = Vue.extend({
    components: {AppNav, AppSidebar, AppDetail},
    methods: {
        showCaseDetail: function (c) {
            this.cinfo = c
        }
    },
    data() {
        return {
            cinfo: {}
        }
    },
    computed: {
        hasCInfo: function () {
            return this.cinfo.name !== undefined
        }
    }
})
export default App
</script>

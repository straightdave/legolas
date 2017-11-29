<template>
<div>
    <app-sidebar :whatFor="whatFor" @item-clicked="showCaseDetail"></app-sidebar>
    <app-detail v-if="hasCInfo" :case-info="cinfo"></app-detail>
</div>
</template>

<script>
import Vue from 'vue'
import AppSidebar from './app-sidebar.vue'
import AppDetail from './app-detail.vue'

var AppCaseMain = Vue.extend({
    components: {AppSidebar, AppDetail},
    methods: {
        showCaseDetail: function (c) {
            this.cinfo = c
        }
    },
    data() {
        return {
            whatFor: 'cases',
            cinfo: {}
        }
    },
    computed: {
        hasCInfo() {
            return this.cinfo.name !== undefined
        }
    }
})
export default AppCaseMain
</script>

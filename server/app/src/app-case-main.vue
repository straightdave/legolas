<template>
<div>
    <app-sidebar :what-for="whatFor" @item-clicked="showCaseDetail"></app-sidebar>
    <app-case-detail v-if="hasCInfo" :case-info="currentCase"></app-case-detail>
</div>
</template>

<script>
import Vue from 'vue'
import AppSidebar from './app-sidebar.vue'
import AppCaseDetail from './app-case-detail.vue'

var AppCaseMain = Vue.extend({
    components: {AppSidebar, AppCaseDetail},
    data() {
        return {
            whatFor: 'cases',
            currentCase: null
        }
    },
    computed: {
        hasCInfo() {
            return this.currentCase !== null
        }
    },
    methods: {
        showCaseDetail(caseObject) {
            this.currentCase = caseObject
        }
    }
})
export default AppCaseMain
</script>

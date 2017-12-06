<template>
<div>
    <app-sidebar :what-for="whatFor" @item-clicked="showCaseDetail" @close-right-panel="closeRightPanel"></app-sidebar>
    <app-case-detail v-if="hasCInfo" :case-info="currentCase" @refresh-sidebar-list="refreshSidebarList"></app-case-detail>
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
        },
        refreshSidebarList() {
            // TODO: for now, refresh the whole page (will miss new action data)
            window.location.reload()
        },
        closeRightPanel() {
            this.currentCase = null
        }
    }
})
export default AppCaseMain
</script>

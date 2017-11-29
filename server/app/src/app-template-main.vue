<template>
<div>
    <app-sidebar :what-for="whatFor" @item-clicked="showDetail"></app-sidebar>
    <app-template-detail v-if="hasTpl" :template-object="tplObject"></app-template-detail>
</div>
</template>

<script>
import Vue from 'vue'
import AppSidebar from './app-sidebar.vue'
import AppTemplateDetail from './app-template-detail.vue'

var AppTemplateMain = Vue.extend({
    components: {AppSidebar, AppTemplateDetail},
    data() {
        return {
            whatFor: 'templates',
            tplObject: null
        }
    },
    computed: {
        hasTpl() {
            return this.tplObject !== null
        }
    },
    methods: {
        showDetail(tplObject) {
            this.tplObject = tplObject
        }
    }
})
export default AppTemplateMain
</script>

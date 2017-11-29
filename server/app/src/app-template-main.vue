<template>
<div>
    <app-sidebar :whatFor="whatFor"></app-sidebar>
</div>
</template>

<script>
import Vue from 'vue'
import AppSidebar from './app-sidebar.vue'

var AppTemplateMain = Vue.extend({
    components: {AppSidebar},

    data() {
        return {
            whatFor: 'templates'
        }
    }
})
export default AppTemplateMain
</script>

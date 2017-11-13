<!-- app content -->
<template>
<div id="app-content">
    <h1>{{ caseData.name }}</h1>
</div>
</template>

<script>
import AppDetail from './app-detail.vue'

var AppContent = Vue.extend({
    props: {
        caseData: {
            type: Object,
            required: true
        }
    },
    components: {AppDetail}
})
export default AppContent
</script>

<style>
div#app-content {
    float: left;
    padding: 10px;
}
</style>

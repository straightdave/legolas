<!-- app -->
<template>
<div id="app">
    <app-nav></app-nav>
    <div>
        <app-sidebar @case-clicked="caseClicked"></app-sidebar>
        <app-content :case-data="c"></app-content>
    </div>
</div>
</template>

<script>
import AppNav from './app-nav.vue'
import AppSidebar from './app-sidebar.vue'
import AppContent from './app-content.vue'

var App = Vue.extend({
    components: {AppNav, AppSidebar, AppContent},
    methods: {
        caseClicked: function (id) {
            this.c.name = id
        }
    },
    data() {
        return {
            c: {
                name: "test"
            }
        }
    }
})
export default App
</script>

<template>
<div id="template-detail">
    {{ templateObject.name }}
</div>
</template>

<script>
import Vue from 'vue'
import VueResource from 'vue-resource'
Vue.use(VueResource)

var AppTemplateDetail = Vue.extend({
    props: {
        templateObject: {
            type: Object,
            required: true
        }
    }
})
export default AppTemplateDetail
</script>

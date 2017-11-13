<!-- action item -->
<template>
<div id="action" v-on:click="$emit('viewAction', name)">
    <div id="titlebar">
        <p><span>{{ name }}</span></p>
        <p><span>$/{{ casePath }}/{{ caseName }}#{{ name }}</span></p>
    </div>
</div>
</template>

<script>

var AppAction = Vue.extend({
    props: {
        action: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            name: this.action.name,
            casePath: this.action.casePath,
            caseName: this.action.caseName,
            index: this.action.index
        }
    }
})
export default AppAction

</script>

<style>
div#action {
    width: 100%;
    height: 100px;
    border-bottom: solid 2px #ececec;
    margin-bottom: 5px;
    padding:3px;
}
</style>

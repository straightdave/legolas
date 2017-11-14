<!-- action item -->
<template>
<div id="action" v-on:click="$emit('actionClicked', action.name)">
    <div id="titlebar">
        <p><span id="name">{{ action.name }}</span></p>
        <p><span id="desc">This is an action</span></p>
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
    }
})
export default AppAction

</script>

<style>
div#action {
    height: 60px;
    margin-bottom: 5px;
    padding:10px;
    cursor: pointer;
    background-color: #ececec;
    overflow: hidden;
}

div#titlebar span#name {
    font-size: 20px;
    font-weight: 200;
}
div#titlebar span#desc {
    color: gray;
    font-weight: 200;
}
</style>

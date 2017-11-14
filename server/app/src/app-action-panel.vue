<!-- app action panel -->
<template>
<div id="action-panel" v-if="actionObject">
    <div><input type="text" v-model="actionObject.name"></input></div>
    <div><input type="text" v-model="actionObject.desc"></input></div>
    <div id="snippet">
        <textarea v-model="actionObject.snippet" rows="25" placeholder="Snippet..."></textarea>
    </div>
</div>
</template>

<script>

var AppActionPanel = Vue.extend({
    props: {
        actionObject: {
            type: Object,
            required: true
        }
    },
    watch: {
        actionObject: function (newActionObject) {
            console.log('changed to ' + newActionObject.name)
        }
    }
})
export default AppActionPanel

</script>

<style scoped>
div#action-panel {
    width: 600px;
    float: left;
    margin-left: 10px;
}
div#action-panel > div {
    margin-bottom: 10px;
}
div#action-panel input[type="text"]
 {
    font-size: 15px;
    padding:3px;
    width: 200px;
    border: solid 2px #ececec;
    background-color: #fff;
}
div#action-panel textarea {
    font-size: 16px;
    margin-bottom: 5px;
    border: solid 2px #ececec;
    width: 100%;
    min-width: 500px;
}
</style>

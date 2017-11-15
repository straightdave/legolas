<!-- app action panel -->
<template>
<div id="action-panel" v-if="localActionObject">
    <div><input type="text" v-model="localActionObject.name" placeholder="Action name"></input></div>
    <div><input type="text" v-model="localActionObject.desc" placeholder="Desc goes here..."></input></div>
    <div id="snippet">
        <textarea v-model="localActionObject.snippet" rows="25" placeholder="Snippet..."></textarea>
    </div>
    <div>
        <button v-on:click="saveAction">Save</button>
        <button v-on:click="deleteAction">Delete</button>
    </div>
</div>
</template>

<script>
import $ from 'jquery'

var AppActionPanel = Vue.extend({
    props: {
        actionObject: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            localActionObject: {
                cpath: this.actionObject.cpath,
                cname: this.actionObject.cname,
                name: this.actionObject.name,
                desc: this.actionObject.desc,
                snippet: this.actionObject.snippet
            },
            isNew: this.actionObject.isNew
        }
    },
    watch: {
        // caused by use selecting another action in the list
        actionObject: function (newActionObject) {
            console.log('changed to ' + newActionObject.name)

            this.localActionObject.cpath = newActionObject.cpath
            this.localActionObject.cname = newActionObject.cname
            this.localActionObject.name = newActionObject.name
            this.localActionObject.desc = newActionObject.desc
            this.localActionObject.snippet = newActionObject.snippet
            this.isNew = newActionObject.isNew
        }
    },
    methods: {
        saveAction: function () {
            if (this.isNew) {
                console.log('save new action')
                var data = {
                    case_path: this.localActionObject.cpath,
                    case_name: this.localActionObject.cname,
                    name: this.localActionObject.name,
                    desc: this.localActionObject.desc,
                    snippet: this.localActionObject.snippet
                }
                $.post('/actions', JSON.stringify(data), function (resp) {
                    console.log('success create action: ' + JSON.stringify(resp))
                    this.isNew = false
                }, "json")
            }
            else {
                console.log('save existing action')
            }
            this.$emit('action-list-refresh-needed')
        },
        deleteAction: function () {
            console.log('to delete this action')

            this.$emit('action-list-refresh-needed')
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
    border: none;
}
div#action-panel textarea {
    font-size: 16px;
    margin-bottom: 5px;
    border: solid 2px #ececec;
    width: 100%;
    min-width: 500px;
}

div#action-panel button {
    height: 35px;
    font-size: 15px;
    font-weight: 200;
    width: 80px;
    border: none;
    background-color: #ececec;
    cursor: pointer;
}
</style>

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
                case_path: this.actionObject.case_path,
                case_name: this.actionObject.case_name,
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
            console.log('changed to ' + JSON.stringify(newActionObject))

            this.localActionObject.case_path = newActionObject.case_path
            this.localActionObject.case_name = newActionObject.case_name
            this.localActionObject.name = newActionObject.name
            this.localActionObject.desc = newActionObject.desc
            this.localActionObject.snippet = newActionObject.snippet
            this.isNew = newActionObject.isNew
        }
    },
    methods: {
        saveAction: function () {
            var self = this
            if (this.isNew) {
                console.log('save new action')
                $.post('/actions', JSON.stringify(self.localActionObject), function (resp) {
                    console.log('response: ' + JSON.stringify(resp))
                    if (resp.error) {
                        console.log('error happened')
                        self.isNew = true
                    }
                    else {
                        console.log('save new succeeded')
                        self.isNew = false

                        // refresh action list and close current action panel
                        self.$emit('action-list-refresh-needed', true)
                    }
                }, "json")
            }
            else {
                console.log('save existing action')
                // TODO
            }
        },
        deleteAction: function () {
            if (this.isNew) {
                console.log('delete a new (unsafed) action')
                // refresh action list and close current action panel
                this.$emit('action-list-refresh-needed', true)
                return
            }

            var cpath = this.actionObject.case_path
            var cname = this.actionObject.case_name
            var name = this.actionObject.name
            console.log(`delete action: ${cpath}/${cname}#${name}`)
            var url = `/case/${encodeURI(cpath)}/${encodeURI(cname)}/${encodeURI(name)}`
            var self = this
            $.ajax({
                url: url,
                type: 'DELETE',
                success: function (data) {
                    console.log('resp: ' + JSON.stringify(data))
                    if (data.error) {
                        console.log('error')
                        // TODO: leave error messages somewhere on page
                    }
                    else {
                        console.log('succeeded')
                        // refresh list and close current action panel
                        self.$emit('action-list-refresh-needed', true)
                    }
                }
            })
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

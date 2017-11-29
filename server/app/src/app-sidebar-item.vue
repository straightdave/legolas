<!-- listed item -->
<template>
<div id="item" @click="$emit('view-item-detail')">
    <div id="item-inner">
        <span id="item-title">
            <span v-if="isNew" class="newitem">
                <i class="fa fa-circle"></i>
            </span>
            <span v-else class="green">
                <i class="fa fa-check"></i>
            </span>
            {{ itemName }}
        </span>
        <span id="item-group">{{ itemPath }}</span>
    </div>
</div>
</template>

<script>
export default {
    props: {
        itemObject: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            itemId:   this.itemObject._id,
            itemName: this.itemObject.name,
            itemPath: this.itemObject.path,
            isNew:    this.itemObject.isNew
        }
    }
}
</script>

<style scoped>
div#item {
    width: 100%;
    margin-top: 5px;
}

.green {
    color: #00B140;
}
.newitem {
    font-size: small;
    color: #4AACF2;
}

div#item-inner {
    height: 100px;
    background-color: #fff;
    border-left: solid 10px #00B140;
    padding-top: 10px;
    padding-left: 10px;
    cursor: pointer;
    position: relative;
}

span#item-title {
    display: block;
    font-size: 20px;
    overflow: hidden;
}

span#item-group {
    display: block;
    position: absolute;
    left: 10px;
    bottom: 10px;
    color: gray;
}
</style>

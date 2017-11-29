<!-- case listed item -->
<template>
<div id="case" @click="$emit('view-case', name)">
    <div id="case-inner">
        <span id="case-title">
            <span v-if="isNew" class="newitem">
                <i class="fa fa-circle"></i>
            </span>
            <span v-else class="green">
                <i class="fa fa-check"></i>
            </span>
            {{ name }}
        </span>
        <span id="case-group">{{ path }}</span>
    </div>
</div>
</template>

<script>
export default {
    props: {
        case: {
            type: Object,
            required: true
        }
    },
    data() {
        return {
            name: this.case.name,
            path: this.case.path,
            isNew: this.case.isNew
        }
    }
}
</script>

<style scoped>
div#case {
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

div#case-inner {
    height: 100px;
    background-color: #fff;
    border-left: solid 10px #00B140;
    padding-top: 10px;
    padding-left: 10px;
    cursor: pointer;
    position: relative;
}

span#case-title {
    display: block;
    font-size: 20px;
    overflow: hidden;
}

span#case-group {
    display: block;
    position: absolute;
    left: 10px;
    bottom: 10px;
    color: gray;
}
</style>

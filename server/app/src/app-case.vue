<!-- item -->
<template>
<div class="case" @click="$emit('view-case', id)">
    <div class="case-inner">
        <span class="case-title">
            <span v-if="lastResult == 'passed'" class="green">
                <i class="fa fa-check"></i>
            </span>
            {{ name }}
        </span>
        <span class="case-group">$/{{ group }}</span>
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
            id: this.case.id,
            name: this.case.name,
            lastResult: this.case.lastResult,
            group: this.case.group
        }
    }
}
</script>

<style scoped>
div.case {
    width: 100%;
    margin-top: 5px;
}

.green {
    color: #00B140;
}

div.case-inner {
    height: 100px;
    background-color: #fff;
    border-left: solid 10px #00B140;
    padding-top: 10px;
    padding-left: 10px;
    cursor: pointer;
    position: relative;
}


span.case-title {
    display: block;
    font-size: 20px;
}

span.case-group {
    display: block;
    position: absolute;
    left: 10px;
    bottom: 10px;
    color: gray;
}
</style>

<script>
export default {
    props: {
        msg: String // Riceve il messaggio dal componente padre
    },
    mounted() {
        this.startTimer();
    },
    watch: {
        msg() {
            this.startTimer();
        }
    },
    methods: {
        startTimer() {
            if (this.timeout) clearTimeout(this.timeout);
            if (this.msg) {
                this.timeout = setTimeout(() => {
                    this.$emit('clear');
                }, 5000);
            }
        }
    },
    unmounted() {
        if (this.timeout) clearTimeout(this.timeout);
    }
}
</script>
<template>
    <div class="error-box" v-if="msg">
        ⚠️ {{ msg }}
    </div>
</template>
<style scoped>
.error-box {
    background-color: #ffe6e6;
    color: #cc0000;
    padding: 10px;
    border: 1px solid #cc0000;
    border-radius: 5px;
    margin-top: 10px;
    font-weight: bold;
}
</style>

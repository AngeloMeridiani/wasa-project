<script>
    import ErrorMsg from '@/components/ErrorMsg.vue';

    export default {
        components: {
            ErrorMsg
        },
        data : function() {
            return {
                username : "",
                errormsg : "",
                loading: false
            }
        },
        methods: {
            async doLogin() {
                if (this.username == "" || this.username.length < 3) {
                    this.errormsg = "Il nome deve avere almeno 3 caratteri.";
                    return;
                }
                this.loading = true;
                this.errormsg = "";

                try {
                    let response = await this.$axios.post("/session", 
                        { name : this.username }
                    );
                    const identifier = response.data.identifier;

                    localStorage.setItem("token", identifier);
                    localStorage.setItem("username", this.username);

                    this.$router.push("/home");
                } 
                catch (e) {
                    if (e.response && e.response.status == 400) {
                        this.errormsg = "Nome non valido (deve essere compreso tra i 3 e 25 caratteri)";
                    } else {
                        this.errormsg = e.toString();
                    }
                    console.error("Errore di login:", e);
                } finally {
                    this.loading = false;
                }
            }
        }
    }
</script>

<template>
    <div class="login-container">
        <p>
            Inserire nome utente per accedere o registrarsi:
        </p>
        <input 
            type="text"
            v-model="username"
            placeholder="Il tuo nome:"
            @keyup.enter="doLogin"
            :disabled="loading"
        >
        <button @click="doLogin" style="margin-top: 10px;" :disabled="loading">
            <span v-if="loading">Attendi...</span>
            <span v-else>Accedi</span>
        </button>
        
        <ErrorMsg :msg="errormsg" />
    </div>
</template>

<style>

</style>
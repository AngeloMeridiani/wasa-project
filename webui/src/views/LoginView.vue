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
    <div class="login-page">
        <div class="login-box">
            <div class="login-header">
                <h2>WasaText</h2>
                <p>Inserisci il tuo nome per accedere o registrarti</p>
            </div>
            
            <div class="login-form">
                <input 
                    type="text"
                    class="input-custom login-input"
                    v-model="username"
                    placeholder="Il tuo nome"
                    @keyup.enter="doLogin"
                    :disabled="loading"
                    autocomplete="off"
                >
                <button 
                    class="btn-primary-custom login-btn" 
                    @click="doLogin" 
                    :disabled="loading"
                >
                    <span v-if="loading">Attendi...</span>
                    <span v-else>Accedi</span>
                </button>
            </div>
            
            <ErrorMsg v-if="errormsg" :msg="errormsg" class="login-error" />
        </div>
    </div>
</template>

<style scoped>
.login-page {
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--color-bg-app);
}

.login-box {
    background-color: var(--color-bg-white);
    padding: 40px;
    border-radius: var(--radius-xl);
    box-shadow: var(--shadow-md);
    width: 100%;
    max-width: 400px;
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.login-header {
    text-align: center;
}

.login-header h2 {
    color: var(--color-primary);
    font-weight: 600;
    margin-bottom: 8px;
    letter-spacing: -0.5px;
}

.login-header p {
    color: var(--color-text-muted);
    font-size: 0.95rem;
    margin: 0;
}

.login-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.login-input {
    width: 100%;
    padding: 12px 16px;
    font-size: 1rem;
}

.login-btn {
    width: 100%;
    padding: 12px;
    font-size: 1rem;
    font-weight: 600;
}

.login-error {
    margin-top: -8px;
}
</style>
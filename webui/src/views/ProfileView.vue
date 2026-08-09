<script>
import ErrorMsg from '@/components/ErrorMsg.vue'; 
    
export default {
    components: {
            ErrorMsg
    },
    data: function() {
        return {
            userId : "",
            username : "",
            errormsg : "",
            loading: false,
            photoUrl: ""
        }
    },
    methods: {
        async loadProfile() {
            try {
                this.loading = true;
                this.userId = localStorage.getItem("token");
                let response = await this.$axios.get(`/users/${this.userId}`);
                this.username = response.data.name;
            }
            catch(e) {
                this.errormsg = e.toString();
            }
            finally {
                this.loading = false;
            }
        },
        async setUserName(newUsername) {
            if (this.username == "" || this.username.length < 3) {
                this.errormsg = "Il nome deve avere almeno 3 caratteri.";
                return;
            }
            this.loading = true;
            this.errormsg = "";

            try {
                let response = await this.$axios.put(`/users/${this.userId}`, { name : newUsername });
                this.username = newUsername;
                localStorage.setItem("username", newUsername);
            }
            catch(e) {
                if (e.response && e.response.status == 400) {
                    this.errormsg = "Nome non valido (deve essere compreso tra i 3 e 25 caratteri)";
                } else {
                    this.errormsg = e.toString();
                }
                console.error("Errore di login:", e);
            }
            finally {
                this.loading = false;
            }
        },
        mounted() {
            this.loadProfile();
        }
    }
}
</script>

<template>
    <div class="profile-view">
        <h1>Profilo Utente</h1>
        <p>Qui vedrai le informazioni del profilo.</p>
        
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

        <div class="user-informations">
            <h3> Nome utente: {{ username }} </h3>
            <h3> ID Utente: {{ userId }} </h3>
            <h3> Immagine profilo utente: </h3>
            <img v-if="photoUrl" :src="photoUrl">
        </div>

        <div class="change-name">
            <input 
                type="text"
                class="input-custom login-input"
                v-model="username"
                placeholder="Il tuo nuovo nome..."
                @keyup.enter="setUserName(username)"
                :disabled="loading"
                autocomplete="off"
            >
            <button 
                class="btn-primary-custom login-btn" 
                @click="setUserName(username)" 
                :disabled="loading"
            >
                <span v-if="loading">Attendi...</span>
                <span v-else>Conferma</span>
            </button>
        </div> 
    </div>
</template>

<style scoped>
.profile-view {
    padding: 20px;
    max-width: 600px;
    margin: 0 auto;
}

/* Stile per la card delle informazioni utente */
.user-informations {
    display: flex;
    flex-direction: column;
    align-items: center;
    background-color: var(--color-bg-white);
    padding: 30px;
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
    margin-bottom: 30px;
    border: 1px solid var(--color-border);
}
</style>

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
            newUsername: "",
            errormsg : "",
            loading: false,
            newPhotoUrl: "",
            photoUrl: ""
        }
    },
    methods: {
        async loadProfile() {
            try {
                this.loading = true;
                this.userId = localStorage.getItem("token");
                
                // Se non c'è il token, ributtiamo l'utente al Login!
                if (!this.userId) {
                    this.$router.push("/");
                    return;
                }

                let response = await this.$axios.get(`/users/${this.userId}`);
                this.photoUrl = response.data.userPhoto;
                this.username = response.data.name;
                this.newUsername = response.data.name;
            }
            catch(e) {
                this.errormsg = e.toString();
            }
            finally {
                this.loading = false;
            }
        },
        async setUserName() {
            if (this.newUsername == "" || this.newUsername.length < 3) {
                this.errormsg = "Il nome deve avere almeno 3 caratteri.";
                return;
            }
            
            if (this.newUsername == this.username) {
                return;
            }

            this.loading = true;
            this.errormsg = "";

            try {
                let response = await this.$axios.put(`/users/${this.userId}/name`, { name : this.newUsername });
                this.username = this.newUsername;
                localStorage.setItem("username", this.newUsername);
                await this.loadProfile();
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
        async setPhoto() {
            this.loading = true;
            try {
                let response = await this.$axios.put(`/users/${this.userId}/photo`, { photoUrl : this.newPhotoUrl });
                this.photoUrl = response.data.photoUrl;
                await this.loadProfile();
            }
            catch(e) {
                this.errormsg = e.toString();
            }
            finally {
                this.loading = false;
            }
        }
    },
    mounted() {
        this.loadProfile();
    }
}
</script>

<template>
    <div class="profile-view">
        <h1 class="page-title">Profilo Utente</h1>
        
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
                v-model="newUsername"
                placeholder="Il tuo nuovo nome..."
                @keyup.enter="setUserName"
                :disabled="loading"
                autocomplete="off"
            >
            <button 
                class="btn-primary-custom login-btn" 
                @click="setUserName" 
                :disabled="loading"
            >
                <span v-if="loading">Attendi...</span>
                <span v-else>Conferma</span>
            </button>
        </div> 

        <div class="change-photo">
            <input 
                type="text"
                class="input-custom login-input"
                v-model="newPhotoUrl"
                placeholder="La tua nuova foto..."
                @keyup.enter="setPhoto"
                :disabled="loading"
                autocomplete="off"
            >
            <button 
                class="btn-primary-custom login-btn" 
                @click="setPhoto" 
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

.page-title {
    color: var(--color-primary);
    margin-bottom: 20px;
}

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

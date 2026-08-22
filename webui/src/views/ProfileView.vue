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
            photoUrl: "",
            selectedFile: null
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
                if (e.response && e.response.status == 409) {
                    this.errormsg = "Questo nome utente è già in uso!";
                } else if (e.response && e.response.status == 400) {
                    this.errormsg = e.response.data || "Nome non valido (deve essere compreso tra i 3 e 25 caratteri)";
                } else {
                    this.errormsg = e.toString();
                }
                console.error("Errore modifica nome:", e);
            }
            finally {
                this.loading = false;
            }
        },
        handlePhotoUpload(event) {
            this.selectedFile = event.target.files[0];
        },
        async setPhoto() {
            if (!this.selectedFile) {
                this.errormsg = "Seleziona prima un'immagine!";
                return;
            }
            this.loading = true;
            try {
                let response = await this.$axios.put(`/users/${this.userId}/photo`, this.selectedFile, {
                    headers: {
                        'Content-Type': this.selectedFile.type
                    }
                });
                this.photoUrl = response.data.photoUrl;
                await this.loadProfile();
                this.selectedFile = null;
                if (this.$refs.fileInput) this.$refs.fileInput.value = "";
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
            <h3> Nome utente: <strong>{{ username }}</strong> </h3>
            <h3> ID Utente: <strong>{{ userId }}</strong> </h3>
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
                type="file"
                ref="fileInput"
                class="input-custom login-input"
                accept="image/png, image/jpeg"
                @change="handlePhotoUpload"
                :disabled="loading"
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
    height: 100%;
    overflow-y: auto;
    padding-bottom: 80px;
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

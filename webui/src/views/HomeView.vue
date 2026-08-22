<script>
import ChatsList from '@/components/ChatsList.vue';
import ChatWindow from '@/components/ChatWindow.vue';
export default {
	components : {
		ChatsList,
		ChatWindow
	},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			chats : null,
			messages : null,
			sentMessages : [],
			selectedChatId: null,
			selectedChatTitle: '',
			selectedGroupId: null,
			newMessage: '',
			replyingToMessageId: null,
			forwardingMessageId: null,
			availableEmojis: ['👍', '❤️', '💔', '😴', '😂', '😮', '😢', '🔥', '🤐'],
			pollIntervalId: null,
			isCreatingGroup: false,
			isAddingUserToGroup: false,
			isChangingGroupName: false,
			newGroupNameInput: '',
			groupSearchQuery: '',
			groupSearchResults: [],
			newGroupData: {
				name: '',
				description: '',
				photoUrl: '',
				selectedMembers: []
			}
		}
	},
	methods: {
		async refresh(silent = false) {
			if (!silent) {
				this.loading = true;
				this.errormsg = null;
			}
			try {
				let response = await this.$axios.get("/chats");
				this.chats = response.data.chats || [];
			} catch (e) {
				if (!silent) this.errormsg = e.toString();
			}
			if (!silent) this.loading = false;
		},
		async CreateChat(targetUserName) {
			this.loading = true;
			this.errormsg = null;
			if (!targetUserName || targetUserName.length < 3) {
				this.errormsg = "Il nome utente deve avere almeno 3 caratteri.";
				this.loading = false;
				return;
			}

			try {
				let userRes = await this.$axios.get(`/users?name=${targetUserName}`);
				if (!userRes.data || userRes.data.length === 0) {
					this.errormsg = "Utente non trovato.";
					this.loading = false;
					return;
				}
				
				const targetUserId = userRes.data[0].userId;
				const payload = {targetUserId: targetUserId};
				let response = await this.$axios.post("/chats", payload);
				
				if (!this.chats) this.chats = [];
				
				const chatExists = this.chats.find(c => c.chatId === response.data.chatId);
				if (!chatExists) {
					this.chats.unshift(response.data);
				}
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},
		async deleteChat(chatId) {
			this.loading = true;
			this.errormsg = null;
			
			if (!chatId || chatId.length !== 20) {
				this.errormsg = "L'ID chat deve essere lungo esattamente 20 caratteri.";
				this.loading = false;
				return;
			}

			try {
				await this.$axios.delete(`/chats/${chatId}`);
				
				this.chats = this.chats.filter(c => c.chatId !== chatId);
				
				if (this.selectedChatId === chatId) {
					this.selectedChatId = null;
					this.messages = null;
				}
			} catch (e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
			}
		},
		async selectChat(chatId, silent = false) {
			if (!silent) {
				this.loading = true;
				this.errormsg = null;
			}
			try {
				if (this.selectedChatId !== chatId) {
					this.messages = null;
					this.selectedChatId = chatId;
					this.replyingToMessageId = null;
				}
				let response = await this.$axios.get(`/chats/${chatId}?limit=100`);
				if (!response.data.messages) {
					this.messages = [];
				} else {
					this.messages = response.data.messages; 
					this.messages.forEach(m => this.checkMessageStatus(m));
					this.updateMessageStatus();
				}
				this.selectedChatId = chatId;
				this.selectedChatTitle = response.data.title || '';
				this.selectedGroupId = response.data.groupId || null;
			} catch (e) {
				if (!silent) this.errormsg = e.toString();
			}
			if (!silent) this.loading = false;
		},
		async sendMessage() {
			if (!this.selectedChatId || !this.newMessage.trim()) return;

			this.errormsg = null;
			try {
				let payload = {
					chatId: this.selectedChatId,
					text: this.newMessage,
				};
				if (this.replyingToMessageId) {
					payload.replyTo = this.replyingToMessageId;
				}

                let response = await this.$axios.post(`/messages`, payload);
                    
                this.newMessage = '';
				this.replyingToMessageId = null;
                if (!this.messages) this.messages = [];
                this.messages.push(response.data);
            } catch (e) {
                this.errormsg = e.toString();
            }
        },
		async sendImageMessage(file) {
            if (!this.selectedChatId || !file) return;
            try {
				let formData = new FormData();
				formData.append('chatId', this.selectedChatId);
				formData.append('image', file);
				
				if (this.replyingToMessageId) {
					formData.append('replyTo', this.replyingToMessageId);
				}

                let response = await this.$axios.post(`/messages`, formData, {
					headers: {
						'Content-Type': 'multipart/form-data'
					}
				});
                    
				this.replyingToMessageId = null;
                if (!this.messages) this.messages = [];
                this.messages.push(response.data);
            } catch (e) {
                this.errormsg = e.toString();
            }
		},
		forwardMessage(messageId) {
			this.forwardingMessageId = messageId;
		},
		async confirmForward(targetChatId) {
			if (!targetChatId || !this.forwardingMessageId) return;

			this.loading = true;
			this.errormsg = null;
			try {
				const payload = {
					targetChatId: targetChatId 
				};
				let response = await this.$axios.post(`/messages/${this.forwardingMessageId}/forward`, payload);
				
				if (response.data) {
					if (targetChatId === this.selectedChatId) {
						if (!this.messages) this.messages = [];
						this.messages.push(response.data);
					}
				}
			} catch(e) {
				this.errormsg = e.toString();
			} finally {
				this.loading = false;
				this.forwardingMessageId = null;
			}
		},
		async deleteMessage(messageId) {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.delete(`/messages/${messageId}`);

				this.messages = this.messages.filter(m => m.messageId !== messageId);
			}
			catch(e) {
				this.errmsg = e.toString();
			}
			finally {
				this.loading = false;
			}
		},
		async updateMessage(messageId, newText) {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.put(`/messages/${messageId}`, { text: newText });
			}
			catch(e) {
				this.errormsg = e.toString();
			}
			finally {
				this.loading = false;
			}
		},
		async addReaction(messageId, emoji) {
			try {
				let response = await this.$axios.post(`/messages/${messageId}/comments`, { emoji: emoji });
				let msg = this.messages.find(m => m.messageId === messageId);
				if (msg) {
					if (!msg.comments) msg.comments = [];
					msg.comments.push(response.data);
				}
			} catch(e) {
				this.errormsg = e.toString();
			}
		},
		async removeReaction(messageId, commentId) {
			try {
				await this.$axios.delete(`/messages/${messageId}/comments/${commentId}`);
				let msg = this.messages.find(m => m.messageId === messageId);
				if (msg && msg.comments) {
					msg.comments = msg.comments.filter(c => c.commentId !== commentId);
				}
			} catch(e) {
				this.errormsg = e.toString();
			}
		},
		async createGroup() {
			if (!this.newGroupData.name || !this.newGroupData.name.trim()) {
				this.errormsg = "Il nome del gruppo è obbligatorio";
				return;
			}
			this.loading = true;
			this.errormsg = null;
			try {
				const body = { 
					groupName: this.newGroupData.name, 
					groupDescription: this.newGroupData.description, 
					photoUrl: this.newGroupData.photoUrl,
					members: this.newGroupData.selectedMembers
				};
				let response = await this.$axios.post(`/groups`, body);
				this.isCreatingGroup = false;
				this.newGroupData = { name: '', description: '', photoUrl: '', selectedMembers: [] };
				this.refresh();
			}
			catch(e) {
				this.errormsg = e.toString();
			}
			finally {
				this.loading = false;
			}
		},
		async leaveGroup(chatId) {
			this.loading = true;
			this.errormsg = null;
			try {
				const myUserId = localStorage.getItem("token");
				let response = await this.$axios.delete(`/groups/${chatId}/members/${myUserId}`);
				if (this.selectedChatId === chatId) {
					this.selectedChatId = null;
					this.messages = null;
				}
				this.refresh();
			}
			catch(e) {
				this.errormsg = e.toString();
			}
			finally {
				this.loading = false;
			}
		},
		async updateGroupPhoto(file) {
			if (!this.selectedGroupId || !file) return;
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.put(`/groups/${this.selectedGroupId}/photo`, file, {
					headers: {
						'Content-Type': file.type
					}
				});
				this.refresh();
			}
			catch(e) {
				this.errormsg = e.toString();
			}
			finally {
				this.loading = false;
			}
		},
		async confirmChangeGroupName() {
			if (this.newGroupNameInput && this.newGroupNameInput.trim() !== "") {
				this.isChangingGroupName = false;
				await this.updateGroupName(this.newGroupNameInput.trim());
			} else {
				this.errormsg = "Il nome del gruppo non può essere vuoto.";
			}
		},
		async updateGroupName(newName) {
			// Apre il modale se chiamato dal click del bottone (newName non è una stringa)
			if (typeof newName !== 'string') {
				this.newGroupNameInput = this.selectedChatTitle;
				this.isChangingGroupName = true;
				return;
			}
			if (!this.selectedGroupId) return;
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.put(`/groups/${this.selectedGroupId}/name`, { name: newName });
				this.selectedChatTitle = newName;
				this.refresh();
			}
			catch(e) {
				if (e.response && e.response.status === 403) {
					this.errormsg = "Solo l'admin può cambiare nome al gruppo!";
				} else {
					this.errormsg = e.toString();
				}
			}
			finally {
				this.loading = false;
			}
		},
		async addUserToGroup(userId) {
			if (!this.selectedGroupId) return;
			this.loading = true;
			this.errormsg = null;
			try {
				await this.$axios.post(`/groups/${this.selectedGroupId}/members`, { userId: userId });
				this.isAddingUserToGroup = false;
				this.refresh();
				alert("Utente aggiunto al gruppo con successo");
			}
			catch(e) {
				if (e.response && e.response.status === 409) {
					this.errormsg = "Questo utente fa già parte del gruppo!";
				} else if (e.response && e.response.status === 403) {
					this.errormsg = "Solo l'admin può aggiungere membri al gruppo";
				} else if (e.response && e.response.data) {
					this.errormsg = e.response.data;
				} else {
					this.errormsg = e.toString();
				}
			}
			finally {
				this.loading = false;
			}
		},
		async searchUsersForGroup() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(`/users?name=${this.groupSearchQuery}`);
				const myUserId = localStorage.getItem("token");
				if (response.data) {
					this.groupSearchResults = response.data.filter(u => u.userId !== myUserId);
				} else {
					this.groupSearchResults = [];
				}
			}
			catch(e) {
				this.errormsg = e.toString();
			}
			finally {
				this.loading = false;
			}
		},
		selectMembers(userId) {
			const index = this.newGroupData.selectedMembers.indexOf(userId);
			if (index === -1) {
				this.newGroupData.selectedMembers.push(userId);
			} else {
				this.newGroupData.selectedMembers.splice(index, 1);
			}
		},
		checkMessageStatus(message) {
			const myUserId = localStorage.getItem("token");
			if (message.state === "sent" && myUserId !== message.senderId) {
				if (!this.sentMessages.includes(message.messageId)) {
					this.sentMessages.push(message.messageId);
				}
			}
		},
		async updateMessageStatus() {
			if (!this.sentMessages || this.sentMessages.length === 0) return;
			
			this.errormsg = null;
			try {
				const unreadIds = [...this.sentMessages];
				this.sentMessages = [];
				
				for (let mid of unreadIds) {
					await this.$axios.put(`/messages/${mid}/status`, { status: "read" });
				}
				this.refresh();
			}
			catch(e) {
				this.errormsg = e.toString();
			}
		}
	},
	mounted() {
		this.refresh();
		this.pollIntervalId = setInterval(() => {
			this.refresh(true);
			if (this.selectedChatId) {
				this.selectChat(this.selectedChatId, true);
			}
		}, 3000);
	},
	unmounted() {
		if (this.pollIntervalId) {
			clearInterval(this.pollIntervalId);
		}
	}
}
</script>

<template>
	<LoadingSpinner :loading="loading">
		<div class="chat-layout">
			<div class="sidebar-container" v-if="!selectedChatId">
				<ChatsList 
					:chats="chats" 
					@select-chat="selectChat"
					@delete-chat="deleteChat"
					@create-chat="CreateChat"
					@open-create-group="isCreatingGroup = true"
				/>
			</div>
			<div class="main-container" v-if="selectedChatId">
				<ChatWindow 
					:messages="messages" 
					:selectedChatId="selectedChatId" 
					:chatTitle="selectedChatTitle" 
					:groupId="selectedGroupId"
					:availableEmojis="availableEmojis"
					v-model:draftMessage="newMessage" 
					v-model:replyingToMessageId="replyingToMessageId"
					@send-message="sendMessage" 
					@delete-message="deleteMessage" 
					@forward-message="forwardMessage" 
					@change-group-name="updateGroupName"
					@change-group-photo="updateGroupPhoto"
					@add-user-to-group="isAddingUserToGroup = true"
					@send-image-message="sendImageMessage"
					@add-reaction="addReaction"
					@remove-reaction="removeReaction"
					@reply-message="replyingToMessageId = $event"
					@close-chat="selectedChatId = null" 
				/>
			</div>
		</div>
	</LoadingSpinner>

	<ErrorMsg v-if="errormsg" :msg="errormsg" class="global-error"></ErrorMsg>
		
		<div v-if="forwardingMessageId" class="modal-overlay" @click.self="forwardingMessageId = null">
			<div class="modal-content">
				<div class="modal-header">
					<h3>Inoltra Messaggio</h3>
					<button class="close-btn" @click="forwardingMessageId = null">
						<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
					</button>
				</div>
				<div class="modal-body">
					<p v-if="!chats || chats.length === 0" class="no-chats">Non hai nessuna chat.</p>
					<div v-else class="simple-chats-list">
						<div v-for="chat in chats" :key="chat.chatId" class="simple-chat-item" @click="confirmForward(chat.chatId)">
							<svg class="feather" style="width: 14px; height: 14px; margin-right: 8px;"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
							{{ chat.title || chat.chatId }}
						</div>
					</div>
				</div>
			</div>
		</div>

		<!-- Modal Crea Gruppo -->
		<div v-if="isCreatingGroup" class="modal-overlay" @click.self="isCreatingGroup = false">
			<div class="modal-content">
				<div class="modal-header">
					<h3>Crea Nuovo Gruppo</h3>
					<button class="close-btn" @click="isCreatingGroup = false">
						<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
					</button>
				</div>
				<div class="modal-body group-creation-body">
					<div class="form-group">
						<input type="text" class="input-custom full-width" v-model="newGroupData.name" placeholder="Nome del gruppo (obbligatorio)">
					</div>
					<div class="form-group">
						<input type="text" class="input-custom full-width" v-model="newGroupData.description" placeholder="Descrizione del gruppo">
					</div>
					
					<hr>
					<h4>Aggiungi Membri (Opzionale)</h4>
					<div class="search-users-box">
						<input
							type="text"
							class="input-custom full-width search-input"
							v-model="groupSearchQuery"
							placeholder="Cerca utenti per nome..."
							@keyup="searchUsersForGroup"
						/>
					</div>
					
					<div class="users-list-scrollable" v-if="groupSearchResults.length > 0">
						<div v-for="user in groupSearchResults" :key="user.userId" class="user-selection-item" @click="selectMembers(user.userId)">
							<input 
								type="checkbox" 
								:checked="newGroupData.selectedMembers.includes(user.userId)"
								@click.stop="selectMembers(user.userId)"
							>
							<div class="user-info">
								<span class="user-name">{{ user.name }}</span>
								<span class="user-id">#{{ user.userId.substring(0,6) }}</span>
							</div>
						</div>
					</div>
					<p v-else class="no-chats">Nessun utente trovato.</p>
					
					<div class="selected-count" v-if="newGroupData.selectedMembers.length > 0">
						{{ newGroupData.selectedMembers.length }} membri selezionati
					</div>
				</div>
				<div class="modal-footer">
					<button class="btn-primary-custom full-width" @click="createGroup" :disabled="loading">
						<span v-if="loading">Creazione...</span>
						<span v-else>Crea Gruppo</span>
					</button>
				</div>
			</div>
		</div>

		<!-- Aggiungi Utente a Gruppo -->
		<div v-if="isAddingUserToGroup" class="modal-overlay" @click.self="isAddingUserToGroup = false">
			<div class="modal-content">
				<div class="modal-header">
					<h3>Aggiungi Utente</h3>
					<button class="close-btn" @click="isAddingUserToGroup = false">
						<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
					</button>
				</div>
				<div class="modal-body group-creation-body">
					<div class="search-users-box">
						<input
							type="text"
							class="input-custom full-width search-input"
							v-model="groupSearchQuery"
							placeholder="Cerca per nome..."
							@keyup="searchUsersForGroup"
						/>
					</div>
					
					<div class="users-list-scrollable" v-if="groupSearchResults.length > 0">
						<div v-for="user in groupSearchResults" :key="user.userId" class="user-selection-item" @click="addUserToGroup(user.userId)">
							<div class="user-info">
								<span class="user-name">{{ user.name }}</span>
								<span class="user-id">#{{ user.userId.substring(0,6) }}</span>
							</div>
							<button class="btn-primary-custom" style="padding: 4px 8px; font-size: 0.8rem;">Aggiungi</button>
						</div>
					</div>
					<p v-else class="no-chats">Nessun utente trovato.</p>
				</div>
			</div>
		</div>

		<!-- Modal Cambia Nome Gruppo -->
		<div v-if="isChangingGroupName" class="modal-overlay" @click.self="isChangingGroupName = false">
			<div class="modal-content">
				<div class="modal-header">
					<h3>Cambia Nome Gruppo</h3>
					<button class="close-btn" @click="isChangingGroupName = false">
						<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
					</button>
				</div>
				<div class="modal-body group-creation-body">
					<div class="form-group">
						<input type="text" class="input-custom full-width" v-model="newGroupNameInput" placeholder="Nuovo nome del gruppo" @keyup.enter="confirmChangeGroupName">
					</div>
				</div>
				<div class="modal-footer">
					<button class="btn-primary-custom full-width" @click="confirmChangeGroupName" :disabled="loading">
						<span v-if="loading">Salvataggio...</span>
						<span v-else>Salva Nome</span>
					</button>
				</div>
			</div>
		</div>

</template>

<style scoped>
.chat-layout {
	display: flex;
	height: calc(100vh - 60px);
	overflow: hidden;
	width: 100%;
}
.sidebar-container {
	width: 100%;
	max-width: 600px;
	margin: 0 auto;
	flex-shrink: 0;
	background-color: var(--color-bg-app);
	display: flex;
	flex-direction: column;
	overflow-y: auto;
}

.main-container {
	flex-grow: 1;
	background-color: var(--color-bg-chat);
	display: flex;
	flex-direction: column;
	position: relative;
}
.global-error {
	position: fixed;
	bottom: 20px;
	left: 50%;
	transform: translateX(-50%);
	z-index: 9999;
	box-shadow: var(--shadow-md);
}

/* Stili Modale in-line */
.modal-overlay {
    position: fixed;
    top: 0; left: 0; width: 100vw; height: 100vh;
    background-color: rgba(0, 0, 0, 0.4);
    display: flex; align-items: center; justify-content: center;
    z-index: 9999;
}
.modal-content {
    background-color: var(--color-bg-white);
    border-radius: var(--radius-lg);
    width: 100%; max-width: 400px; max-height: 80vh;
    display: flex; flex-direction: column;
    box-shadow: var(--shadow-lg);
}
.modal-header {
    padding: 16px 20px;
    border-bottom: 1px solid var(--color-border);
    display: flex; justify-content: space-between; align-items: center;
}
.modal-header h3 { margin: 0; font-size: 1.1rem; font-weight: 600; }
.close-btn {
    background: none; border: none; cursor: pointer; color: var(--color-text-muted); padding: 4px; border-radius: 50%;
}
.close-btn:hover { background-color: var(--color-border); color: var(--color-text-main); }
.modal-body { padding: 8px; overflow-y: auto; }
.no-chats { text-align: center; color: var(--color-text-muted); padding: 20px; }

/* Lista chat semplificata */
.simple-chats-list { display: flex; flex-direction: column; gap: 4px; }
.simple-chat-item {
    display: flex; align-items: center; padding: 12px; border-radius: var(--radius-md); 
    cursor: pointer; font-weight: 500; transition: background-color 0.2s; color: var(--color-text-main);
}
.simple-chat-item:hover { background-color: #F5F3ED; color: var(--color-primary); }

/* Stili Modale Gruppo */
.group-creation-body { padding: 16px 20px; }
.form-group { margin-bottom: 12px; }
.full-width { width: 100%; box-sizing: border-box; }
.search-users-box { margin-bottom: 12px; }
.users-list-scrollable { max-height: 200px; overflow-y: auto; border: 1px solid var(--color-border); border-radius: var(--radius-md); }
.user-selection-item { 
	display: flex; align-items: center; padding: 10px 12px; 
	border-bottom: 1px solid var(--color-border); cursor: pointer; transition: background 0.2s;
}
.user-selection-item:last-child { border-bottom: none; }
.user-selection-item:hover { background-color: var(--color-bg-chat); }
.user-info { display: flex; flex-direction: column; margin-left: 10px; }
.user-name { font-weight: 500; font-size: 0.95rem; }
.user-id { font-size: 0.8rem; color: var(--color-text-muted); }
.selected-count { text-align: center; font-size: 0.9rem; color: var(--color-primary); margin-top: 10px; font-weight: 500; }
.modal-footer { padding: 12px 20px; border-top: 1px solid var(--color-border); }
</style>

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
			selectedChatId: null,
			selectedChatTitle: '',
			newMessage: '',
			forwardingMessageId: null,
			availableEmojis: ['👍', '❤️', '💔', '😴', '😂', '😮', '😢', '🔥', '🤐'],
			pollIntervalId: null
		}
	},
	methods: {
		async refresh(silent = false) {
			if (!silent) this.loading = true;
			this.errormsg = null;
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
			if (!silent) this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(`/chats/${chatId}`);
				if (!response.data.messages) {
					this.messages = [];
				} else {
					this.messages = response.data.messages; 
				}
				this.selectedChatId = chatId;
				this.selectedChatTitle = response.data.title || '';
			} catch (e) {
				if (!silent) this.errormsg = e.toString();
			}
			if (!silent) this.loading = false;
		},
		async sendMessage() {
            if (!this.selectedChatId || !this.newMessage.trim()) return;
            try {
                const payload = {
                    chatId: this.selectedChatId,
                    text: this.newMessage
                };

                let response = await this.$axios.post(`/messages`, payload);
                    
                this.newMessage = '';
                if (!this.messages) this.messages = [];
                this.messages.push(response.data);
            } catch (e) {
                this.errmsg = e.toString();
            }
			finally {
				this.loading = false;
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
	},
	mounted() {
		this.refresh();
		this.pollIntervalId = setInterval(() => {
			this.refresh(true);
			if (this.selectedChatId) {
				this.selectChat(this.selectedChatId, true);
			}
		}, 3000); // 3 secondi di intervallo
	},
	unmounted() {
		if (this.pollIntervalId) {
			clearInterval(this.pollIntervalId);
		}
	}
}
</script>

<template>
	<div class="chat-layout">
		<div class="sidebar-container" v-if="!selectedChatId">
			<ChatsList :chats="chats" @select-chat="selectChat" @delete-chat="deleteChat" @create-chat="CreateChat" />
		</div>
		<div class="main-container" v-if="selectedChatId">
			<ChatWindow :messages="messages" :selectedChatId="selectedChatId" :chatTitle="selectedChatTitle" v-model:draftMessage="newMessage" @send-message="sendMessage" @delete-message="deleteMessage" @forward-message="forwardMessage" @close-chat="selectedChatId = null" />
		</div>
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
</style>

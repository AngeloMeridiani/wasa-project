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
			newMessage: '',
			availableEmojis: ['👍', '❤️', '💔', '😴', '😂', '😮', '😢', '🔥', '🤐']
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/chats");
				this.chats = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async CreateChat(targetUserId) {
			this.loading = true;
			this.errormsg = null;
			if (!targetUserId || targetUserId.length !== 20) {
				this.errormsg = "L'ID utente deve essere lungo esattamente 20 caratteri.";
				this.loading = false;
				return;
			}
			const isAlphanumeric = /^[a-zA-Z0-9]+$/.test(targetUserId);
			if (!isAlphanumeric) {
				this.errormsg = "L'ID utente può contenere solo lettere e numeri.";
				this.loading = false;
				return;
			}

			try {
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
		async selectChat() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/messages");
				this.messages = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
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
        },
	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div class="row">
		
		<div class="col-md-4 border-end">
			<ChatsList />
		</div>
		
		<div class="col-md-8">
			<ChatWindow />
		</div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>

<script>
export default {
    name: "ChatsList",
    props: {
        chats: {
            type: Array,
            default: () => []
        }
    },
    data() {
        return {
            targetUserName: "",
            myUserId: localStorage.getItem("token")
        }
    }
}
</script>

<template>
    <div class="chats-list-container">
        <div class="chats-header">
            <h3 class="chats-title">Chat</h3>
        </div>
        
        <div class="create-chat-box">
            <input 
                type="text" 
                class="input-custom search-input" 
                v-model="targetUserName" 
                placeholder="Nome utente" 
                @keyup.enter="$emit('create-chat', targetUserName)"
            />
            <button class="btn-primary-custom create-btn" @click="$emit('create-chat', targetUserName)" title="Crea Nuova Chat Singola">
                <svg class="feather" style="width: 16px; height: 16px;"><use href="/feather-sprite-v4.29.0.svg#user-plus"/></svg>
            </button>
            <button class="btn-primary-custom create-btn group-btn" @click="$emit('open-create-group')" title="Crea Nuovo Gruppo">
                <svg class="feather" style="width: 16px; height: 16px;"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
            </button>
        </div>

        <div class="chats-scroll-area">
            <ul v-if="chats && chats.length" class="chats-ul">
                <li 
                    v-for="chat in chats" 
                    :key="chat.chatId" 
                    class="chat-card"
                    @click="$emit('select-chat', chat.chatId)"
                >
                    <div class="chat-card-content">
                        <div class="chat-avatar">
                            {{ (chat.title || chat.chatId).substring(0, 1).toUpperCase() }}
                        </div>
                        <div class="chat-info">
                            <h4 class="chat-name">{{ chat.title || chat.chatId }}</h4>
                            <p v-if="chat.lastMessage && chat.lastMessage.contentType === 'image'" class="chat-last-msg italic"> Immagine </p>
                            <p v-else-if="chat.lastMessage" class="chat-last-msg">{{ chat.lastMessage.text }}</p>
                            <p v-else class="chat-last-msg italic">Nessun messaggio</p>
                        </div>
                    </div>
                    <!-- Mostra il badge solo se ci sono messaggi non letti -->
                    <div v-if="chat.unreadCount > 0" class="unread-badge">
                        {{ chat.unreadCount }}
                    </div>
                    <button class="delete-btn" @click.stop="$emit('delete-chat', chat.chatId)" title="Elimina chat">
                        <svg class="feather" style="width: 16px; height: 16px;"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
                    </button>
                </li>
            </ul>
            
            <div v-else class="empty-state">
                <svg class="feather empty-icon"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
                <p>Nessuna chat disponibile.</p>
                <small>Inserisci un nome utente per avviare una nuova chat con lui.</small>
            </div>
        </div>
    </div>
</template>

<style scoped>
.chats-list-container {
    display: flex;
    flex-direction: column;
    height: 100%;
}
.chats-header {
    padding: 20px 20px 10px 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.chats-title {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--color-text-main);
}
.user-id-badge {
    font-size: 0.75rem;
    background-color: var(--color-bg-app);
    padding: 6px 10px;
    border-radius: var(--radius-sm);
    color: var(--color-primary);
    font-weight: 500;
    display: flex;
    align-items: center;
    cursor: help;
    border: 1px solid var(--color-border);
}

.unread-badge{
    font-size: 0.75rem;
    background-color: var(--color-bg-app);
    padding: 6px 10px;
    border-radius: var(--radius-sm);
    color: var(--color-primary);
    font-weight: 500;
    display: flex;
    align-items: center;
    cursor: help;
    border: 1px solid var(--color-border);
}

.create-chat-box {
    padding: 10px 20px 20px 20px;
    display: flex;
    gap: 10px;
    border-bottom: 1px solid var(--color-border);
}
.search-input {
    flex-grow: 1;
    font-size: 0.9rem;
    padding: 8px 14px;
}
.create-btn {
    padding: 8px 12px;
    display: flex;
    align-items: center;
    justify-content: center;
}
.chats-scroll-area {
    flex-grow: 1;
    overflow-y: auto;
    padding: 10px;
}
.chats-ul {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 5px;
}
.chat-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px;
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: background-color 0.2s ease, transform 0.1s ease;
}
.chat-card:hover {
    background-color: #F5F3ED;
}
.chat-card:active {
    transform: scale(0.99);
}
.chat-card-content {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-grow: 1;
    overflow: hidden;
}
.chat-avatar {
    width: 44px;
    height: 44px;
    background-color: var(--color-primary);
    color: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 1.2rem;
    flex-shrink: 0;
}
.chat-info {
    overflow: hidden;
    display: flex;
    flex-direction: column;
    gap: 4px;
}
.chat-name {
    margin: 0;
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text-main);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.chat-last-msg {
    margin: 0;
    font-size: 0.85rem;
    color: var(--color-text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
.italic {
    font-style: italic;
    opacity: 0.7;
}
.delete-btn {
    background: none;
    border: none;
    color: #FF5252;
    padding: 8px;
    border-radius: 50%;
    cursor: pointer;
    opacity: 0;
    transition: opacity 0.2s ease, background-color 0.2s ease;
}
.chat-card:hover .delete-btn {
    opacity: 0.7;
}
.delete-btn:hover {
    opacity: 1 !important;
    background-color: rgba(255, 82, 82, 0.1);
}
.empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--color-text-muted);
    text-align: center;
    padding: 20px;
}
.empty-icon {
    width: 48px;
    height: 48px;
    margin-bottom: 16px;
    color: var(--color-primary);
    opacity: 0.8;
}
</style>

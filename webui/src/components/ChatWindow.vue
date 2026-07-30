<script>
export default {
    name: "ChatWindow",
    props: {
        messages: {
            type: Array,
            default: () => []
        },
        selectedChatId: {
            type: String,
            default: null
        },
        draftMessage: {
            type: String,
            default: ""
        },
        chatTitle: {
            type: String,
            default: ""
        }
    },
    data() {
        return {
            myUserId: localStorage.getItem("token")
        }
    },
    computed: {
        localMessage: {
            get() {
                return this.draftMessage;
            },
            set(value) {
                this.$emit('update:draftMessage', value);
            }
        },
        displayTitle() {
            return this.chatTitle || this.selectedChatId || 'Chat Sconosciuta';
        }
    },
    methods: {
        send() {
            if (this.localMessage.trim()) {
                this.$emit('send-message');
            }
        }
    }
}
</script>

<template>
    <div class="chat-window-container">
        
        <div v-if="!selectedChatId" class="no-chat-selected">
            <div class="welcome-box">
                <svg class="feather welcome-icon"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
                <h2>WasaText per il Web</h2>
                <p>Seleziona una chat dalla lista per iniziare a inviare messaggi.</p>
            </div>
        </div>
        
        <div v-else class="chat-active">
            <!-- Header della Chat -->
            <div class="chat-header">
                <div class="chat-header-info">
                    <div class="chat-avatar-small">
                        {{ displayTitle.substring(0, 1).toUpperCase() }}
                    </div>
                    <div class="chat-title-group">
                        <h3 class="chat-title">{{ displayTitle }}</h3>
                        <small class="chat-subtitle">ID: {{ selectedChatId.substring(0,8) }}...</small>
                    </div>
                </div>
                <div class="chat-header-actions">
                    <!-- Eventuali bottoni per chiamate, info, ecc. -->
                </div>
            </div>
            
            <!-- Area Messaggi -->
            <div class="messages-area">
                <div v-if="messages && messages.length" class="messages-list">
                    <div 
                        v-for="msg in messages" 
                        :key="msg.messageId" 
                        :class="['message-wrapper', msg.senderId === myUserId ? 'message-out' : 'message-in']"
                    >
                        <div class="message-bubble">
                            <!-- Nome mittente visibile solo nei gruppi o se non è l'utente corrente, ma per minimalismo lo nascondiamo se è il mio -->
                            <div v-if="msg.senderId !== myUserId" class="message-sender">
                                {{ msg.senderId.substring(0,6) }}
                            </div>
                            <div class="message-text">
                                {{ msg.text }}
                            </div>
                            <div class="message-meta">
                                <span class="message-time">{{ msg.createdAt ? msg.createdAt.substring(11, 16) : '' }}</span>
                                
                                <div class="message-actions">
                                    <button @click="$emit('forward-message', msg.messageId)" class="msg-action-btn" title="Inoltra">
                                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#corner-up-right"/></svg>
                                    </button>
                                    <button @click="$emit('delete-message', msg.messageId)" class="msg-action-btn delete-icon" title="Elimina">
                                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div v-else class="empty-chat">
                    Nessun messaggio in questa chat. Scrivi qualcosa!
                </div>
            </div>

            <!-- Area Input -->
            <div class="input-area">
                <div class="input-container">
                    <button class="attach-btn" title="Allega">
                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#paperclip"/></svg>
                    </button>
                    <input 
                        type="text" 
                        class="chat-input" 
                        v-model="localMessage" 
                        placeholder="Scrivi un messaggio..." 
                        @keyup.enter="send" 
                    />
                    <button class="send-btn" @click="send" :disabled="!localMessage.trim()">
                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#send"/></svg>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.chat-window-container {
    height: 100%;
    width: 100%;
    display: flex;
    flex-direction: column;
}

/* Stato vuoto (Nessuna chat selezionata) */
.no-chat-selected {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: var(--color-bg-app);
}
.welcome-box {
    text-align: center;
    max-width: 400px;
    padding: 40px;
    background-color: var(--color-bg-white);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
}
.welcome-icon {
    width: 64px;
    height: 64px;
    color: var(--color-primary);
    margin-bottom: 20px;
}
.welcome-box h2 {
    font-weight: 300;
    margin-bottom: 10px;
}
.welcome-box p {
    color: var(--color-text-muted);
}

/* Chat attiva */
.chat-active {
    display: flex;
    flex-direction: column;
    height: 100%;
}

/* Header */
.chat-header {
    height: 70px;
    background-color: var(--color-bg-white);
    border-bottom: 1px solid var(--color-border);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 24px;
    flex-shrink: 0;
}
.chat-header-info {
    display: flex;
    align-items: center;
    gap: 16px;
}
.chat-avatar-small {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: var(--color-primary);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
}
.chat-title-group {
    display: flex;
    flex-direction: column;
}
.chat-title {
    margin: 0;
    font-size: 1.1rem;
    font-weight: 600;
}
.chat-subtitle {
    color: var(--color-text-muted);
    font-size: 0.8rem;
}

/* Area Messaggi */
.messages-area {
    flex-grow: 1;
    overflow-y: auto;
    padding: 24px;
    display: flex;
    flex-direction: column;
    /* Un pattern di sfondo tenue opzionale per la chat */
    background-color: var(--color-bg-chat);
}
.messages-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}
.empty-chat {
    text-align: center;
    color: var(--color-text-muted);
    margin-top: auto;
    margin-bottom: auto;
    font-style: italic;
}

/* Bolle Messaggi */
.message-wrapper {
    display: flex;
    width: 100%;
}
.message-out {
    justify-content: flex-end;
}
.message-in {
    justify-content: flex-start;
}

.message-bubble {
    max-width: 65%;
    padding: 10px 14px;
    position: relative;
    box-shadow: var(--shadow-msg);
    display: flex;
    flex-direction: column;
}

.message-out .message-bubble {
    background-color: var(--color-message-out);
    color: white;
    border-radius: var(--radius-md) var(--radius-md) 0 var(--radius-md);
}

.message-in .message-bubble {
    background-color: var(--color-message-in);
    color: var(--color-text-main);
    border-radius: var(--radius-md) var(--radius-md) var(--radius-md) 0;
}

.message-sender {
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--color-primary);
    margin-bottom: 4px;
}

.message-text {
    font-size: 0.95rem;
    line-height: 1.4;
    word-wrap: break-word;
}

.message-meta {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 10px;
    margin-top: 4px;
    font-size: 0.7rem;
}

.message-out .message-meta {
    color: rgba(255, 255, 255, 0.8);
}

.message-in .message-meta {
    color: var(--color-text-muted);
}

.message-actions {
    display: flex;
    gap: 8px;
    opacity: 0;
    transition: opacity 0.2s;
}

.message-bubble:hover .message-actions {
    opacity: 1;
}

.msg-action-btn {
    background: none;
    border: none;
    padding: 0;
    cursor: pointer;
    color: inherit;
    opacity: 0.7;
}

.msg-action-btn:hover {
    opacity: 1;
}
.delete-icon:hover {
    color: #FF5252;
}
.msg-action-btn svg {
    width: 14px;
    height: 14px;
}

/* Area Input */
.input-area {
    padding: 16px 24px;
    background-color: var(--color-bg-app);
    border-top: 1px solid var(--color-border);
    flex-shrink: 0;
}
.input-container {
    display: flex;
    align-items: center;
    gap: 12px;
    background-color: var(--color-bg-white);
    border-radius: var(--radius-xl);
    padding: 8px 16px;
    box-shadow: var(--shadow-sm);
}
.attach-btn {
    background: none;
    border: none;
    color: var(--color-text-muted);
    cursor: pointer;
    display: flex;
    align-items: center;
    padding: 4px;
}
.attach-btn:hover {
    color: var(--color-primary);
}
.chat-input {
    flex-grow: 1;
    border: none;
    outline: none;
    font-size: 1rem;
    background: transparent;
    padding: 8px 0;
}
.send-btn {
    background-color: var(--color-primary);
    color: white;
    border: none;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: background-color 0.2s, transform 0.1s;
}
.send-btn:hover:not(:disabled) {
    background-color: var(--color-primary-hover);
    transform: scale(1.05);
}
.send-btn:disabled {
    background-color: var(--color-border);
    cursor: not-allowed;
}
.send-btn svg {
    width: 18px;
    height: 18px;
    margin-left: -2px; /* Centra otticamente l'icona paper-plane */
}
</style>

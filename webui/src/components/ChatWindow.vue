<script>
import MessageStatus from './MessageStatus.vue';

export default {
    name: "ChatWindow",
    components: {
        MessageStatus
    },
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
        },
        groupId: {
            type: String,
            default: null
        },
        availableEmojis: {
            type: Array,
            default: () => []
        },
        replyingToMessageId: {
            type: String,
            default: null
        }
    },
    data() {
        return {
            myUserId: localStorage.getItem("token"),
            activeReactionMessageId: null
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
            return this.chatTitle || 'Chat Senza Nome';
        }
    },
    methods: {
        send() {
            if (this.localMessage.trim()) {
                this.$emit('send-message');
                this.scrollToBottom();
            }
        },
        triggerGroupPhotoUpload() {
            this.$refs.groupPhotoInput.click();
        },
        onGroupPhotoSelected(event) {
            const file = event.target.files[0];
            if (file) {
                this.$emit('change-group-photo', file);
            }
            if (this.$refs.groupPhotoInput) {
                this.$refs.groupPhotoInput.value = "";
            }
        },
        triggerImageMessageUpload() {
            this.$refs.imageMessageInput.click();
        },
        onImageMessageSelected(event) {
            const file = event.target.files[0];
            if (file) {
                this.$emit('send-image-message', file);
            }
            if (this.$refs.imageMessageInput) {
                this.$refs.imageMessageInput.value = "";
            }
        },
        formatMessageDate(dateStr) {
            if (!dateStr) return '';
            // Formato atteso: 2026-08-18T13:58:27Z
            // Output desiderato: 18/08 13:58
            try {
                const dateObj = new Date(dateStr);
                const day = String(dateObj.getDate()).padStart(2, '0');
                const month = String(dateObj.getMonth() + 1).padStart(2, '0');
                const hours = String(dateObj.getHours()).padStart(2, '0');
                const minutes = String(dateObj.getMinutes()).padStart(2, '0');
                
                // Controlliamo se è oggi per eventualmente omettere la data
                const today = new Date();
                if (dateObj.getDate() === today.getDate() && 
                    dateObj.getMonth() === today.getMonth() && 
                    dateObj.getFullYear() === today.getFullYear()) {
                    return `${hours}:${minutes}`;
                }
                
                
                return `${day}/${month} ${hours}:${minutes}`;
            } catch (e) {
                return dateStr.substring(11, 16); // fallback
            }
        },
        groupedReactions(comments) {
            if (!comments || comments.length === 0) return [];
            let groups = {};
            comments.forEach(c => {
                if (!groups[c.emoji]) {
                    groups[c.emoji] = { emoji: c.emoji, count: 0, me: false, myCommentId: null, authors: [] };
                }
                groups[c.emoji].count++;
                groups[c.emoji].authors.push(c.senderName || c.senderId);
                if (c.senderId === this.myUserId) {
                    groups[c.emoji].me = true;
                    groups[c.emoji].myCommentId = c.commentId;
                }
            });
            return Object.values(groups).sort((a,b) => b.count - a.count);
        },
        handleReactionClick(messageId, group) {
            if (group.me) {
                this.$emit('remove-reaction', messageId, group.myCommentId);
            } else {
                this.$emit('add-reaction', messageId, group.emoji);
            }
            this.activeReactionMessageId = null;
        },
        selectEmoji(messageId, emoji) {
            this.$emit('add-reaction', messageId, emoji);
            this.activeReactionMessageId = null;
        },
        getMessageTextById(messageId) {
            const m = this.getMessageById(messageId);
            if (!m) return 'Messaggio originale non trovato';
            if (m.contentType === 'image') return '📷 Immagine';
            return m.text.length > 50 ? m.text.substring(0, 50) + '...' : m.text;
        },
        getMessageById(messageId) {
            if (!this.messages) return null;
            return this.messages.find(msg => msg.messageId === messageId);
        },
        cancelReply() {
            this.$emit('update:replyingToMessageId', null);
        },
        scrollToBottom() {
            this.$nextTick(() => {
                if (this.$refs.messagesArea) {
                    this.$refs.messagesArea.scrollTop = this.$refs.messagesArea.scrollHeight;
                }
            });
        }
    },
    watch: {
        messages: {
            deep: true,
            handler() {
                const el = this.$refs.messagesArea;
                if (!el) return;
                
                // Se non ci sono messaggi (chat vuota o in caricamento), scrolla giù
                if (!this.messages || this.messages.length === 0) {
                    this.scrollToBottom();
                    return;
                }

                // Verifica se l'utente si trova vicino in fondo (tolleranza 100px)
                const isAtBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 100;
                
                if (isAtBottom) {
                    this.scrollToBottom();
                }
            }
        }
    },
    mounted() {
        this.scrollToBottom();
    }
}
</script>

<template>
    <div class="chat-window-container">
            <!-- Header della Chat -->
            <div class="chat-header">
                <div class="chat-header-info">
                    <button class="back-btn" @click="$emit('close-chat')" title="Indietro">
                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#arrow-left"/></svg>
                    </button>
                    <div class="chat-avatar-small">
                        {{ displayTitle.substring(0, 1).toUpperCase() }}
                    </div>
                    <div class="chat-title-group">
                        <h3 class="chat-title">{{ displayTitle }}</h3>
                    </div>
                    <!-- Bottoni azioni gruppo, visibili solo se la chat è un gruppo (ha groupId) -->
                    <div v-if="groupId" class="group-actions ml-auto d-flex align-items-center">
                        <button class="action-btn" @click="$emit('add-user-to-group')" title="Aggiungi Utenti">
                            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user-plus"/></svg>
                        </button>
                        <button class="action-btn" @click="$emit('change-group-name')" title="Cambia Nome Gruppo">
                            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit-2"/></svg>
                        </button>
                        <button class="action-btn" @click="triggerGroupPhotoUpload" title="Cambia Foto Gruppo">
                            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#camera"/></svg>
                        </button>
                        <input type="file" ref="groupPhotoInput" style="display: none" accept="image/png, image/jpeg" @change="onGroupPhotoSelected">
                    </div>
                </div>
            </div>
            
            <!-- Area Messaggi -->
            <div class="messages-area" ref="messagesArea">
                <div v-if="!messages || messages.length === 0" class="empty-chat">
                    Nessun messaggio in questa chat.
                </div>
                
                <div v-else class="messages-list">
                    <div 
                        v-for="msg in messages" 
                        :key="msg.messageId" 
                        :class="['message-wrapper', msg.senderId === myUserId ? 'message-out' : 'message-in']"
                    >
                        <div class="message-bubble">
                            
                            <!-- Anteprima Risposta -->
                            <div v-if="msg.replyTo" class="reply-preview-bubble">
                                <div class="reply-preview-text">
                                    <small>Risposta a:</small>
                                    <img v-if="getMessageById(msg.replyTo)?.contentType === 'image'" :src="getMessageById(msg.replyTo).photoUrl" class="reply-preview-img-small" />
                                    <div v-else>{{ getMessageTextById(msg.replyTo) }}</div>
                                </div>
                            </div>

                            <div v-if="msg.contentType === 'image'" class="message-image-container">
                                <img :src="msg.photoUrl" alt="Immagine inviata" class="message-img" />
                            </div>
                            <div v-else class="message-text">
                                {{ msg.text }}
                            </div>
                            <div class="message-meta">
                                <span class="message-time">{{ formatMessageDate(msg.createdAt) }}</span>
                                
                                <MessageStatus v-if="msg.senderId === myUserId" :status="msg.state" />

                                <div class="message-actions">
                                    <button @click="$emit('reply-message', msg.messageId)" class="msg-action-btn" title="Rispondi">
                                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#corner-up-left"/></svg>
                                    </button>
                                    <button @click="activeReactionMessageId = activeReactionMessageId === msg.messageId ? null : msg.messageId" class="msg-action-btn" title="Reagisci">
                                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#smile"/></svg>
                                    </button>
                                    <button @click="$emit('forward-message', msg.messageId)" class="msg-action-btn" title="Inoltra">
                                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#corner-up-right"/></svg>
                                    </button>
                                    <button @click="$emit('delete-message', msg.messageId)" class="msg-action-btn delete-icon" title="Elimina">
                                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
                                    </button>
                                </div>
                            </div>
                            
                            <!-- Area Picker Emoji (visibile solo se attivata) -->
                            <div v-if="activeReactionMessageId === msg.messageId" class="emoji-picker">
                                <span v-for="emoji in availableEmojis" :key="emoji" @click="selectEmoji(msg.messageId, emoji)" class="emoji-option">
                                    {{ emoji }}
                                </span>
                            </div>

                            <!-- Area Reazioni Attive -->
                            <div v-if="msg.comments && msg.comments.length > 0" class="reactions-container">
                                <div 
                                    v-for="group in groupedReactions(msg.comments)" 
                                    :key="group.emoji" 
                                    @click="handleReactionClick(msg.messageId, group)"
                                    :class="['reaction-badge', { 'reaction-me': group.me }]"
                                    :title="group.authors.join(', ')"
                                >
                                    <span class="reaction-emoji">{{ group.emoji }}</span>
                                    <span class="reaction-count" v-if="group.count > 1">{{ group.count }}</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Area Input Messaggio -->
            <div class="chat-input-area">
                <div v-if="replyingToMessageId" class="replying-banner">
                    <div class="replying-text">
                        <small>Stai rispondendo a:</small>
                        <img v-if="getMessageById(replyingToMessageId)?.contentType === 'image'" :src="getMessageById(replyingToMessageId).photoUrl" class="reply-preview-img-small" />
                        <div v-else>{{ getMessageTextById(replyingToMessageId) }}</div>
                    </div>
                    <button class="cancel-reply-btn" @click="cancelReply" title="Annulla">
                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x"/></svg>
                    </button>
                </div>
                <div class="input-container">
                    <button class="attach-btn" title="Allega Immagine" @click="triggerImageMessageUpload">
                        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image"/></svg>
                    </button>
                    <input type="file" ref="imageMessageInput" style="display: none" accept="image/png, image/jpeg" @change="onImageMessageSelected">
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
</template>

<style scoped>
.chat-window-container {
    height: 100%;
    width: 100%;
    display: flex;
    flex-direction: column;
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
    gap: 12px;
}
.back-btn {
    background: none;
    border: none;
    color: var(--color-text-main);
    cursor: pointer;
    padding: 8px; /* Area cliccabile più ampia */
    margin-right: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
}
.back-btn svg {
    width: 28px; /* Ingrandisce l'icona */
    height: 28px;
}
.back-btn:hover {
    color: var(--color-primary);
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

/* Pulsanti Azione Gruppo */
.group-actions {
    margin-left: auto;
    display: flex;
    gap: 12px;
}
.action-btn {
    background: transparent;
    border: none;
    color: var(--color-text-muted);
    cursor: pointer;
    padding: 6px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s, color 0.2s;
}
.action-btn:hover {
    background-color: #F5F3ED;
    color: var(--color-primary);
}
.action-btn svg {
    width: 20px;
    height: 20px;
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
    margin-left: -2px;
}

.message-image-container {
    padding: 2px;
    border-radius: 8px;
    background: transparent;
}

.message-img {
    max-width: 100%;
    max-height: 250px;
    border-radius: 8px;
    object-fit: cover;
}

/* Reazioni ed Emoji */
.emoji-picker {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 8px;
    padding: 8px;
    background-color: #222222; /* Sfondo scuro per contrasto */
    color: white; /* Per i font emoji monocromatici (es. su Linux) */
    border-radius: var(--radius-md);
    border: 1px solid var(--color-border);
}

.emoji-option {
    font-size: 1.25rem;
    cursor: pointer;
    transition: transform 0.2s ease;
}

.emoji-option:hover {
    transform: scale(1.2);
}

.reactions-container {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    margin-top: 6px;
}

.reaction-badge {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    background-color: var(--color-bg-app);
    padding: 2px 6px;
    border-radius: 12px;
    font-size: 0.85rem;
    cursor: pointer;
    border: 1px solid var(--color-border);
    transition: background-color 0.2s ease;
}

.reaction-badge:hover {
    background-color: #F5F3ED;
}

.reaction-me {
    background-color: rgba(255, 126, 64, 0.1);
    border-color: var(--color-primary);
}

.reaction-emoji {
    font-size: 1rem;
}

/* Risposte (Reply) */
.reply-preview-bubble {
    background-color: rgba(0, 0, 0, 0.05);
    border-left: 3px solid var(--color-primary);
    padding: 6px 10px;
    margin-bottom: 6px;
    border-radius: 4px;
    font-size: 0.85rem;
    color: var(--color-text-muted);
}
.message-out .reply-preview-bubble {
    background-color: rgba(255, 255, 255, 0.2);
    border-left-color: white;
    color: white;
}
.reply-preview-text small {
    display: block;
    font-weight: 600;
    margin-bottom: 2px;
}

.replying-banner {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 16px;
    background-color: var(--color-bg-white);
    border-bottom: 1px solid var(--color-border);
    border-radius: var(--radius-xl) var(--radius-xl) 0 0;
    margin-bottom: -10px;
    z-index: 1;
}
.replying-text small {
    color: var(--color-primary);
    font-weight: 600;
}
.replying-text div {
    font-size: 0.9rem;
    color: var(--color-text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 300px;
}
.cancel-reply-btn {
    background: none;
    border: none;
    color: var(--color-text-muted);
    cursor: pointer;
    padding: 4px;
}
.cancel-reply-btn:hover {
    color: var(--color-error);
}

.reply-preview-img-small {
    max-height: 50px;
    max-width: 150px;
    border-radius: 4px;
    object-fit: cover;
}

.reaction-count {
    font-weight: 600;
    color: var(--color-text-main);
}
</style>

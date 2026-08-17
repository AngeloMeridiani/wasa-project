<script setup>
import { defineProps, computed } from 'vue';

const props = defineProps({
  status: {
    type: String,
    required: true,
    validator: (value) => ['sent', 'received', 'read'].includes(value)
  }
});

const statusText = computed(() => {
  switch (props.status) {
    case 'sent': return 'Inviato';
    case 'received': return 'Ricevuto';
    case 'read': return 'Letto';
    default: return '';
  }
});
</script>

<template>
  <div class="message-status" :title="statusText">
    <!-- Spunta singola: Inviato -->
    <svg v-if="status === 'sent'" viewBox="0 0 24 24" width="14" height="14">
      <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" fill="currentColor"/>
    </svg>

    <!-- Doppia spunta: Consegnato/Ricevuto o Letto -->
    <svg v-else-if="status === 'received' || status === 'read'" viewBox="0 0 24 24" width="14" height="14">
      <path d="M18 7l-1.41-1.41-6.34 6.34 1.41 1.41L18 7zm4.24-1.41L11.66 16.17 7.48 12l-1.41 1.41L11.66 19l12-12-1.42-1.41zM.41 13.41L6 19l1.41-1.41L1.83 12 .41 13.41z" fill="currentColor"/>
    </svg>
  </div>
</template>

<style scoped>
.message-status {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 2px;
}

.message-status svg {
  color: currentColor; /* Torna nero esattamente come prima */
}
</style>
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api, type Note } from './api'

const notes = ref<Note[]>([])
const newText = ref('')
const loading = ref(false)
const submitting = ref(false)
const error = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    notes.value = await api.list()
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}

async function addNote() {
  const text = newText.value.trim()
  if (!text || submitting.value) return
  submitting.value = true
  try {
    await api.create(text)
    newText.value = ''
    await load()
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    submitting.value = false
  }
}

async function removeNote(id: number) {
  try {
    await api.remove(id)
    await load()
  } catch (e) {
    error.value = (e as Error).message
  }
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString('ru-RU')
}

onMounted(load)
</script>

<template>
  <main>
    <header>
      <h1>📝 Заметки <span class="badge">{{ notes.length }}</span></h1>
      <p class="subtitle">demo: Vue + TypeScript + Go + PostgreSQL</p>
    </header>

    <form class="add-form" @submit.prevent="addNote">
      <input v-model="newText" type="text" placeholder="О чём хотите не забыть?" />
      <button type="submit" :disabled="!newText.trim() || submitting">
        {{ submitting ? 'Добавляем…' : 'Добавить' }}
      </button>
    </form>

    <p v-if="error" class="error">⚠️ {{ error }}</p>

    <div v-if="loading" class="skeleton">
      <div class="skeleton-row" v-for="i in 3" :key="i" />
    </div>

    <TransitionGroup v-else tag="ul" name="list" class="notes">
      <li v-for="note in notes" :key="note.id" class="note-card">
        <span class="text">{{ note.text }}</span>
        <span class="date">{{ formatDate(note.created_at) }}</span>
        <button class="delete-btn" @click="removeNote(note.id)" aria-label="Удалить заметку">✕</button>
      </li>
    </TransitionGroup>

    <div v-if="!loading && notes.length === 0" class="empty-state">
      <p class="empty-emoji">🗒️</p>
      <p>Заметок пока нет</p>
      <p class="hint">Добавьте первую заметку выше ↑</p>
    </div>
  </main>
</template>

<style scoped>
main {
  max-width: 640px;
  margin: 2rem auto;
  font-family: system-ui, sans-serif;
  padding: 0 1rem 3rem;
  color: #1a1a2e;
}

header {
  margin-bottom: 1.5rem;
}

h1 {
  display: flex;
  align-items: center;
  gap: 0.6rem;
  margin: 0;
  background: linear-gradient(135deg, #6366f1, #ec4899);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.badge {
  background: #6366f1;
  color: #fff;
  font-size: 0.9rem;
  font-weight: 600;
  padding: 0.1rem 0.6rem;
  border-radius: 999px;
  -webkit-text-fill-color: #fff;
}

.subtitle {
  margin: 0.25rem 0 0;
  color: #6b7280;
  font-size: 0.9rem;
}

.add-form {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

input {
  flex: 1;
  padding: 0.6rem 0.8rem;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 1rem;
}

input:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgb(99 102 241 / 0.15);
}

button {
  padding: 0.6rem 1.2rem;
  border: none;
  border-radius: 8px;
  background: #6366f1;
  color: #fff;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

button:hover:not(:disabled) {
  background: #4f46e5;
}

button:disabled {
  background: #c7c8f5;
  cursor: not-allowed;
}

.notes {
  list-style: none;
  padding: 0;
  margin: 0;
}

.note-card {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.9rem 1rem;
  margin-bottom: 0.6rem;
  background: #fff;
  border: 1px solid #eee;
  border-radius: 10px;
  box-shadow: 0 1px 2px rgb(0 0 0 / 0.04);
}

.text {
  flex: 1;
  word-break: break-word;
}

.date {
  color: #9ca3af;
  font-size: 0.8rem;
  white-space: nowrap;
}

.delete-btn {
  background: transparent;
  color: #9ca3af;
  padding: 0.2rem 0.5rem;
  font-weight: 400;
  border-radius: 6px;
}

.delete-btn:hover:not(:disabled) {
  background: #fee2e2;
  color: #dc2626;
}

.error {
  color: #dc2626;
  background: #fef2f2;
  border: 1px solid #fecaca;
  padding: 0.6rem 0.8rem;
  border-radius: 8px;
}

.empty-state {
  text-align: center;
  color: #9ca3af;
  padding: 2.5rem 1rem;
}

.empty-emoji {
  font-size: 2.5rem;
  margin: 0 0 0.5rem;
}

.hint {
  font-size: 0.85rem;
}

.skeleton {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.skeleton-row {
  height: 3rem;
  border-radius: 10px;
  background: linear-gradient(90deg, #f3f4f6 25%, #e5e7eb 50%, #f3f4f6 75%);
  background-size: 200% 100%;
  animation: shimmer 1.4s infinite;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

.list-enter-active,
.list-leave-active {
  transition: all 0.25s ease;
}

.list-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}

.list-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

.list-leave-active {
  position: absolute;
}
</style>

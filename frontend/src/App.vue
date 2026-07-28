<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api, type Note } from './api'

const notes = ref<Note[]>([])
const newText = ref('')
const loading = ref(false)
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
  if (!text) return
  try {
    await api.create(text)
    newText.value = ''
    await load()
  } catch (e) {
    error.value = (e as Error).message
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
    <h1>Заметки (demo: Vue + TS + Go + Postgres)</h1>

    <form @submit.prevent="addNote">
      <input v-model="newText" type="text" placeholder="Новая заметка" />
      <button type="submit">Добавить</button>
    </form>

    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="loading">Загрузка...</p>

    <ul>
      <li v-for="note in notes" :key="note.id">
        <span class="text">{{ note.text }}</span>
        <span class="date">{{ formatDate(note.created_at) }}</span>
        <button @click="removeNote(note.id)">Удалить</button>
      </li>
    </ul>

    <p v-if="!loading && notes.length === 0">Заметок пока нет.</p>
  </main>
</template>

<style scoped>
main {
  max-width: 640px;
  margin: 2rem auto;
  font-family: system-ui, sans-serif;
  padding: 0 1rem;
}

form {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

input {
  flex: 1;
  padding: 0.5rem;
}

button {
  padding: 0.5rem 1rem;
  cursor: pointer;
}

ul {
  list-style: none;
  padding: 0;
}

li {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem 0;
  border-bottom: 1px solid #ddd;
}

.text {
  flex: 1;
}

.date {
  color: #888;
  font-size: 0.85rem;
}

.error {
  color: #c00;
}
</style>

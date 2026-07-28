export interface Note {
  id: number
  text: string
  created_at: string
}

async function handle<T>(res: Response): Promise<T> {
  if (!res.ok) {
    throw new Error(`Ошибка запроса: ${res.status} ${await res.text()}`)
  }
  if (res.status === 204) {
    return undefined as T
  }
  return (await res.json()) as T
}

export const api = {
  list(): Promise<Note[]> {
    return fetch('/api/notes').then((r) => handle<Note[]>(r))
  },
  create(text: string): Promise<Note> {
    return fetch('/api/notes', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ text }),
    }).then((r) => handle<Note>(r))
  },
  remove(id: number): Promise<void> {
    return fetch(`/api/notes/${id}`, { method: 'DELETE' }).then((r) => handle<void>(r))
  },
}

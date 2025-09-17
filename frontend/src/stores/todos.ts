import { defineStore } from 'pinia'
import { api } from '../api/client'
import type { Todo } from '../types'

export const useTodos = defineStore('todos', {
  state: () => ({
    items: [] as Todo[],
    loading: false,
    error: '' as string | null,
  }),
  actions: {
    async fetch() {
      this.loading = true
      this.error = null
      try {
        const { data } = await api.get<Todo[]>('/todos')
        this.items = data
      } catch {
        this.error = 'Failed to load todos'
      } finally {
        this.loading = false
      }
    },
    async add(title: string) {
      const { data } = await api.post<Todo>('/todos', { title })
      this.items.unshift(data)
    },
    async toggle(id: number, done: boolean) {
      const { data } = await api.patch<Todo>(`/todos/${id}`, { done })
      const idx = this.items.findIndex(t => t.id === id)
      if (idx >= 0) this.items[idx] = data
    },
    async rename(id: number, title: string) {
      const { data } = await api.patch<Todo>(`/todos/${id}`, { title })
      const idx = this.items.findIndex(t => t.id === id)
      if (idx >= 0) this.items[idx] = data
    },
    async remove(id: number) {
      await api.delete(`/todos/${id}`)
      this.items = this.items.filter(t => t.id !== id)
    },
  },
})

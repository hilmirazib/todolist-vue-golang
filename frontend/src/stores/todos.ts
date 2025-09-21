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
      this.error = null
      try {
        const { data } = await api.post<Todo>('/todos', { title })
        this.items.unshift(data)
      } catch (err: unknown) {
        this.error = 'Failed to create todo'
        throw err
      }
    },
    async toggle(id: number, done: boolean) {
      const idx = this.items.findIndex((t) => t.id === id)
      if (idx < 0) return

      const previous = { ...this.items[idx] }
      this.items[idx] = { ...this.items[idx], done }
      this.error = null

      try {
        const { data } = await api.patch<Todo>(`/todos/${id}`, { done })
        this.items[idx] = data
      } catch (err: unknown) {
        this.items[idx] = previous
        this.error = 'Failed to update todo status'
        throw err
      }
    },
    async rename(id: number, title: string) {
      const idx = this.items.findIndex((t) => t.id === id)
      if (idx < 0) return

      const previous = { ...this.items[idx] }
      this.items[idx] = { ...this.items[idx], title }
      this.error = null

      try {
        const { data } = await api.patch<Todo>(`/todos/${id}`, { title })
        this.items[idx] = data
      } catch (err: unknown) {
        this.items[idx] = previous
        this.error = 'Failed to rename todo'
        throw err
      }
    },
    async remove(id: number) {
      const previous = [...this.items]
      this.items = this.items.filter((t) => t.id !== id)
      this.error = null

      try {
        await api.delete(`/todos/${id}`)
      } catch (err: unknown) {
        this.items = previous
        this.error = 'Failed to remove todo'
        throw err
      }
    },
  },
})

<template>
  <div class="pa-6" style="max-width: 720px; margin: 0 auto;">
    <v-card elevation="2">
      <v-card-title class="text-h6">Golang TodoList</v-card-title>
      <v-card-text>
        <div class="d-flex ga-2 mb-4">
          <v-text-field
            v-model="title"
            label="Add a task..."
            density="comfortable"
            hide-details="auto"
            @keyup.enter="create"
          />
          <v-btn color="primary" @click="create">Add</v-btn>
        </div>

        <v-alert
          v-if="store.error"
          type="error"
          variant="tonal"
          class="mb-3"
        >
          {{ store.error }}
        </v-alert>

        <v-skeleton-loader
          v-if="store.loading"
          type="list-item-two-line, list-item-two-line, list-item-two-line"
          class="mb-4"
        />

        <v-list v-else :lines="'one'" class="py-0">
          <template v-if="store.items.length">
            <v-list-item
              v-for="t in store.items"
              :key="t.id"
              class="py-2"
            >
              <template #prepend>
                <v-checkbox
                  :model-value="t.done"
                  hide-details
                  density="compact"
                  @update:model-value="(v:boolean) => toggle(t, v)"
                />
              </template>

              <v-list-item-title>
                <span :style="t.done ? 'text-decoration: line-through; opacity:.6' : ''">
                  {{ t.title }}
                </span>
              </v-list-item-title>

              <template #append>
                <div class="d-flex ga-1">
                  <v-btn size="small" variant="text" @click="edit(t)">Edit</v-btn>
                  <v-btn
                    size="small"
                    variant="text"
                    color="error"
                    @click="askDelete(t)"
                  >
                    Delete
                  </v-btn>
                </div>
              </template>
            </v-list-item>
          </template>

          <template v-else>
            <div class="text-medium-emphasis text-center py-6">
              No tasks yet.
            </div>
          </template>
        </v-list>
      </v-card-text>
    </v-card>

    <!-- Dialog Edit -->
    <v-dialog v-model="editing.visible" max-width="480">
      <v-card>
        <v-card-title class="text-h6">Edit task</v-card-title>
        <v-card-text>
          <v-text-field v-model="editing.title" label="Title" hide-details="auto" />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="editing.visible = false">Cancel</v-btn>
          <v-btn color="primary" @click="saveEdit">Save</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Dialog Konfirmasi Delete -->
    <v-dialog v-model="confirmDelete.visible" max-width="420">
      <v-card>
        <v-card-title class="text-h6">Delete this task?</v-card-title>
        <v-card-text>
          This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="confirmDelete.visible = false">Cancel</v-btn>
          <v-btn color="error" @click="confirmRemove">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useTodos } from '../stores/todos'

const store = useTodos()

const title = ref('')
const editing = reactive<{ visible: boolean; id: number | null; title: string }>({
  visible: false,
  id: null,
  title: ''
})

const confirmDelete = reactive<{ visible: boolean; id: number | null }>({
  visible: false,
  id: null
})

onMounted(() => {
  store.fetch()
})

const create = async () => {
  const value = title.value.trim()
  if (!value) return
  await store.add(value)
  title.value = ''
}

const toggle = async (t: any, v: boolean) => {
  await store.toggle(t.id, v)
}

const edit = (t: any) => {
  editing.visible = true
  editing.id = t.id
  editing.title = t.title
}

const saveEdit = async () => {
  const val = editing.title.trim()
  if (!val || editing.id == null) return
  await store.rename(editing.id, val)
  editing.visible = false
}

const askDelete = (t: any) => {
  confirmDelete.visible = true
  confirmDelete.id = t.id
}

const confirmRemove = async () => {
  if (confirmDelete.id == null) return
  await store.remove(confirmDelete.id)
  confirmDelete.visible = false
  confirmDelete.id = null
}
</script>

<style scoped>
/* opsional styling tambahan */
</style>

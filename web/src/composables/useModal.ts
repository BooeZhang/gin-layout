import type { Ref, WritableComputedRef } from 'vue'

interface ModalExpose {
  okLoading: boolean
  open: (options?: Record<string, unknown>) => void
}

export function useModal(): [Ref<ModalExpose | null>, WritableComputedRef<boolean | undefined>] {
  const modalRef = ref<ModalExpose | null>(null)
  const okLoading = computed({
    get() {
      return modalRef.value?.okLoading
    },
    set(v: boolean | undefined) {
      if (modalRef.value)
        modalRef.value.okLoading = !!v
    },
  })
  return [modalRef, okLoading]
}

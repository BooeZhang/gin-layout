import type { FormInst, FormRules } from 'naive-ui'
import type { Ref } from 'vue'
import { cloneDeep } from 'lodash-es'

type FormRecord = Record<string, unknown>

export function useForm<T extends FormRecord = FormRecord>(initFormData = {} as T): [
  Ref<FormInst | null>,
  Ref<T>,
  () => ReturnType<FormInst['validate']> | undefined,
  FormRules,
] {
  const formRef = ref<FormInst | null>(null)
  const formModel = ref(cloneDeep(initFormData)) as Ref<T>
  const rules: FormRules = {
    required: {
      required: true,
      message: '此为必填项',
      trigger: ['blur', 'change'],
    },
  }
  const validation = (): ReturnType<FormInst['validate']> | undefined => {
    return formRef.value?.validate()
  }
  return [formRef, formModel, validation, rules]
}

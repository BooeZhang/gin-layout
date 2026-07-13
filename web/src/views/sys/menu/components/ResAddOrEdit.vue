<template>
  <MeModal ref="modalRef">
    <n-form
      ref="modalFormRef"
      label-placement="left"
      require-mark-placement="left"
      :label-width="100"
      :model="modalForm"
    >
      <n-grid :cols="24" :x-gap="24">
        <n-form-item-gi :span="12" label="所属菜单" path="parentId">
          <n-tree-select
            v-model:value="modalForm.parentId"
            :options="menuOptions"
            :disabled="parentIdDisabled"
            label-field="name"
            key-field="id"
            placeholder="根菜单"
            clearable
          />
        </n-form-item-gi>
        <n-form-item-gi :span="12" path="name" :rule="required">
          <template #label>
            <QuestionLabel label="名称" content="标题" />
          </template>
          <n-input v-model:value="modalForm.name" />
        </n-form-item-gi>
        <n-form-item-gi :span="12" path="code" :rule="required">
          <template #label>
            <QuestionLabel
              label="编码"
              content="如果是菜单则对应前端路由的name，使用大驼峰"
            />
          </template>
          <n-input
            v-model:value="modalForm.code"
            :disabled="modalAction === 'edit'"
          />
        </n-form-item-gi>
        <n-form-item-gi
          v-if="isMenuForm"
          :span="12"
          path="path"
          :rule="{
            trigger: ['blur', 'change'],
            type: 'string',
            message: '必须是/、http、https开头',
            validator(rule, value) {
              if (value) {
                return /\/|http|https/.test(value);
              }
              return true;
            },
          }"
        >
          <template #label>
            <QuestionLabel label="路由地址" content="父级菜单可不填" />
          </template>
          <n-input v-model:value="modalForm.path" />
        </n-form-item-gi>
        <n-form-item-gi v-if="isMenuForm" :span="12" path="icon">
          <template #label>
            <QuestionLabel
              label="菜单图标"
              content="如material-symbols:help，图标库地址: https://icones.js.org/collection/all"
            />
          </template>
          <n-select
            v-model:value="modalForm.icon"
            :options="iconOptions"
            clearable
            filterable
          />
        </n-form-item-gi>
        <n-form-item-gi v-if="isMenuForm" :span="24" path="component">
          <template #label>
            <QuestionLabel
              label="组件路径"
              content="前端组件的路径，以 /src 开头，父级菜单可不填"
            />
          </template>
          <n-select
            v-model:value="modalForm.component"
            :options="componentOptions"
            clearable
            filterable
            tag
          />
        </n-form-item-gi>
        <n-form-item-gi v-if="isPermissionForm" :span="12" label="请求方法" path="method">
          <n-select
            v-model:value="modalForm.method"
            :options="methodOptions"
            clearable
          />
        </n-form-item-gi>
        <n-form-item-gi v-if="isPermissionForm" :span="12" label="接口路径" path="apiPath">
          <n-input v-model:value="modalForm.apiPath" placeholder="/api/v1/users" />
        </n-form-item-gi>
        <n-form-item-gi v-if="isPermissionForm" :span="12" label="权限编码" path="permCode">
          <n-input v-model:value="modalForm.permCode" placeholder="sys:user:create" />
        </n-form-item-gi>

        <n-form-item-gi v-if="isMenuForm" :span="12" path="hidden">
          <template #label>
            <QuestionLabel
              label="显示状态"
              content="控制是否在菜单栏显示，不影响路由注册"
            />
          </template>
          <n-switch
            v-model:value="modalForm.hidden"
            :checked-value="false"
            :unchecked-value="true"
          >
            <template #checked> 显示 </template>
            <template #unchecked> 隐藏 </template>
          </n-switch>
        </n-form-item-gi>
        <n-form-item-gi :span="12" path="enabled">
          <template #label>
            <QuestionLabel
              label="状态"
              content="如果是菜单，禁用后将不添加到路由表，无法进入此页面"
            />
          </template>
          <n-switch v-model:value="modalForm.enabled">
            <template #checked> 启用 </template>
            <template #unchecked> 禁用 </template>
          </n-switch>
        </n-form-item-gi>
        <n-form-item-gi v-if="isMenuForm" :span="12" path="cache">
          <template #label>
            <QuestionLabel
              label="KeepAlive"
              content="设置keepAlive需将组件的name设置成当前菜单的code"
            />
          </template>
          <n-switch v-model:value="modalForm.cache">
            <template #checked> 是 </template>
            <template #unchecked> 否 </template>
          </n-switch>
        </n-form-item-gi>
        <n-form-item-gi
          v-if="isMenuForm"
          :span="12"
          label="排序"
          path="sort"
          :rule="{
            type: 'number',
            required: true,
            message: '此为必填项',
            trigger: ['blur', 'change'],
          }"
        >
          <n-input-number v-model:value="modalForm.sort" />
        </n-form-item-gi>
      </n-grid>
    </n-form>
  </MeModal>
</template>

<script setup>
import icons from "isme:icons";
import pagePathes from "isme:page-pathes";
import { MeModal } from "@/components";
import { useForm, useModal } from "@/composables";
import api from "@/api/sys/resource";
import QuestionLabel from "./QuestionLabel.vue";

const props = defineProps({
  menus: {
    type: Array,
    required: true,
  },
});
const emit = defineEmits(["refresh"]);

const menuOptions = computed(() => {
  return [{ name: "根菜单", id: "", children: props.menus || [] }];
});
const componentOptions = pagePathes.map((path) => ({
  label: path,
  value: path,
}));
const iconOptions = icons.map((item) => ({
  label: () =>
    h("span", { class: "flex items-center" }, [
      h("i", { class: `${item} text-18 mr-8` }),
      item,
    ]),
  value: item,
}));
const methodOptions = ["GET", "POST", "PUT", "PATCH", "DELETE"].map((item) => ({
  label: item,
  value: item,
}));
const required = {
  required: true,
  message: "此为必填项",
  trigger: ["blur", "change"],
};

const defaultForm = { enabled: true, hidden: false, cache: true };
const [modalFormRef, modalForm, validation] = useForm();
const [modalRef, okLoading] = useModal();

const modalAction = ref("");
const parentIdDisabled = ref(false);
function handleOpen(options = {}) {
  const { action, row = {}, ...rest } = options;
  modalAction.value = action;
  modalForm.value = { ...defaultForm, ...row, layout: "" };
  parentIdDisabled.value = !!row.parentId && row.type === "button";
  modalRef.value.open({ ...rest, onOk: onSave });
}

async function onSave() {
  await validation();
  okLoading.value = true;
  try {
    let newFormData;
    if (!modalForm.value.parentId) modalForm.value.parentId = null;
    if (modalAction.value === "add") {
      newFormData = await api.addMenu(modalForm.value);
    } else if (modalAction.value === "edit") {
      await api.updateMenu(modalForm.value.id, modalForm.value);
    }
    okLoading.value = false;
    $message.success("保存成功");
    emit(
      "refresh",
      modalAction.value === "add" ? newFormData : modalForm.value,
    );
  } catch (error) {
    console.error(error);
    okLoading.value = false;
    return false;
  }
}

defineExpose({
  handleOpen,
});

const isMenuForm = computed(
  () => modalForm.value.type === "menu" || modalForm.value.type === "catalog",
);
const isPermissionForm = computed(
  () => modalForm.value.type === "menu" || modalForm.value.type === "button",
);
</script>

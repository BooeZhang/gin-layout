<template>
  <CommonPage>
    <template #action>
      <NButton v-access="'AddUser'" type="primary" @click="handleAdd()">
        <i class="i-material-symbols:add mr-4 text-18" />
        创建新用户
      </NButton>
    </template>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :scroll-x="1200"
      :columns="columns"
      :get-data="api.list"
    >
      <MeQueryItem label="用户名" :label-width="50">
        <n-input
          v-model:value="queryItems.account"
          type="text"
          placeholder="请输入用户名"
          clearable
        />
      </MeQueryItem>

      <MeQueryItem label="性别" :label-width="50">
        <n-select
          v-model:value="queryItems.gender"
          clearable
          :options="genders"
        />
      </MeQueryItem>

      <MeQueryItem label="状态" :label-width="50">
        <n-select
          v-model:value="queryItems.enabled"
          clearable
          :options="[
            { label: '启用', value: 1 },
            { label: '停用', value: 0 },
          ]"
        />
      </MeQueryItem>
    </MeCrud>

    <MeModal ref="modalRef" width="520px">
      <n-form
        ref="modalFormRef"
        label-placement="left"
        label-align="left"
        :label-width="80"
        :model="modalForm"
        :disabled="modalAction === 'view'"
      >
        <n-form-item
          label="用户名"
          path="account"
          :rule="{
            required: true,
            message: '请输入用户名',
            trigger: ['input', 'blur'],
          }"
        >
          <n-input
            v-model:value="modalForm.account"
            :disabled="modalAction !== 'add'"
          />
        </n-form-item>
        <n-form-item
          v-if="['add', 'edit'].includes(modalAction)"
          :label="modalAction === 'add' ? '初始密码' : '密码'"
          path="password"
          :rule="{
            required: modalAction === 'add',
            message:
              modalAction === 'add' ? '请输入密码' : '留空则保持原密码不变',
            trigger: ['input', 'blur'],
          }"
        >
          <n-input
            v-model:value="modalForm.password"
            type="password"
            show-password-on="mousedown"
            :placeholder="
              modalAction === 'add' ? '请输入初始密码' : '留空则保持原密码不变'
            "
          />
        </n-form-item>

        <n-form-item
          v-if="['add', 'edit'].includes(modalAction)"
          label="角色"
          path="roleIds"
        >
          <n-select
            v-model:value="modalForm.roleIds"
            :options="roles"
            label-field="name"
            value-field="id"
            clearable
            filterable
            multiple
          />
        </n-form-item>
        <n-form-item
          v-if="['add', 'edit'].includes(modalAction)"
          label="手机号"
          path="phone"
        >
          <n-input v-model:value="modalForm.phone" placeholder="请输入手机号" />
        </n-form-item>
        <n-form-item
          v-if="['add', 'edit'].includes(modalAction)"
          label="状态"
          path="enabled"
        >
          <NSwitch v-model:value="modalForm.enabled">
            <template #checked> 启用 </template>
            <template #unchecked> 停用 </template>
          </NSwitch>
        </n-form-item>
      </n-form>
    </MeModal>
  </CommonPage>
</template>

<script setup lang="ts">
import type { DataTableColumns } from "naive-ui";
import type { Ref } from "vue";
import type { PasswordPayload } from "@/types/auth";
import type { Id } from "@/types/common";
import type { RoleRecord } from "@/types/role";
import type { UserInfo, UserPayload } from "@/types/user";
import { h, onMounted, ref } from "vue";
import { NAvatar, NButton, NSwitch, NTag } from "naive-ui";
import { MeCrud, MeModal, MeQueryItem } from "@/components";
import { useCrud } from "@/composables";
import { withAccess } from "@/directives";
import { formatDateTime } from "@/utils";
import api from "@/api/sys/user";
import defaultAvatar from "@/assets/images/avatar.png";

defineOptions({ name: "UserMgt" });

type GenderValue = 1 | 2;
type UserModalAction = "add" | "edit" | "view";

interface QueryItems {
  account?: string;
  gender?: GenderValue;
  enabled?: 0 | 1;
}

interface UserForm extends UserPayload, PasswordPayload {
  account: string;
  password: string;
  roleIds: Id[];
  phone: string;
  enabled: boolean;
}

interface UserRow extends UserInfo {
  id: Id;
  account?: string;
  username?: string;
  phone?: string;
  email?: string;
  gender?: GenderValue;
  enabled?: boolean;
  enableLoading?: boolean;
  createTime?: string | number | Date;
  createdAt?: string;
  created_at?: string;
  roles?: RoleRecord[];
}

interface CrudTableExpose {
  handleSearch: () => void;
}

type UserTableColumn = DataTableColumns<UserRow>[number] & {
  hideInExcel?: boolean;
};

const $table = ref<CrudTableExpose | null>(null);
const queryItems = ref<QueryItems>({});

onMounted(() => {
  $table.value?.handleSearch();
});

const genders = [
  { label: "男", value: 1 },
  { label: "女", value: 2 },
] as const;
const roles = ref<RoleRecord[]>([]);
api.getAllRoles().then((data = []) => (roles.value = data));

const {
  modalForm,
  modalRef,
  modalFormRef,
  modalAction: baseModalAction,
  handleAdd,
  handleDelete,
  handleOpen,
  handleSave,
} = useCrud<UserForm>({
  name: "用户",
  initForm: {
    account: "",
    password: "",
    roleIds: [] as number[],
    phone: "",
    enabled: true,
  },
  doCreate: api.create,
  doDelete: api.delete,
  doUpdate: api.update,
  refresh: () => $table.value?.handleSearch(),
});

const modalAction = baseModalAction as Ref<UserModalAction | "">;
const openUserModal = (options: {
  action?: UserModalAction;
  row?: Partial<UserForm>;
  title?: string;
  onOk?: () => Promise<unknown> | unknown;
}) => {
  const row = options.row as Partial<UserForm & UserRow> | undefined;
  const normalizedRow = row
    ? {
        ...row,
        roleIds: row.roleIds ?? row.roles?.map((item) => item.id) ?? [],
      }
    : undefined;

  return handleOpen({
    ...options,
    row: normalizedRow,
    onOk: options.onOk ?? (options.action === "edit" ? onSave : undefined),
  } as Parameters<typeof handleOpen>[0]);
};

const columns: UserTableColumn[] = [
  {
    title: "头像",
    key: "avatar",
    width: 80,
    render: (row) =>
      h(NAvatar, {
        size: "medium",
        src: row.avatar || defaultAvatar,
      }),
  },
  { title: "用户名", key: "account", width: 150, ellipsis: { tooltip: true } },
  {
    title: "角色",
    key: "roles",
    width: 200,
    ellipsis: { tooltip: true },
    render: (row) => {
      const userRoles = row.roles;
      if (userRoles?.length) {
        return userRoles.map((item, index) =>
          h(
            NTag,
            { type: "success", style: index > 0 ? "margin-left: 8px;" : "" },
            { default: () => item.name },
          ),
        );
      }
      return "暂无角色";
    },
  },
  {
    title: "性别",
    key: "gender",
    width: 80,
    render: (row) =>
      genders.find((item) => row.gender === item.value)?.label ?? "",
  },
  { title: "邮箱", key: "email", width: 150, ellipsis: { tooltip: true } },
  { title: "手机", key: "phone", width: 150, ellipsis: { tooltip: true } },
  {
    title: "创建时间",
    key: "createDate",
    width: 180,
    render(row) {
      return h("span", formatDateTime(row.createdAt ?? row.created_at));
    },
  },
  {
    title: "状态",
    key: "enabled",
    width: 120,
    render: (row) =>
      h(
        NSwitch,
        {
          size: "small",
          rubberBand: false,
          value: row.enabled,
          loading: !!row.enableLoading,
          onUpdateValue: () => handleEnable(row),
        },
        {
          checked: () => "启用",
          unchecked: () => "停用",
        },
      ),
  },
  {
    title: "操作",
    key: "actions",
    width: 420,
    align: "right",
    fixed: "right",
    hideInExcel: true,
    render(row) {
      return [
        withAccess(
          h(
            NButton,
            {
              size: "small",
              type: "primary",
              secondary: true,
            },
            {
              default: () => "超管专属",
              icon: () => h("i", { class: "i-carbon:user-role text-14" }),
            },
          ),
          "SuperAdmin",
        ),
        h(
          NButton,
          {
            size: "small",
            type: "primary",
            style: "margin-left: 12px;",
            secondary: true,
            onClick: () =>
              openUserModal({ action: "edit", title: "编辑用户", row }),
          },
          {
            default: () => "编辑",
            icon: () =>
              h("i", { class: "i-material-symbols:edit-outline text-14" }),
          },
        ),
        h(
          NButton,
          {
            size: "small",
            type: "error",
            style: "margin-left: 12px;",
            disabled: row.account === "admin",
            onClick: () => handleDelete(row.id),
          },
          {
            default: () => "删除",
            icon: () =>
              h("i", { class: "i-material-symbols:delete-outline text-14" }),
          },
        ),
      ];
    },
  },
];

async function handleEnable(row: UserRow) {
  row.enableLoading = true;
  try {
    await api.update({ id: row.id, enabled: !row.enabled });
    row.enableLoading = false;
    $message.success("操作成功");
    $table.value?.handleSearch();
  } catch (error) {
    console.error(error);
    row.enableLoading = false;
  }
}

function onSave() {
  if (modalAction.value === "edit") {
    return handleSave({
      api: () => {
        const payload: UserPayload = { ...modalForm.value };
        if (!payload.password) {
          delete payload.password;
        }
        return api.update(payload);
      },
      cb: () => {
        $message.success("保存成功");
        $table.value?.handleSearch();
      },
    });
  }
  handleSave();
}
</script>

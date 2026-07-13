<template>
  <CommonPage>
    <div class="flex">
      <n-spin size="small" :show="treeLoading">
        <MenuTree
          v-model:current-menu="selectedMenu"
          class="w-320 shrink-0"
          :tree-data="menuTreeData"
          @refresh="initData"
        />
      </n-spin>

      <div class="ml-40 w-0 flex-1">
        <template v-if="currentMenu">
          <div class="flex justify-between">
            <h3 class="mb-12">
              {{ currentMenu.name }}
            </h3>
            <NButton
              size="small"
              type="primary"
              @click="handleEdit(currentMenu)"
            >
              <i class="i-material-symbols:edit-outline mr-4 text-14" />
              编辑
            </NButton>
          </div>
          <n-descriptions label-placement="left" bordered :column="2">
            <n-descriptions-item label="编码">
              {{ currentMenu.code }}
            </n-descriptions-item>
            <n-descriptions-item label="名称">
              {{ currentMenu.name }}
            </n-descriptions-item>
            <n-descriptions-item label="路由地址">
              {{ currentMenu.path ?? "--" }}
            </n-descriptions-item>
            <n-descriptions-item label="组件路径">
              {{ currentMenu.component ?? "--" }}
            </n-descriptions-item>
            <n-descriptions-item label="菜单图标">
              <span v-if="currentMenu.icon" class="flex items-center">
                <i :class="`${currentMenu.icon}?mask text-22 mr-8`" />
                <span class="opacity-50">{{ currentMenu.icon }}</span>
              </span>
              <span v-else>无</span>
            </n-descriptions-item>
            <n-descriptions-item label="是否显示">
              {{ currentMenu.hidden ? "否" : "是" }}
            </n-descriptions-item>
            <n-descriptions-item label="是否启用">
              {{ currentMenu.enabled ? "是" : "否" }}
            </n-descriptions-item>
            <n-descriptions-item label="接口映射">
              {{ currentMenu.method && currentMenu.apiPath ? `${currentMenu.method} ${currentMenu.apiPath}` : "--" }}
            </n-descriptions-item>
            <n-descriptions-item label="权限编码">
              {{ currentMenu.permCode || "--" }}
            </n-descriptions-item>
            <n-descriptions-item label="KeepAlive">
              {{ currentMenu.cache ? "是" : "否" }}
            </n-descriptions-item>
            <n-descriptions-item label="排序">
              {{ currentMenu.sort ?? "--" }}
            </n-descriptions-item>
          </n-descriptions>

          <div class="mt-32 flex justify-between">
            <h3 class="mb-12">按钮</h3>
            <NButton size="small" type="primary" @click="handleAddBtn">
              <i class="i-fe:plus mr-4 text-14" />
              新增
            </NButton>
          </div>

          <n-data-table
            :columns="btnsColumns"
            :data="buttonList"
            :pagination="false"
            :scroll-x="1"
            remote
          />
        </template>
        <n-empty
          v-else
          class="h-450 f-c-c"
          size="large"
          description="请选择菜单查看详情"
        />
      </div>
    </div>
    <ResAddOrEdit ref="modalRef" :menus="treeData" @refresh="initData" />
  </CommonPage>
</template>

<script setup>
import { computed } from "vue";
import { NButton, NSwitch } from "naive-ui";
import api from "@/api/sys/resource";
import MenuTree from "./components/MenuTree.vue";
import ResAddOrEdit from "./components/ResAddOrEdit.vue";

const treeData = ref([]);
const treeLoading = ref(false);
const currentMenu = ref(null);
const menuTreeData = computed(() => filterMenuTree(treeData.value));
const selectedMenu = computed({
  get: () => currentMenu.value,
  set: (value) => {
    if (!value) {
      currentMenu.value = null;
      return;
    }
    currentMenu.value =
      findMenuNode(treeData.value, value.id, value.code) || value;
  },
});
const buttonList = computed(() =>
  (currentMenu.value?.children || []).filter((item) => item.type === "button"),
);

async function initData(data) {
  treeLoading.value = true;
  treeData.value = await api.getMenuTree();
  treeLoading.value = false;

  if (data) currentMenu.value = data;
}
initData();

const modalRef = ref(null);
function handleEdit(item = {}) {
  modalRef.value?.handleOpen({
    action: "edit",
    title: `编辑菜单 - ${item.name}`,
    row: item,
    okText: "保存",
  });
}

const btnsColumns = [
  { title: "名称", key: "name" },
  { title: "编码", key: "code" },
  {
    title: "接口",
    key: "apiPath",
    render: (row) => row.apiPath ? `${row.method || "--"} ${row.apiPath}` : "--",
  },
  {
    title: "状态",
    key: "enabled",
    render: (row) =>
      h(
        NSwitch,
        {
          size: "small",
          rubberBand: false,
          value: row.enabled,
          loading: !!row.enabledLoading,
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
    width: 320,
    align: "right",
    fixed: "right",
    render(row) {
      return [
        h(
          NButton,
          {
            size: "small",
            type: "primary",
            style: "margin-left: 12px;",
            onClick: () => handleEditBtn(row),
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
            onClick: () => handleDeleteBtn(row.id),
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

function handleAddBtn() {
  modalRef.value?.handleOpen({
    action: "add",
    title: "新增按钮",
    row: { type: "button", parentId: currentMenu.value.id },
    okText: "保存",
  });
}

function handleEditBtn(row) {
  modalRef.value?.handleOpen({
    action: "edit",
    title: `编辑按钮 - ${row.name}`,
    row,
    okText: "保存",
  });
}

function handleDeleteBtn(id) {
  const d = $dialog.warning({
    content: "确定删除？",
    title: "提示",
    positiveText: "确定",
    negativeText: "取消",
    async onPositiveClick() {
      try {
        d.loading = true;
        await api.deleteMenu(id);
        $message.success("删除成功");
        await initData(currentMenu.value);
        d.loading = false;
      } catch (error) {
        console.error(error);
        d.loading = false;
      }
    },
  });
}

async function handleEnable(item) {
  try {
    item.enabledLoading = true;
    await api.updateMenu(item.id, {
      enabled: !item.enabled,
    });
    $message.success("操作成功");
    item.enabled = !item.enabled;
    item.enabledLoading = false;
  } catch (error) {
    console.error(error);
    item.enabledLoading = false;
  }
}

function filterMenuTree(nodes = []) {
  return nodes
    .filter((item) => item.type !== "button")
    .map((item) => ({
      ...item,
      children: filterMenuTree(item.children || []),
    }));
}

function findMenuNode(nodes = [], id, code) {
  for (const item of nodes) {
    if ((id !== undefined && item.id === id) || (code && item.code === code)) {
      return item;
    }
    const child = findMenuNode(item.children || [], id, code);
    if (child) {
      return child;
    }
  }
  return null;
}
</script>

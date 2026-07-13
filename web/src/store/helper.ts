import { cloneDeep } from "lodash-es";
import userApi from "@/api/sys/user";
import { baseAccessMenus } from "@/settings";
import type { AccessMenuRecord } from "@/types/app";
import type { UserInfo } from "@/types/user";

const hiddenMenuCodes = new Set([
  "ExternalLink",
  "ShowDocs",
  "ApiFoxDocs",
  "NaiveUI",
  "MyBlog",
  "Base",
  "Demo",
  "BusinessDemo",
  "BusinessExample",
  "ImgUpload",
  "TestModal",
  "Unocss",
  "UnocssIcon",
  "KeepAlive",
]);

// 获取用户信息
export async function getUserInfo(): Promise<UserInfo> {
  return userApi.getUser();
}

// 获取用户可访问菜单
export async function getUserAccessMenus() {
  let remoteAccessMenus: AccessMenuRecord[] = [];
  try {
    remoteAccessMenus = await userApi.getUserMenu();
  } catch (error) {
    console.error(error);
  }
  return cloneDeep(baseAccessMenus).concat(filterHiddenAccessMenus(remoteAccessMenus));
}

function filterHiddenAccessMenus(
  accessMenus: AccessMenuRecord[],
): AccessMenuRecord[] {
  return accessMenus.reduce<AccessMenuRecord[]>((result, accessMenu) => {
    if (shouldHideAccessMenu(accessMenu)) return result;

    const children = accessMenu.children
      ? filterHiddenAccessMenus(accessMenu.children)
      : undefined;

    result.push({
      ...accessMenu,
      ...(children ? { children } : {}),
    });

    return result;
  }, []);
}

function shouldHideAccessMenu(accessMenu: AccessMenuRecord): boolean {
  return (
    (accessMenu.type === "catalog" || accessMenu.type === "menu") &&
    (hiddenMenuCodes.has(accessMenu.code) ||
      isRemovedMenuPath(accessMenu.path) ||
      isRemovedMenuPath(accessMenu.component))
  );
}

function isRemovedMenuPath(path?: string): boolean {
  return (
    !!path &&
    (/^https?:|mailto:|tel:/.test(path) ||
      path.includes("/src/views/base/") ||
      path.includes("/src/views/demo/") ||
      path.includes("/src/views/iframe/") ||
      path.startsWith("/base") ||
      path.startsWith("/demo") ||
      path.startsWith("/iframe"))
  );
}

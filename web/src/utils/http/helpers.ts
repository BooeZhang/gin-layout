import { useAuthStore } from "@/store";

let isConfirming = false;

// 处理登录过期
function handleAuthExpired(
  content: string,
  needTip: boolean,
): false | undefined {
  if (isConfirming || !needTip) return;
  isConfirming = true;
  $dialog.confirm({
    title: "提示",
    type: "info",
    content,
    confirm() {
      useAuthStore().logout();
      window.$message?.success("已退出登录");
      isConfirming = false;
    },
    cancel() {
      isConfirming = false;
    },
  });
  return false;
}

// 处理请求错误
export function resolveResError(
  code: number | string | undefined,
  message: string | undefined,
  needTip = true,
): string | false | undefined {
  switch (code) {
    case 401:
      return handleAuthExpired("登录已过期，是否重新登录？", needTip);
    case 11007:
    case 11008:
      return handleAuthExpired(`${message}，是否重新登录？`, needTip);
    case 403:
      message = "请求被拒绝";
      break;
    case 404:
      message = "请求资源或接口不存在";
      break;
    case 500:
      message = "服务器发生异常";
      break;
    default:
      message = message ?? `【${code}】: 未知异常!`;
      break;
  }
  needTip && window.$message?.error(message);
  return message;
}

<template>
  <div class="login-page">
    <div class="ambient-layer" aria-hidden="true">
      <span class="beam beam-a" />
      <span class="beam beam-b" />
      <span class="beam beam-c" />
    </div>

    <header class="login-topbar">
      <div class="brand-mark">
        <img src="@/assets/images/logo.png" alt="Gin Layout" class="brand-logo">
        <div>
          <strong>{{ title }}</strong>
          <span>Admin Studio</span>
        </div>
      </div>
      <div class="status-pill">
        <i class="i-fe:shield" />
        <span>Secure Access</span>
      </div>
    </header>

    <main class="login-shell">
      <section class="identity-panel" aria-label="系统概览">
        <div class="hero-copy">
          <p class="eyebrow">
            Control room is ready
          </p>
          <h1>从这里点亮你的管理现场。</h1>
          <p>
            不只是进入后台，更像进入一间为权限、用户和资源编排而生的工作室。
          </p>
        </div>

        <div class="identity-stage" aria-hidden="true">
          <div class="stage-ring ring-outer" />
          <div class="stage-ring ring-inner" />
          <div class="stage-core">
            <img src="@/assets/images/logo.png" alt="">
          </div>
          <div class="stage-tag tag-user">
            <i class="i-fe:users" />
            <span>User Layer</span>
          </div>
          <div class="stage-tag tag-role">
            <i class="i-fe:key" />
            <span>Role Matrix</span>
          </div>
          <div class="stage-tag tag-api">
            <i class="i-fe:git-branch" />
            <span>API Routes</span>
          </div>
        </div>
      </section>

      <section class="form-panel" aria-label="管理员登录">
        <div
          class="form-card"
          :class="{
            'is-account-focus': focusField === 'account',
            'is-password-focus': focusField === 'password',
            'is-typing': isTyping,
          }"
        >
          <div class="card-glow" aria-hidden="true" />
          <div class="form-heading">
            <p class="eyebrow">
              Sign in
            </p>
            <h2>欢迎回来</h2>
            <span>输入管理员身份，继续你的控制台会话。</span>
          </div>

          <div class="form-stack">
            <label class="field-group" for="login-account">
              <span>账号</span>
              <n-input
                id="login-account"
                v-model:value="loginInfo.account"
                autofocus
                class="login-input"
                placeholder="请输入用户名"
                :maxlength="20"
                @focus="setFocus('account')"
                @blur="clearFocus"
                @input="markTyping"
                @keydown.enter="handleLogin()"
              >
                <template #prefix>
                  <i class="input-icon i-fe:user" />
                </template>
              </n-input>
            </label>

            <label class="field-group" for="login-password">
              <span>密码</span>
              <n-input
                id="login-password"
                v-model:value="loginInfo.password"
                class="login-input"
                :type="passwordVisible ? 'text' : 'password'"
                placeholder="请输入密码"
                :maxlength="20"
                @focus="setFocus('password')"
                @blur="clearFocus"
                @input="markTyping"
                @keydown.enter="handleLogin()"
              >
                <template #prefix>
                  <i class="input-icon i-fe:lock" />
                </template>
                <template #suffix>
                  <button
                    class="password-toggle"
                    type="button"
                    :aria-label="passwordVisible ? '隐藏密码' : '显示密码'"
                    @mousedown.prevent
                    @click="passwordVisible = !passwordVisible"
                  >
                    <i :class="passwordVisible ? 'i-fe:eye-off' : 'i-fe:eye'" />
                  </button>
                </template>
              </n-input>
            </label>

            <!-- <div class="captcha-row">
              <n-input
                v-model:value="loginInfo.captcha"
                class="login-input captcha-input"
                placeholder="请输入验证码"
                :maxlength="4"
                @focus="setFocus('captcha')"
                @blur="clearFocus"
                @input="markTyping"
                @keydown.enter="handleLogin()"
              >
                <template #prefix>
                  <i class="input-icon i-fe:key" />
                </template>
              </n-input>
              <img
                v-if="captchaUrl"
                :src="captchaUrl"
                alt="验证码"
                class="captcha-image"
                @click="initCaptcha"
              />
            </div> -->

            <div class="login-meta">
              <n-checkbox
                :checked="isRemember"
                label="记住我"
                :on-update:checked="(val: boolean) => (isRemember = val)"
              />
              <button type="button" class="text-button" @click="initCaptcha">
                <i class="i-fe:refresh-cw" />
                刷新验证码
              </button>
            </div>

            <n-button
              class="login-button"
              type="primary"
              :loading="loading"
              @click="handleLogin()"
            >
              <template #icon>
                <i class="i-fe:log-in" />
              </template>
              登录
            </n-button>
          </div>

          <div class="form-footnote">
            <div>
              <i class="i-fe:activity" />
              <span>Session route</span>
            </div>
            <strong>{{ route.query.redirect ? 'Back to target page' : 'Console home' }}</strong>
          </div>
        </div>
      </section>
    </main>

    <TheFooter class="login-footer" />
  </div>
</template>

<script setup lang="ts">
import type { LoginPayload, TokenPayload } from '@/types/app'
import { useStorage } from '@vueuse/core'
import { onBeforeUnmount, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api/auth'
import { useAuthStore } from '@/store'
import { lStorage, throttle } from '@/utils'

type LoginFocusField = 'account' | 'password' | 'captcha' | ''

interface LoginError {
  code?: number
}

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const title = import.meta.env.VITE_TITLE

const loginInfo = reactive<LoginPayload>({
  account: '',
  password: '',
  captcha: '',
})

const captchaUrl = ref('')
const initCaptcha = throttle(() => {
  captchaUrl.value = `${import.meta.env.VITE_AXIOS_BASE_URL}/auth/captcha?${Date.now()}`
}, 500)

const localLoginInfo = lStorage.get(
  'loginInfo',
) as Partial<LoginPayload> | null
if (localLoginInfo) {
  loginInfo.account = localLoginInfo.account || ''
  loginInfo.password = localLoginInfo.password || ''
}
initCaptcha()

const isRemember = useStorage('isRemember', true)
const loading = ref(false)
const focusField = ref<LoginFocusField>('')
const passwordVisible = ref(false)
const isTyping = ref(false)
let typingTimer: ReturnType<typeof window.setTimeout> | undefined

function setFocus(field: LoginFocusField) {
  focusField.value = field
}

function clearFocus() {
  focusField.value = ''
}

function markTyping() {
  isTyping.value = true
  if (typingTimer)
    window.clearTimeout(typingTimer)
  typingTimer = window.setTimeout(() => {
    isTyping.value = false
  }, 550)
}

async function handleLogin() {
  const { account, password, captcha } = loginInfo
  if (!account || !password)
    return $message.warning('请输入用户名和密码')
  // if (!captcha) return $message.warning("请输入验证码");
  try {
    loading.value = true
    $message.loading('正在验证，请稍后...', { key: 'login' })
    const data = await api.login({
      account,
      password: password.toString(),
      captcha,
    })
    if (isRemember.value)
      lStorage.set('loginInfo', { account, password }, undefined)
    else lStorage.remove('loginInfo')
    onLoginSuccess(data)
  }
  catch (error) {
    if ((error as LoginError)?.code === 10003)
      initCaptcha()
    $message.destroy('login')
    console.error(error)
  }
  loading.value = false
}

async function onLoginSuccess(data: TokenPayload = {}) {
  authStore.setToken(data)
  $message.loading('登录中...', { key: 'login' })
  try {
    $message.success('登录成功', { key: 'login' })
    if (route.query.redirect) {
      const { redirect, ...query } = route.query
      router.push({ path: String(redirect), query })
    }
    else {
      router.push('/')
    }
  }
  catch (error) {
    console.error(error)
    $message.destroy('login')
  }
}

onBeforeUnmount(() => {
  if (typingTimer)
    window.clearTimeout(typingTimer)
})
</script>

<style scoped>
.login-page {
  position: relative;
  min-height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 28px;
  color: #f8fafc;
  background:
    radial-gradient(circle at 18% 18%, rgba(20, 184, 166, 0.24), transparent 30%),
    radial-gradient(circle at 82% 18%, rgba(232, 121, 249, 0.2), transparent 28%),
    radial-gradient(circle at 62% 88%, rgba(251, 191, 36, 0.14), transparent 34%),
    linear-gradient(135deg, #050816 0%, #111827 46%, #240f3d 100%);
}

.login-page::before {
  position: absolute;
  inset: 0;
  background:
    linear-gradient(90deg, rgba(255, 255, 255, 0.05) 1px, transparent 1px),
    linear-gradient(180deg, rgba(255, 255, 255, 0.04) 1px, transparent 1px);
  background-size: 56px 56px;
  content: '';
  mask-image: linear-gradient(180deg, rgba(0, 0, 0, 0.88), transparent 88%);
}

.ambient-layer {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.beam {
  position: absolute;
  display: block;
  width: 42vw;
  height: 18vh;
  border-radius: 999px;
  filter: blur(34px);
  opacity: 0.54;
  transform: rotate(-18deg);
}

.beam-a {
  left: -10%;
  top: 18%;
  background: rgba(20, 184, 166, 0.36);
}

.beam-b {
  right: -8%;
  top: 24%;
  background: rgba(129, 140, 248, 0.34);
  transform: rotate(18deg);
}

.beam-c {
  left: 32%;
  bottom: 5%;
  background: rgba(244, 114, 182, 0.24);
}

.login-topbar {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  width: min(1180px, 100%);
  margin: 0 auto 26px;
}

.brand-mark {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-logo {
  width: 44px;
  height: 44px;
  padding: 6px;
  border: 1px solid rgba(255, 255, 255, 0.34);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 18px 44px rgba(0, 0, 0, 0.22);
}

.brand-mark strong,
.brand-mark span {
  display: block;
}

.brand-mark strong {
  color: #ffffff;
  font-size: 18px;
  line-height: 1.2;
}

.brand-mark span {
  margin-top: 3px;
  color: rgba(226, 232, 240, 0.7);
  font-size: 12px;
  letter-spacing: 0;
  text-transform: uppercase;
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 38px;
  padding: 0 14px;
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 999px;
  color: #ccfbf1;
  background: rgba(15, 23, 42, 0.48);
  backdrop-filter: blur(18px);
  font-size: 13px;
  font-weight: 700;
}

.status-pill i {
  color: #2dd4bf;
  font-size: 17px;
}

.login-shell {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(360px, 430px);
  gap: 34px;
  align-items: center;
  width: min(1180px, 100%);
  min-height: 660px;
  margin: auto;
}

.identity-panel {
  position: relative;
  min-height: 620px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.hero-copy {
  position: relative;
  z-index: 1;
  max-width: 590px;
}

.eyebrow {
  margin: 0 0 12px;
  color: #5eead4;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0;
  text-transform: uppercase;
}

.hero-copy h1 {
  margin: 0;
  color: #ffffff;
  font-size: clamp(42px, 6vw, 72px);
  font-weight: 700;
  line-height: 1.02;
  letter-spacing: 0;
}

.hero-copy p:last-child {
  width: min(500px, 100%);
  margin: 18px 0 0;
  color: rgba(226, 232, 240, 0.72);
  font-size: 16px;
  line-height: 1.8;
}

.identity-stage {
  position: absolute;
  right: 0;
  bottom: 10px;
  width: min(430px, 58vw);
  aspect-ratio: 1;
}

.stage-ring {
  position: absolute;
  inset: 0;
  border: 1px solid rgba(255, 255, 255, 0.16);
  border-radius: 50%;
}

.ring-outer {
  background:
    conic-gradient(from 210deg, transparent, rgba(45, 212, 191, 0.34), transparent 36%), rgba(255, 255, 255, 0.02);
  box-shadow: inset 0 0 60px rgba(255, 255, 255, 0.04);
}

.ring-inner {
  inset: 19%;
  border-color: rgba(255, 255, 255, 0.22);
  background: rgba(15, 23, 42, 0.34);
  backdrop-filter: blur(12px);
}

.stage-core {
  position: absolute;
  left: 50%;
  top: 50%;
  display: grid;
  place-items: center;
  width: 104px;
  height: 104px;
  border: 1px solid rgba(255, 255, 255, 0.24);
  border-radius: 50%;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.24), rgba(255, 255, 255, 0.06)), rgba(15, 23, 42, 0.72);
  box-shadow: 0 24px 70px rgba(20, 184, 166, 0.24);
  transform: translate(-50%, -50%);
}

.stage-core img {
  width: 58px;
  height: 58px;
  border-radius: 8px;
}

.stage-tag {
  position: absolute;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-height: 38px;
  padding: 0 13px;
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 999px;
  color: #e2e8f0;
  background: rgba(15, 23, 42, 0.56);
  box-shadow: 0 18px 44px rgba(0, 0, 0, 0.18);
  backdrop-filter: blur(16px);
  font-size: 12px;
  font-weight: 700;
}

.stage-tag i {
  color: #f0abfc;
  font-size: 16px;
}

.tag-user {
  left: 4%;
  top: 22%;
}

.tag-role {
  right: 0;
  top: 42%;
}

.tag-api {
  left: 18%;
  bottom: 16%;
}

.form-panel {
  display: flex;
  align-items: center;
  justify-content: center;
}

.form-card {
  position: relative;
  width: min(410px, 100%);
  overflow: hidden;
  padding: 30px;
  border: 1px solid rgba(255, 255, 255, 0.22);
  border-radius: 8px;
  background: rgba(15, 23, 42, 0.58);
  box-shadow: 0 30px 90px rgba(0, 0, 0, 0.34);
  backdrop-filter: blur(22px);
  transition:
    border-color 0.18s ease,
    box-shadow 0.18s ease;
}

.form-card.is-account-focus {
  border-color: rgba(129, 140, 248, 0.64);
  box-shadow: 0 30px 90px rgba(99, 102, 241, 0.2);
}

.form-card.is-password-focus {
  border-color: rgba(45, 212, 191, 0.62);
  box-shadow: 0 30px 90px rgba(20, 184, 166, 0.18);
}

.form-card.is-typing .card-glow {
  opacity: 0.9;
}

.card-glow {
  position: absolute;
  inset: -1px;
  opacity: 0.62;
  background:
    linear-gradient(125deg, rgba(45, 212, 191, 0.22), transparent 36%),
    linear-gradient(315deg, rgba(244, 114, 182, 0.16), transparent 34%);
  pointer-events: none;
  transition: opacity 0.18s ease;
}

.form-heading,
.form-stack,
.form-footnote {
  position: relative;
  z-index: 1;
}

.form-heading {
  margin-bottom: 28px;
}

.form-heading h2 {
  margin: 0;
  color: #ffffff;
  font-size: 38px;
  font-weight: 700;
  line-height: 1.12;
  letter-spacing: 0;
}

.form-heading span {
  display: block;
  margin-top: 10px;
  color: rgba(226, 232, 240, 0.7);
  font-size: 14px;
}

.form-stack {
  display: grid;
  gap: 18px;
}

.field-group {
  display: grid;
  gap: 8px;
}

.field-group > span {
  color: rgba(226, 232, 240, 0.82);
  font-size: 13px;
  font-weight: 600;
}

.login-input {
  height: 48px;
  border-radius: 8px;
}

:deep(.login-input.n-input) {
  --n-height: 48px !important;
  --n-line-height-textarea: 48px !important;
  --n-text-color: #f8fafc !important;
  --n-caret-color: #5eead4 !important;
  --n-placeholder-color: rgba(203, 213, 225, 0.5) !important;
  --n-color: rgba(255, 255, 255, 0.08) !important;
  --n-color-focus: rgba(255, 255, 255, 0.1) !important;
  --n-border: 1px solid rgba(255, 255, 255, 0.18) !important;
  --n-border-hover: 1px solid rgba(94, 234, 212, 0.56) !important;
  --n-border-focus: 1px solid rgba(94, 234, 212, 0.72) !important;
  background: rgba(255, 255, 255, 0.08);
}

.login-input :deep(.n-input__border),
.login-input :deep(.n-input__state-border) {
  border-radius: 8px;
}

.login-input :deep(.n-input-wrapper),
.login-input :deep(.n-input__input) {
  height: 100%;
}

.login-input :deep(.n-input-wrapper) {
  align-items: center;
  padding-left: 14px;
  padding-right: 10px;
}

.login-input :deep(.n-input__input-el) {
  height: 48px;
  color: #f8fafc;
  caret-color: #5eead4;
  line-height: 48px;
  -webkit-text-fill-color: #f8fafc;
}

.login-input :deep(.n-input__placeholder) {
  color: rgba(203, 213, 225, 0.5);
}

.login-input :deep(.n-input__input-el::placeholder) {
  color: rgba(203, 213, 225, 0.5);
  -webkit-text-fill-color: rgba(203, 213, 225, 0.5);
}

.input-icon {
  margin-right: 10px;
  color: #5eead4;
  font-size: 18px;
}

.password-toggle,
.text-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 0;
  color: #5eead4;
  background: transparent;
  cursor: pointer;
  transition:
    color 0.18s ease,
    background-color 0.18s ease;
}

.password-toggle {
  width: 30px;
  height: 30px;
  border-radius: 8px;
  font-size: 18px;
}

.password-toggle:hover,
.password-toggle:focus-visible {
  color: #ffffff;
  background: rgba(45, 212, 191, 0.16);
}

.captcha-row {
  display: grid;
  grid-template-columns: 1fr 104px;
  gap: 12px;
  align-items: end;
}

.captcha-image {
  width: 104px;
  height: 48px;
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 8px;
  object-fit: cover;
  background: rgba(255, 255, 255, 0.08);
  cursor: pointer;
}

.login-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  min-height: 28px;
}

.text-button {
  gap: 6px;
  padding: 0;
  font-size: 13px;
  white-space: nowrap;
}

.text-button:hover,
.text-button:focus-visible {
  color: #ffffff;
}

.login-button {
  width: 100%;
  height: 48px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 700;
}

.login-button :deep(.n-button__content) {
  gap: 8px;
}

.form-footnote {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: 22px;
  padding: 14px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 8px;
  color: rgba(226, 232, 240, 0.76);
  background: rgba(15, 23, 42, 0.42);
  font-size: 12px;
  line-height: 1.6;
}

.form-footnote div {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.form-footnote i {
  color: #5eead4;
}

.form-footnote strong {
  color: #ffffff;
  font-size: 12px;
  white-space: nowrap;
}

.login-footer {
  position: relative;
  z-index: 1;
  flex: 0 0 auto;
  width: min(1180px, 100%);
  margin: 22px auto 0;
  color: rgba(226, 232, 240, 0.62);
}

@media (prefers-reduced-motion: reduce) {
  .form-card,
  .card-glow,
  .password-toggle,
  .text-button {
    transition: none;
  }
}

@media (max-width: 980px) {
  .login-page {
    overflow-y: auto;
  }

  .login-topbar {
    margin-bottom: 14px;
  }

  .login-shell {
    grid-template-columns: 1fr;
    gap: 22px;
    min-height: auto;
  }

  .identity-panel {
    min-height: 470px;
  }

  .identity-stage {
    right: 0;
    bottom: -10px;
    width: min(360px, 72vw);
  }
}

@media (max-width: 640px) {
  .login-page {
    padding: 14px;
  }

  .login-topbar {
    align-items: flex-start;
    flex-direction: column;
  }

  .status-pill {
    min-height: 34px;
  }

  .identity-panel {
    min-height: 420px;
  }

  .hero-copy h1 {
    font-size: 40px;
  }

  .hero-copy p:last-child {
    font-size: 14px;
  }

  .identity-stage {
    width: 300px;
    opacity: 0.72;
  }

  .stage-tag {
    display: none;
  }

  .form-card {
    padding: 24px;
  }

  .form-heading h2 {
    font-size: 32px;
  }

  .login-meta {
    align-items: flex-start;
    flex-direction: column;
  }

  .form-footnote {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>

<template>
  <div
    ref="stageRef"
    class="cute-login-characters"
    :class="{
      'is-password': focusField === 'password',
      'is-captcha': focusField === 'captcha',
      'is-typing': typing,
      'is-secret-visible': passwordVisible,
    }"
    @pointermove="handlePointerMove"
    @pointerleave="resetPointer"
  >
    <i class="sparkle sparkle-a" />
    <i class="sparkle sparkle-b" />
    <i class="sparkle sparkle-c" />
    <div class="ground" />

    <div class="buddy buddy-main" :style="eyeStyle">
      <i class="ear ear-left" />
      <i class="ear ear-right" />
      <i class="eye eye-left" />
      <i class="eye eye-right" />
      <i class="cheek cheek-left" />
      <i class="cheek cheek-right" />
      <i class="smile" />
      <i class="paw paw-left" />
      <i class="paw paw-right" />
    </div>

    <div class="buddy buddy-friend" :style="friendEyeStyle">
      <i class="ear ear-left" />
      <i class="ear ear-right" />
      <i class="eye eye-left" />
      <i class="eye eye-right" />
      <i class="cheek cheek-left" />
      <i class="cheek cheek-right" />
      <i class="smile" />
      <i class="paw paw-left" />
      <i class="paw paw-right" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

type LoginFocusField = 'account' | 'password' | 'captcha' | ''

interface Props {
  focusField?: LoginFocusField
  passwordVisible?: boolean
  typing?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  focusField: '',
  passwordVisible: false,
  typing: false,
})

const stageRef = ref<HTMLElement>()
const pointer = ref({ x: 0, y: 0 })

const focusOffset = computed(() => {
  if (props.passwordVisible)
    return { x: -4, y: -2 }
  if (props.focusField === 'password')
    return { x: 3, y: 1 }
  if (props.focusField === 'captcha')
    return { x: 0, y: 4 }
  if (props.focusField === 'account')
    return { x: -2, y: 0 }
  return pointer.value
})

const eyeStyle = computed(() => ({
  '--eye-x': `${focusOffset.value.x}px`,
  '--eye-y': `${focusOffset.value.y}px`,
}))

const friendEyeStyle = computed(() => ({
  '--eye-x': `${props.passwordVisible ? 4 : focusOffset.value.x * -0.6}px`,
  '--eye-y': `${props.passwordVisible ? 1 : focusOffset.value.y * 0.8}px`,
}))

function handlePointerMove(event: PointerEvent) {
  if (!stageRef.value || props.focusField)
    return
  const rect = stageRef.value.getBoundingClientRect()
  const x = ((event.clientX - rect.left) / rect.width - 0.5) * 8
  const y = ((event.clientY - rect.top) / rect.height - 0.5) * 6
  pointer.value = {
    x: Number(x.toFixed(2)),
    y: Number(y.toFixed(2)),
  }
}

function resetPointer() {
  pointer.value = { x: 0, y: 0 }
}
</script>

<style scoped>
.cute-login-characters {
  position: relative;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  min-height: 360px;
  padding: 44px 0 14px;
  isolation: isolate;
}

.ground {
  position: absolute;
  right: 14%;
  bottom: 18px;
  left: 14%;
  height: 30px;
  border-radius: 50%;
  background: rgba(31, 111, 120, 0.13);
  filter: blur(3px);
}

.sparkle {
  position: absolute;
  z-index: -1;
  width: 12px;
  height: 12px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.9);
  transform: rotate(45deg);
  animation: login-sparkle 2.8s ease-in-out infinite;
}

.sparkle-a {
  left: 14%;
  top: 18%;
}

.sparkle-b {
  right: 18%;
  top: 35%;
  animation-delay: -1s;
}

.sparkle-c {
  left: 30%;
  bottom: 27%;
  animation-delay: -1.7s;
}

.buddy {
  --eye-x: 0px;
  --eye-y: 0px;
  position: relative;
  width: 164px;
  height: 184px;
  margin: 0 -4px;
  border: 3px solid rgba(255, 255, 255, 0.82);
  border-radius: 48% 48% 38% 38%;
  background:
    radial-gradient(circle at 34% 22%, rgba(255, 255, 255, 0.72), transparent 12%),
    linear-gradient(165deg, #ffe39f, #ffb65f);
  box-shadow: 0 24px 42px rgba(40, 54, 75, 0.18);
  transform-origin: bottom center;
  animation: login-bob 3.4s ease-in-out infinite;
}

.buddy-friend {
  width: 148px;
  height: 166px;
  margin-left: -24px;
  background:
    radial-gradient(circle at 34% 22%, rgba(255, 255, 255, 0.72), transparent 13%),
    linear-gradient(165deg, #d9c8ff, #9b83ff);
  animation-delay: -0.75s;
  transform: translateY(20px) rotate(3deg);
}

.ear {
  position: absolute;
  top: -18px;
  width: 42px;
  height: 48px;
  border: 3px solid rgba(255, 255, 255, 0.8);
  border-radius: 55% 55% 46% 46%;
  background: inherit;
}

.ear-left {
  left: 25px;
  transform: rotate(-18deg);
}

.ear-right {
  right: 25px;
  transform: rotate(18deg);
}

.eye {
  position: absolute;
  top: 70px;
  width: 34px;
  height: 30px;
  overflow: hidden;
  border: 3px solid rgba(28, 35, 48, 0.08);
  border-radius: 50%;
  background: #fff;
  transition: height 0.2s ease, transform 0.2s ease;
}

.eye-left {
  left: 43px;
}

.eye-right {
  right: 43px;
}

.buddy-friend .eye {
  top: 64px;
  width: 30px;
  height: 27px;
}

.buddy-friend .eye-left {
  left: 38px;
}

.buddy-friend .eye-right {
  right: 38px;
}

.eye::after {
  position: absolute;
  top: calc(8px + var(--eye-y));
  left: calc(13px + var(--eye-x));
  width: 11px;
  height: 11px;
  border-radius: 50%;
  background: #172033;
  box-shadow: 3px -3px 0 0 rgba(255, 255, 255, 0.76) inset;
  content: "";
  transition: top 0.18s ease, left 0.18s ease;
}

.smile {
  position: absolute;
  left: 50%;
  top: 112px;
  width: 38px;
  height: 18px;
  border-bottom: 4px solid rgba(31, 35, 48, 0.72);
  border-radius: 0 0 999px 999px;
  transform: translateX(-50%);
  transition: width 0.2s ease, height 0.2s ease;
}

.cheek {
  position: absolute;
  top: 100px;
  width: 22px;
  height: 12px;
  border-radius: 50%;
  background: rgba(255, 127, 127, 0.36);
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.cheek-left {
  left: 31px;
}

.cheek-right {
  right: 31px;
}

.paw {
  position: absolute;
  top: 112px;
  width: 32px;
  height: 44px;
  border: 3px solid rgba(255, 255, 255, 0.65);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.28);
  transition: transform 0.2s ease;
}

.paw-left {
  left: -12px;
  transform: rotate(12deg);
}

.paw-right {
  right: -12px;
  transform: rotate(-12deg);
}

.is-typing .buddy-main .paw-left {
  transform: translate(8px, -10px) rotate(38deg);
}

.is-typing .buddy-friend .paw-right {
  transform: translate(-8px, -10px) rotate(-38deg);
}

.is-password .buddy-main .cheek,
.is-secret-visible .buddy .cheek {
  opacity: 0.9;
  transform: scale(1.12);
}

.is-secret-visible .buddy-main .eye {
  height: 8px;
  transform: translateY(10px);
}

.is-secret-visible .buddy-main .smile {
  width: 26px;
}

.is-secret-visible .buddy-friend {
  animation-duration: 2.2s;
}

@keyframes login-bob {
  0%,
  100% {
    translate: 0 0;
    rotate: -1deg;
  }

  50% {
    translate: 0 -10px;
    rotate: 1deg;
  }
}

@keyframes login-sparkle {
  0%,
  100% {
    opacity: 0.35;
    scale: 0.78;
  }

  50% {
    opacity: 1;
    scale: 1.1;
  }
}

@media (prefers-reduced-motion: reduce) {
  .buddy,
  .sparkle {
    animation: none;
  }
}

@media (max-width: 768px) {
  .cute-login-characters {
    min-height: 260px;
    padding-top: 24px;
  }

  .buddy {
    width: 132px;
    height: 148px;
  }

  .buddy-friend {
    width: 120px;
    height: 134px;
  }
}
</style>

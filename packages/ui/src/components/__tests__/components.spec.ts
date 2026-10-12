import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import YmEmptyState from '../YmEmptyState.vue'
import YmLogo from '../YmLogo.vue'
import YmStatCard from '../YmStatCard.vue'

describe('YmLogo', () => {
  it('渲染品牌标记与文字', () => {
    const wrapper = mount(YmLogo)

    expect(wrapper.find('.ym-logo__mark').text()).toBe('研')
    expect(wrapper.find('.ym-logo__text b').text()).toBe('研脉学术')
    expect(wrapper.find('.ym-logo__text i').text()).toBe('学术知识付费平台')
  })

  it('showText=false 时只渲染标记', () => {
    const wrapper = mount(YmLogo, { props: { showText: false } })
    expect(wrapper.find('.ym-logo__text').exists()).toBe(false)
  })

  it('subtitle 为空时不渲染副标题', () => {
    const wrapper = mount(YmLogo, { props: { subtitle: '' } })
    expect(wrapper.find('.ym-logo__text i').exists()).toBe(false)
  })

  it('size 决定标记尺寸', () => {
    const wrapper = mount(YmLogo, { props: { size: 48 } })
    const style = wrapper.find('.ym-logo__mark').attributes('style') ?? ''
    expect(style).toContain('width: 48px')
    expect(style).toContain('height: 48px')
  })
})

describe('YmStatCard', () => {
  it('展示标签与数值，并应用语义色调类', () => {
    const wrapper = mount(YmStatCard, { props: { label: '今日营收', value: 1280, tone: 'up' } })

    expect(wrapper.find('.ym-stat-card__label').text()).toBe('今日营收')
    expect(wrapper.find('.ym-stat-card__value').text()).toContain('1280')
    expect(wrapper.classes()).toContain('is-up')
  })

  it('loading 时以骨架屏替代数值', () => {
    const wrapper = mount(YmStatCard, { props: { label: '加载中', value: 0, loading: true } })

    expect(wrapper.find('.ym-stat-card__skeleton').exists()).toBe(true)
    expect(wrapper.find('.ym-stat-card__value').text()).not.toContain('0')
  })

  it('未传 hint 时不渲染提示行', () => {
    const wrapper = mount(YmStatCard, { props: { label: '注册用户', value: 3 } })
    expect(wrapper.find('.ym-stat-card__hint').exists()).toBe(false)
  })
})

describe('YmEmptyState', () => {
  it('渲染默认标题并支持自定义插槽', () => {
    const wrapper = mount(YmEmptyState, { slots: { default: '<button>新建</button>' } })

    expect(wrapper.find('.ym-empty__title').text()).toBe('暂无数据')
    expect(wrapper.find('button').text()).toBe('新建')
  })

  it('size=sm 时应用小尺寸类', () => {
    const wrapper = mount(YmEmptyState, { props: { size: 'sm' } })
    expect(wrapper.classes()).toContain('is-sm')
  })
})

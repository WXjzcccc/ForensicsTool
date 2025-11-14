import { ref, nextTick } from 'vue'

/**
 * 计算表格的动态高度
 * @param {Ref<HTMLElement>} resultCardRef - 结果卡片的引用
 * @param {Ref<string>} tableScrollHeight - 表格滚动高度的响应式引用
 * @param {number} offsetHeight - 额外偏移高度，默认为120px
 * @param {number} minHeight - 最小高度，默认为200px
 */
export function useTableHeight(resultCardRef, tableScrollHeight, offsetHeight = 120, minHeight = 200) {
  /**
   * 计算表格高度
   */
  const calculateTableHeight = () => {
    if (resultCardRef.value) {
      // 获取结果卡片的高度
      const cardHeight = resultCardRef.value.$el.offsetHeight
      // 计算表格高度（卡片高度减去其他元素的高度，如标题、按钮等）
      const calculatedHeight = cardHeight - offsetHeight
      
      // 设置最小高度
      const finalHeight = Math.max(calculatedHeight, minHeight)
      
      // 将像素值转换为vh单位
      const vhValue = (finalHeight / window.innerHeight) * 100
      tableScrollHeight.value = `${vhValue+4}vh`
    }
  }

  /**
   * 窗口大小变化处理函数
   */
  const handleResize = () => {
    nextTick(() => {
      calculateTableHeight()
    })
  }

  return {
    calculateTableHeight,
    handleResize
  }
}
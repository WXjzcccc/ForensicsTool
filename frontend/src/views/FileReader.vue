<template>
    <div class="page-container">
        <Card class="form-card">
            <template #content>
                <form class="form-layout">
                    <div class="field full-width">
                        <FloatLabel class="field full-width" variant="over">
                            <label for="selected">选择任务</label>
                            <Select id="selected" v-model="form.selected" :options="options" optionLabel="label"
                                optionValue="value" placeholder="请选择任务" v-tooltip.top="'选择要执行的数据提取任务类型'" />
                        </FloatLabel>
                    </div>

                    <div v-for="(fieldGroup, index) in inputFields" :key="index" class="input-row">
                        <div v-for="field in fieldGroup" :key="field.name" class="field half-width">
                            <FloatLabel variant="on">
                                <InputText v-if="field.type === 'input'" :id="field.name" v-model="form[field.name]"
                                    :placeholder="field.name === 'file' ? '请拖入文件或目录' : '解密密码'" aria-autocomplete="none"
                                    @drop.prevent="field.name === 'file' ? null : null"
                                    @dragover.prevent="field.name === 'file' ? null : null"
                                    v-tooltip.top="field.name === 'file' ? '拖入要提取数据的文件或目录路径' : '输入解密所需的密码（某些任务需要）'" />
                                <label :for="field.name">{{ field.label }}</label>
                            </FloatLabel>
                        </div>
                    </div>
                </form>
                <div class="button-row">
                    <Button label="解析" icon="pi pi-search" @click="handleExtract" class="button equal-width" />
                    <Button label="清空输出" icon="pi pi-trash" @click="handleClear" severity="secondary"
                        class="button equal-width" />
                </div>
            </template>
        </Card>
        <Card class="result-card" ref="resultCardRef" :loading="loading">
            <template #content>
                <div v-if="Object.keys(tableData).length === 0" class="empty">
                    <Empty />
                </div>
                <Tabs v-if="Object.keys(tableData).length !== 0" :value="getActiveTabIndex"
                    @tab-change="handleTabChange">
                    <TabList>
                        <Tab v-for="(tab, index) in tabs" :key="tab.label" :value="index">{{ tab.label }}</Tab>
                    </TabList>
                    <TabPanels>
                        <TabPanel v-for="(tab, index) in tabs" :key="tab.label" :value="index">
                            <DataTable :value="tableData[tab.value]" :columns="columns[tab.value]" dataKey="id"
                                stripedRows removableSort scrollable :scrollHeight="tableScrollHeight"
                                :virtualScrollerOptions="{ itemSize: 46 }" :reorderableColumns="true"
                                @column-reorder="onColumnReorder" class="p-datatable-sm">
                                <Column v-for="col in columns[tab.value]" :key="col.colKey" :field="col.colKey"
                                    :header="col.title" :sortable="true">
                                    <template #body="{ data, field }">
                                        <span @click="handleCellClick(data[field])">{{ data[field] }}</span>
                                    </template>
                                </Column>
                                <template #empty>
                                    <Empty />
                                </template>
                            </DataTable>
                        </TabPanel>
                    </TabPanels>
                </Tabs>
            </template>
        </Card>
        <Toast position="bottom-right" />
    </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useToast } from 'primevue/usetoast'
import { usePageDataStore } from "@/store"
import { watch } from "vue"
import { ClipboardSetText, OnFileDrop, OnFileDropOff } from "../../wailsjs/runtime/runtime.js"
import { ReadLevelDB, ReadMMKV } from "../../wailsjs/go/reader/FileReader.js"
import FloatLabel from 'primevue/floatlabel'
import Select from 'primevue/select'
import Tabs from 'primevue/tabs'
import Empty from '@/components/Empty.vue'
import { useTableHeight } from '@/composables/useTableHeight.js'

// 任务配置对象，定义每个任务需要的输入字段
const taskConfigs = {
    "1": {
        fields: [
            { name: "file", label: "文件/目录", type: "input" }
        ]
    },
    "2": {
        fields: [
            { name: "file", label: "文件/目录", type: "input" },
            { name: "password", label: "密码", type: "input" }
        ]
    }
}

const toast = useToast()
const store = usePageDataStore()
const form = ref(store.fileReaderStore?.formData || {
    selected: '',
    file: '',
    password: '',
})
// 计算当前任务配置
const currentTaskConfig = computed(() => {
    return taskConfigs[form.value.selected] || { fields: [] }
})

// 计算输入字段，每两个一组
const inputFields = computed(() => {
    const fields = currentTaskConfig.value.fields
    const groups = []
    for (let i = 0; i < fields.length; i += 2) {
        groups.push(fields.slice(i, i + 2))
    }
    return groups
})

const tableData = ref(store.fileReaderStore?.tableData || {})
const tabs = ref(store.fileReaderStore?.tabsData || [])
const activeTab = ref(store.fileReaderStore?.tabData || '')
const columns = ref(store.fileReaderStore?.columnsData || {})

// 动态计算DataTable的scrollHeight
const tableScrollHeight = ref('60vh')
const resultCardRef = ref(null)

// 使用公共的表格高度计算函数
const { calculateTableHeight, handleResize } = useTableHeight(resultCardRef, tableScrollHeight)

// 组件挂载后设置监听
onMounted(() => {
    window.addEventListener('resize', handleResize)
    // 初始计算
    nextTick(() => {
        calculateTableHeight()
    })
})

// 组件卸载时移除监听
onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    // 清理文件拖放监听器
    OnFileDropOff()
})

// 计算属性，用于获取当前活动标签的索引
const getActiveTabIndex = computed(() => {
    const index = tabs.value.findIndex(tab => tab.value === activeTab.value)
    return index >= 0 ? index : 0
})

watch([form, tableData, tabs, activeTab], () => {
    store.saveFileReaderData({
        formData: form.value,
        tableData: tableData.value,
        tabsData: tabs.value,
        tabData: activeTab.value,
        columnsData: columns.value
    })
})

const options = ref([
    { label: "leveldb读取，指定文件夹为leveldb数据库文件夹", value: "1" },
    { label: "mmkv读取，指定文件mmkv数据文件或mmkv目录，单文件时支持输入密码解密加密数据", value: "2" },
])

const loading = ref(false)

const handleExtract = () => {
    loading.value = true

    const currentConfig = taskConfigs[form.value.selected]
    if (!currentConfig) {
        toast.add({ severity: 'error', summary: '错误', detail: '请选择一个任务', life: 3000 })
        loading.value = false
        return
    }

    // 动态构建参数对象
    const params = {}
    currentConfig.fields.forEach(field => {
        params[field.name] = form.value[field.name]
    })

    // 从动态参数中获取file和password
    const { file, password } = params

    if (file === "" || file === undefined || file === null) {
        toast.add({ severity: 'error', summary: '错误', detail: '文件参数异常！', life: 3000 })
        loading.value = false
        return
    }

    let func = null
    switch (form.value.selected) {
        case "1":
            func = ReadLevelDB
            break
        case "2":
            func = ReadMMKV
            break
    }

    if (form.value.selected === "2") {
        if (password === undefined || password === null) {
            password = ""
        }
        func(file, password).then((result) => {
            console.log(result)
            try {
                if (result.err !== "") {
                    toast.add({ severity: 'error', summary: '错误', detail: result.err, life: 3000 })
                } else {
                    handleResult(result)
                }
            } finally {
                loading.value = false
            }
        })
    } else {
        func(file).then((result) => {
            try {
                if (result.err !== "") {
                    toast.add({ severity: 'error', summary: '错误', detail: result.err, life: 3000 })
                } else {
                    handleResult(result)
                }
            } finally {
                loading.value = false
            }
        })
    }
}

function handleResult(result) {
    columns.value = {}
    tabs.value = []
    tableData.value = {}
    let idx = 0
    for (let tab in result.data) {
        if (idx === 0) {
            activeTab.value = tab
            idx++
        }
        tabs.value.push({
            value: tab, label: tab
        })
        if (result.data[tab] == null) {
            toast.add({ severity: 'warn', summary: '提示', detail: tab + '无数据', life: 3000 })
            continue
        }
        let inf = result.data[tab][0]
        tableData.value[tab] = {}
        columns.value[tab] = []
        tableData.value[tab] = result.data[tab]
        for (let key in inf) {
            columns.value[tab].push({
                colKey: key, title: key, ellipsis: true, width: "15vw"
            })
        }
    }

    // 数据加载完成后重新计算表格高度
    nextTick(() => {
        calculateTableHeight()
    })
}

const handleClear = () => {
    tableData.value = {}
    tabs.value = []
    activeTab.value = ""
    columns.value = {}
}

const handleTabChange = (event) => {
    activeTab.value = tabs.value[event.index].value
}

const handleCellClick = (value) => {
    // 复制单元格内容
    ClipboardSetText(value).then((ok) => {
        if (ok) {
            toast.add({ severity: 'success', summary: '成功', detail: '已复制到剪贴板', life: 2000 })
        }
    })
}

onMounted(() => {
    OnFileDrop((x, y, paths) => {
        if (paths.length > 0) {
            form.value.file = paths[0]
            toast.add({ severity: 'success', summary: '成功', detail: '文件已添加', life: 3000 })
        }
    }, false)
})

const onColumnReorder = (event) => {
    columns.value[activeTab.value] = event.columns
}
</script>

<style scoped></style>
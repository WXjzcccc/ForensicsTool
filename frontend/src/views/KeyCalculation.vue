<template>
    <div class="page-container">
        <Card class="form-card">
            <template #content>
                <form class="form-layout">
                    <div class="field full-width">
                        <FloatLabel class="field full-width" variant="over">
                            <label for="task">选择任务</label>
                            <Select
                                id="task"
                                v-model="form.selected"
                                :options="options"
                                optionLabel="label"
                                optionValue="value"
                                placeholder="请选择任务"
                                v-tooltip.top="'选择要计算的密钥类型'"
                            />
                        </FloatLabel>
                    </div>

                    <!-- 动态生成的输入字段 -->
                    <div
                        v-for="(row, rowIndex) in inputFields"
                        :key="rowIndex"
                        class="input-row"
                    >
                        <div
                            v-for="field in row"
                            :key="field.name"
                            class="field"
                            :class="
                                row.length === 1 ? 'full-width' : 'half-width'
                            "
                        >
                            <!-- 根据字段类型渲染不同的组件 -->
                            <div
                                v-if="field.type === 'checkbox'"
                                class="flex align-items-center mt-3"
                            >
                                <Checkbox
                                    :id="field.name"
                                    v-model="form[field.name]"
                                    :binary="true"
                                />
                                <label :for="field.name" class="ml-2">{{
                                    field.label
                                }}</label>
                            </div>

                            <FloatLabel v-else variant="on">
                                <InputText
                                    v-if="field.type === 'input'"
                                    :id="field.name"
                                    v-model="form[field.name]"
                                    aria-autocomplete="none"
                                    :placeholder="calPlace[field.name]"
                                    v-tooltip.top="calPlace[field.name]"
                                />

                                <!-- 下拉选择框 -->
                                <Dropdown
                                    v-else-if="field.type === 'dropdown'"
                                    :id="field.name"
                                    v-model="form[field.name]"
                                    :options="field.options"
                                    optionLabel="label"
                                    optionValue="value"
                                    class="w-full"
                                />

                                <!-- 数字输入框 -->
                                <InputNumber
                                    v-else-if="field.type === 'number'"
                                    :id="field.name"
                                    v-model="form[field.name]"
                                    class="w-full"
                                />

                                <!-- 日期选择器 -->
                                <Calendar
                                    v-else-if="field.type === 'calendar'"
                                    :id="field.name"
                                    v-model="form[field.name]"
                                    class="w-full"
                                />

                                <label :for="field.name">{{
                                    field.label
                                }}</label>
                            </FloatLabel>
                        </div>
                    </div>
                </form>

                <div class="button-row">
                    <Button class="button equal-width" @click="handleCalculate">
                        <i class="pi pi-calculator"></i>
                        计算
                    </Button>
                    <Button
                        class="button equal-width"
                        severity="secondary"
                        @click="handleClear"
                    >
                        <i class="pi pi-trash"></i>
                        清空输出
                    </Button>
                </div>
            </template>
        </Card>
        <Card class="result-card">
            <template #content>
                <div v-if="resultText === ''" class="empty">
                    <Empty />
                </div>
                <div class="result-output" v-html="resultText" />
            </template>
        </Card>
    </div>
</template>

<script setup>
import { ref, computed, watch } from "vue";
import {
    CalMostone,
    CalTiktok,
    CalWechat,
    CalWechatIndex,
    CalWildFire,
    CalBatChat,
} from "../../wailsjs/go/passwdCalc/PasswdCalc.js";
import {
    generateNormalTextOutput,
    generateSuccessTextOutput,
} from "@/utils.js";
import { usePageDataStore } from "@/store";
import FloatLabel from "primevue/floatlabel";
import Select from "primevue/select";
import InputText from "primevue/inputtext";
import InputNumber from "primevue/inputnumber";
import Checkbox from "primevue/checkbox";
import Calendar from "primevue/calendar";
import Button from "primevue/button";
import Dropdown from "primevue/dropdown";
import Empty from "@/components/Empty.vue";

// 任务配置对象，定义每个任务需要的输入字段
const taskConfigs = {
    1: {
        // 微信的EnMicroMsg.db
        fields: [
            { name: "uin", label: "微信uin", type: "input" },
            { name: "imei", label: "imei", type: "input" },
        ],
    },
    2: {
        // 微信的FTS5IndexMicroMsg_encrypt.db
        fields: [
            { name: "uin", label: "微信uin", type: "input" },
            { name: "wxid", label: "wxid", type: "input" },
            { name: "imei", label: "imei", type: "input" },
        ],
    },
    3: {
        // 野火IM系应用的data
        fields: [{ name: "token", label: "token", type: "input" }],
    },
    4: {
        // 默往APP的msg.db
        fields: [{ name: "uid", label: "uid", type: "input" }],
    },
    5: {
        // 抖音的聊天数据库
        fields: [{ name: "uid", label: "uid", type: "input" }],
    },
    6: {
        fields: [{ name: "uid", label: "uid", type: "input" }],
    },
};
const calPlace = {
    uin: "微信用户的uin，可能是负值，在shared_prefs/auth_info_key_prefs.xml文件中_auth_uin的值",
    imei: "微信获取到的IMEI或MEID，在shared_prefs/DENGTA_META.xml文件中IMEI_DENGTA的值，在高版本中通常是1234567890ABCDEF，可以为空",
    wxid: "数据库所属的wxid，一般情况下在解密EnMicroMsg.db的时候会一并提取，若无需要，请从shared_prefs/com.tencent.mm_preferences.xml中提取login_weixin_username的值",
    token: "野火IM系应用的用户token，shared_prefs/config.xml的token的值",
    uid: "默往（通常在shared_prefs/im.xml中的userId的值）、抖音（数据库文件名中的id）计算密钥需要的内容、QQ（msf_mmkv_file中QQ号对应的uid）、蝙蝠（数据库名中的数字）",
};
const store = usePageDataStore();
const form = ref(
    store.keyCalculationStore?.formData || {
        selected: "",
        uin: "",
        imei: "",
        wxid: "",
        token: "",
        uid: "",
    },
);
const resultText = ref(store.keyCalculationStore?.resultData || "");
const options = ref([
    { label: "微信的EnMicroMsg.db", value: "1" },
    { label: "微信的FTS5IndexMicroMsg_encrypt.db", value: "2" },
    { label: "野火IM系应用的data", value: "3" },
    { label: "默往APP的msg.db", value: "4" },
    { label: "抖音的聊天数据库", value: "5" },
    { label: "蝙蝠的聊天数据库", value: "6" },
]);

// 根据选择的任务获取当前任务的配置
const currentTaskConfig = computed(() => {
    return taskConfigs[form.value.selected] || { fields: [] };
});

// 动态生成输入字段数组，每行最多两个字段
const inputFields = computed(() => {
    const fields = currentTaskConfig.value.fields || [];
    const result = [];

    for (let i = 0; i < fields.length; i += 2) {
        result.push(fields.slice(i, i + 2));
    }

    return result;
});

watch([form, resultText], () => {
    store.saveKeyCalculationData({
        formData: form.value,
        resultData: resultText.value,
    });
});

const handleCalculate = () => {
    // 从form中获取当前任务需要的字段值
    const currentFields = currentTaskConfig.value.fields;
    const params = {};

    // 只获取当前任务需要的字段值
    currentFields.forEach((field) => {
        params[field.name] = form.value[field.name];
    });

    switch (form.value.selected) {
        case "1": {
            CalWechat(params.uin, params.imei).then((result) => {
                resultText.value += generateSuccessTextOutput(
                    "成功获取到安卓微信EnMicroMsg.db数据库密钥",
                    result,
                );
            });
            break;
        }
        case "2": {
            CalWechatIndex(params.uin, params.wxid, params.imei).then(
                (result) => {
                    if (result !== "") {
                        resultText.value += generateSuccessTextOutput(
                            "成功获取到安卓微信FTS5IndexMicroMsg_encrypt.db数据库密钥",
                            result,
                        );
                    } else {
                        resultText.value += generateNormalTextOutput(
                            "计算失败，请确认参数",
                            "red",
                        );
                    }
                },
            );
            break;
        }
        case "3": {
            CalWildFire(params.token).then((result) => {
                resultText.value += generateSuccessTextOutput(
                    "成功获取到野火IM数据库密钥",
                    result[0],
                );
                resultText.value += generateNormalTextOutput(
                    result[1],
                    "green",
                );
            });
            break;
        }
        case "4": {
            CalMostone(params.uid).then((result) => {
                resultText.value += generateSuccessTextOutput(
                    "成功获取到默往msg.db数据库密钥",
                    result[0],
                );
                resultText.value += generateNormalTextOutput(
                    result[1],
                    "green",
                );
            });
            break;
        }
        case "5": {
            CalTiktok(params.uid).then((result) => {
                resultText.value += generateSuccessTextOutput(
                    "成功获取到抖音聊天数据库密钥",
                    result[0],
                );
                resultText.value += generateNormalTextOutput(
                    result[1],
                    "green",
                );
            });
            break;
        }
        case "6": {
            CalBatChat(params.uid).then((result) => {
                resultText.value += generateSuccessTextOutput(
                    "成功获取到蝙蝠聊天数据库密钥",
                    result[0],
                );
                resultText.value += generateNormalTextOutput(
                    result[1],
                    "green",
                );
            });
            break;
        }
    }
};

const handleClear = () => {
    resultText.value = "";
};

watch(resultText, () => {
    const card = document.querySelector(".result-card");
    if (card) {
        void card.offsetHeight;
        card.scrollTop = card.scrollHeight;
    }
});
</script>
<style scoped></style>

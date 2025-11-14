<script setup>
import { useLayout } from "../composables/useLayout";

const { primaryColors, surfaces, primary, surface, isDarkMode, updateColors } = useLayout();
</script>

<template>
    <div class="config-panel hidden">
        <div class="config-section">
            <div>
                <span class="config-label">主题</span>
                <div class="config-colors">
                    <button
                        v-for="pc of primaryColors"
                        :key="pc.name"
                        type="button"
                        :title="pc.name"
                        :class="['color-button', { selected: primary === pc.name }]"
                        :style="{ backgroundColor: pc.palette['500'] }"
                        @click="updateColors('primary', pc.name)"
                    />
                </div>
            </div>
            <div>
                <span class="config-label">底色</span>
                <div class="config-colors">
                    <button
                        v-for="s of surfaces"
                        :key="s.name"
                        type="button"
                        :title="s.name"
                        :class="[
                            'color-button',
                            {
                                selected: surface
                                    ? surface === s.name
                                    : isDarkMode
                                    ? s.name === 'zinc'
                                    : s.name === 'slate',
                            },
                        ]"
                        :style="{ backgroundColor: s.palette['500'] }"
                        @click="updateColors('surface', s.name)"
                    />
                </div>
            </div>
        </div>
    </div>
</template>

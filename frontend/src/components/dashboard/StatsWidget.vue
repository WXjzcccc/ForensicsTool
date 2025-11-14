<script setup>
const stats = [
    {
        title: "Total Orders",
        icon: "pi-shopping-cart",
        value: "1,234",
        subtitle: "Last 7 days",
    },
    {
        title: "Active Users",
        icon: "pi-users",
        value: "2,573",
        subtitle: "Last 7 days",
    },
    {
        title: "Revenue",
        icon: "pi-dollar",
        value: "$45,200",
        subtitle: "Last 7 days",
    },
    {
        title: "Success Rate",
        icon: "pi-chart-line",
        value: "95%",
        subtitle: "Last 7 days",
    },
];
</script>

<template>
    <div class="stats">
        <div v-for="(stat, index) in stats" :key="index" class="layout-card">
            <div class="stats-header">
                <span class="stats-title">{{ stat.title }}</span>
                <span class="stats-icon-box">
                    <i :class="['pi', stat.icon]"></i>
                </span>
            </div>
            <div class="stats-content">
                <div class="stats-value">{{ stat.value }}</div>
                <div class="stats-subtitle">{{ stat.subtitle }}</div>
            </div>
        </div>
    </div>
</template>

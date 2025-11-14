<script setup>
import { ref, watch, onMounted } from "vue";
import { useLayout } from "../../composables/useLayout";

const { primary, surface, isDarkMode } = useLayout();

const chartData = ref(null);
const chartOptions = ref(null);

function setChartData() {
    const documentStyle = getComputedStyle(document.documentElement);

    return {
        labels: ["Q1", "Q2", "Q3", "Q4"],
        datasets: [
            {
                type: "bar",
                label: "Subscriptions",
                backgroundColor: documentStyle.getPropertyValue("--p-primary-400"),
                data: [4000, 10000, 15000, 4000],
                barThickness: 32,
            },
            {
                type: "bar",
                label: "Advertising",
                backgroundColor: documentStyle.getPropertyValue("--p-primary-300"),
                data: [2100, 8400, 2400, 7500],
                barThickness: 32,
            },
            {
                type: "bar",
                label: "Affiliate",
                backgroundColor: documentStyle.getPropertyValue("--p-primary-200"),
                data: [4100, 5200, 3400, 7400],
                borderRadius: {
                    topLeft: 8,
                    topRight: 8,
                },
                barThickness: 32,
            },
        ],
    };
}

function setChartOptions() {
    return {
        maintainAspectRatio: false,
        responsive: true,
        plugins: {
            legend: {
                position: "top",
            },
        },
        scales: {
            x: {
                stacked: true,
                grid: {
                    color: "transparent",
                    borderColor: "transparent",
                },
            },
            y: {
                stacked: true,
                grid: {
                    color: "transparent",
                    borderColor: "transparent",
                    drawTicks: false,
                },
            },
        },
    };
}

watch([primary, surface, isDarkMode], () => {
    chartData.value = setChartData();
    chartOptions.value = setChartOptions();
});

onMounted(() => {
    chartData.value = setChartData();
    chartOptions.value = setChartOptions();
});
</script>

<template>
    <div class="layout-card col-item-2">
        <div class="chart-header">
            <span class="chart-title">Sales Trend</span>
        </div>
        <div class="chart-content">
            <Chart type="bar" :data="chartData" :options="chartOptions" style="height: 300px" />
        </div>
    </div>
</template>

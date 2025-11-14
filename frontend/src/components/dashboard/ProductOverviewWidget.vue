<script setup>
import { ref, watch, onMounted } from "vue";

const products = ref([
    {
        name: "Laptop Pro",
        category: "Electronics",
        price: 2499,
        status: "In Stock",
    },
    {
        name: "Wireless Mouse",
        category: "Accessories",
        price: 49,
        status: "Low Stock",
    },
    {
        name: "Monitor 4K",
        category: "Electronics",
        price: 699,
        status: "Out of Stock",
    },
    { name: "Keyboard", category: "Accessories", price: 149, status: "In Stock" },
]);

const selectedProduct = ref(null);
const searchQuery = ref("");
const loading = ref(false);
const filteredProducts = ref([]);

const searchProducts = () => {
    loading.value = true;
    filteredProducts.value = products.value.filter(
        (product) =>
            product.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            product.category.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            product.status.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
    setTimeout(() => {
        loading.value = false;
    }, 300);
};

watch(searchQuery, () => {
    searchProducts();
});

onMounted(() => {
    filteredProducts.value = [...products.value];
});
</script>

<template>
    <div class="layout-card">
        <div class="products-header">
            <span class="products-title">Products Overview</span>
            <IconField class="search-field">
                <InputIcon class="pi pi-search" />
                <InputText
                    v-model="searchQuery"
                    placeholder="Search products..."
                    class="products-search"
                    @keyup.enter="searchProducts"
                />
            </IconField>
        </div>
        <div class="products-table-container">
            <DataTable
                :value="filteredProducts"
                v-model:selection="selectedProduct"
                selectionMode="single"
                :loading="loading"
                :rows="5"
                class="products-table"
                :pt="{
                    mask: {
                        class: 'products-table-mask',
                    },
                    loadingIcon: {
                        class: 'products-table-loading',
                    },
                }"
            >
                <Column field="name" header="Name" sortable></Column>
                <Column field="category" header="Category" sortable></Column>
                <Column field="price" header="Price" sortable>
                    <template #body="{ data }"> ${{ data.price }} </template>
                </Column>
                <Column field="status" header="Status">
                    <template #body="{ data }">
                        <Tag
                            :severity="
                                data.status === 'In Stock' ? 'success' : data.status === 'Low Stock' ? 'warn' : 'danger'
                            "
                        >
                            {{ data.status }}
                        </Tag>
                    </template>
                </Column>
            </DataTable>
        </div>
    </div>
</template>

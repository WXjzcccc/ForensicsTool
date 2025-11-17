import { createRouter, createWebHistory } from 'vue-router';

// @ts-ignore
import KeyCalculation from '../views/KeyCalculation.vue';
// @ts-ignore
import DatabaseDecrypt from '../views/DatabaseDecrypt.vue';
// @ts-ignore
import DataExtraction from '../views/DataExtraction.vue';
// @ts-ignore
import RegistryAnalysis from '../views/RegistryAnalysis.vue';
// @ts-ignore
import TimestampParser from '../views/TimestampParser.vue';
// @ts-ignore
import BruteForce from '../views/BruteForce.vue';
// @ts-ignore
import FileReader from '../views/FileReader.vue';
// @ts-ignore
import About from '../views/About.vue';

const routes = [
    {
        path: '/',
        redirect: '/KeyCalculation'
    },
    {
        path: '/KeyCalculation',
        name: 'KeyCalculation',
        component: KeyCalculation
    },
    {
        path: '/DatabaseDecrypt',
        name: 'DatabaseDecrypt',
        component: DatabaseDecrypt
    },
    {
        path: '/DataExtraction',
        name: 'DataExtraction',
        component: DataExtraction
    },
    {
        path: '/RegistryAnalysis',
        name: 'RegistryAnalysis',
        component: RegistryAnalysis
    },
    {
        path: '/TimestampParser',
        name: 'TimestampParser',
        component: TimestampParser
    },
    {
        path: '/BruteForce',
        name: 'BruteForce',
        component: BruteForce
    },
    {
        path: '/FileReader',
        name: 'FileReader',
        component: FileReader
    },
    {
        path: '/About',
        name: 'About',
        component: About
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

export default router;
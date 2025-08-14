export default [
    {
        path: '/',
        redirect: '/key-calculation'
    },
    {
        path: '/key-calculation',
        name: 'KeyCalculation',
        component: () => import('@/views/KeyCalculation.vue')
    },
    {
        path: '/database-decrypt',
        name: 'DatabaseDecrypt',
        component: () => import('@/views/DatabaseDecrypt.vue')
    },
    {
        path: '/data-extraction',
        name: 'DataExtraction',
        component: () => import('@/views/DataExtraction.vue')
    },
    {
        path: '/registry-analysis',
        name: 'RegistryAnalysis',
        component: () => import('@/views/RegistryAnalysis.vue')
    },
    {
        path: '/brute-force',
        name: 'BruteForce',
        component: () => import('@/views/BruteForce.vue')
    },
    {
        path: '/about',
        name: 'About',
        component: () => import('@/views/About.vue')
    }
  ];

import "./assets/styles/main.css";

import { createApp } from "vue";
import PrimeVue from "primevue/config";
import Aura from "@primeuix/themes/aura";
import router from "./router";
// @ts-ignore
import App from "./App.vue";
import ToastService from 'primevue/toastservice';
import Tooltip from 'primevue/tooltip';
import store from "./store";

const app = createApp(App);

app.use(PrimeVue, {
    theme: {
        preset: Aura,
        options: {
            darkModeSelector: ".p-dark",
        },
    },
});
app.use(ToastService);
app.use(store);
app.use(router);

app.directive('tooltip', Tooltip);

app.mount("#app");

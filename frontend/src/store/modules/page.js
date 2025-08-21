import { defineStore } from 'pinia';

export const usePageDataStore = defineStore('pageData', {
  state: () => ({
    keyCalculationStore: null,
    databaseDecryptStore: null,
    dataExtractionStore: null,
    registryStore: null,
    bruteForceStore: null,
    timestampStore: null,
  }),
  actions: {
    saveKeyCalculationData(data) {
      this.keyCalculationStore = data
    },
    saveDatabaseDecryptData(data) {
      this.databaseDecryptStore = data
    },
    saveDataExtractionData(data) {
      this.dataExtractionStore = data
    },
    saveRegistryData(data) {
      this.registryStore = data
    },
    saveBruteForceData(data) {
      this.bruteForceStore = data
    },
    saveTimestampData(data) {
      this.timestampStore = data
    }
  },
});

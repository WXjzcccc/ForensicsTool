import { defineStore } from 'pinia';

interface PageDataState {
  keyCalculationStore: any | null;
  databaseDecryptStore: any | null;
  dataExtractionStore: any | null;
  registryStore: any | null;
  bruteForceStore: any | null;
  timestampStore: any | null;
}

export const usePageDataStore = defineStore('pageData', {
  state: (): PageDataState => ({
    keyCalculationStore: null,
    databaseDecryptStore: null,
    dataExtractionStore: null,
    registryStore: null,
    bruteForceStore: null,
    timestampStore: null,
  }),
  actions: {
    saveKeyCalculationData(data: any) {
      this.keyCalculationStore = data;
    },
    saveDatabaseDecryptData(data: any) {
      this.databaseDecryptStore = data;
    },
    saveDataExtractionData(data: any) {
      this.dataExtractionStore = data;
    },
    saveRegistryData(data: any) {
      this.registryStore = data;
    },
    saveBruteForceData(data: any) {
      this.bruteForceStore = data;
    },
    saveTimestampData(data: any) {
      this.timestampStore = data;
    },
    // 获取爆破状态
    getBruteForceCrackingState(): boolean {
      return this.bruteForceStore?.cracking || false;
    },
    // 更新爆破状态
    updateBruteForceCrackingState(cracking: boolean) {
      if (!this.bruteForceStore) {
        this.bruteForceStore = {};
      }
      this.bruteForceStore.cracking = cracking;
    }
  },
});
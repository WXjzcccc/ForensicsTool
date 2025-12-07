import { defineStore } from 'pinia';

interface PageDataState {
  keyCalculationStore: any | null;
  databaseDecryptStore: any | null;
  dataExtractionStore: any | null;
  registryStore: any | null;
  bruteForceStore: any | null;
  timestampStore: any | null;
  fileReaderStore: any | null;
  ipLocationStore: any | null;
  isCardCollapsed: boolean | null;
}

export const usePageDataStore = defineStore('pageData', {
  state: (): PageDataState => ({
    keyCalculationStore: null,
    databaseDecryptStore: null,
    dataExtractionStore: null,
    registryStore: null,
    bruteForceStore: null,
    timestampStore: null,
    fileReaderStore: null,
    ipLocationStore: null,
    isCardCollapsed: null,
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
    saveFileReaderData(data: any) {
      this.fileReaderStore = data;
    },
    saveIPLocationData(data: any) {
      this.ipLocationStore = data;
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
    },
    // 设置卡片折叠状态
    setCardCollapsed(collapsed: boolean) {
      this.isCardCollapsed = collapsed;
    },
    // 获取卡片折叠状态
    getCardCollapsed(): boolean {
      return this.isCardCollapsed || false;
    }
  },
});
export const state = (_) => ({
  inx: '',
  isload: null, //判断加载状态，初次进入和刷新合一
});
export const mutations = {
  setInx: (state, data) => (state.inx = data),
  setIsload: (state, isload) => (state.isload = isload),
};

export const actions = {
  async nuxtServerInit({ state }, { req }) {
    commit('setIsload', true);
  },
};

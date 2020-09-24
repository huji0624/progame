let stores = null;
export default function ({ store, $axios, app, redirect }, inject) {
  const axios = $axios;
  if (!stores) stores = store;
  axios.defaults.baseURL = 'http://localhost:8888/';
  axios.defaults.timeout = 10000;

  axios.onResponse(
    (res) => {
      return res.data;
    },
    (err) => {
      console.log('遇到了错误：' + err);
    }
  );
}

export { stores };

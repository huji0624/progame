let stores = null;
export default function ({ store, $axios, app, redirect }, inject) {
  const axios = $axios;
  if (!stores) stores = store;
  const url =
    process.env.NODE_ENV === 'production'
      ? 'https://pgame.51wnl-cq.com/'
      : 'localhost:8881/';
  axios.defaults.baseURL = url;
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

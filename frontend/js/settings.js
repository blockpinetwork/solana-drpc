
const request = window.axios.create({
  baseURL: window.axiosBaseUrl,
  timeout: 60000
});

request.interceptors.response.use(
  response => {
    const { data } = response;
    if (typeof data === 'string') {
      return data;
    }
    if (data.code === 'SUCCESS') {
      return data;
    } else {
      Message.error(data.message);
      return Promise.reject(data.message);
    }
  },
  error => {
    if (error.message) {
      Message({
        message: error.message,
        type: 'error',
        showClose: true
      });
    }
    return Promise.reject(error);
  }
);
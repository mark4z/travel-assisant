import type {AxiosInstance} from 'axios'
import axios, {AxiosError} from 'axios';
import {ElMessage} from 'element-plus'

const axiosInstance: AxiosInstance = axios.create({
    baseURL: '',
});

export default axiosInstance;

axiosInstance.interceptors.response.use(
    onFulfilled => {
        return onFulfilled
    },
    error => {
        console.log(error.response.data)
        return Promise.reject(error.response.data.error)
    }
)



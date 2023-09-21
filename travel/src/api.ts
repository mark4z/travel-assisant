import axiosInstance from './axios-instance';

export async function get<T>(url: string, params?: any): Promise<T> {
    const response = await axiosInstance.get<T>(url, {params});
    return response.data;
}

export async function post<T>(url: string, data?: any): Promise<T> {
    const response = await axiosInstance.post<T>(url, data);
    return response.data;
}


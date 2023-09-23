import axiosInstance from './axios-instance';

export async function get<T>(url: string, params?: any): Promise<T> {
    const response = await axiosInstance.get<T>(url, {params});
    return response.data;
}

export async function post<T>(url: string, data?: any): Promise<T> {
    const response = await axiosInstance.post<T>(url, data);
    return response.data;
}

export interface Train {
    train_no: string
    train_code: string
    start_time: string
    end_time: string
    start_station: string
    start_station_name: string
    end_station: string
    end_station_name: string
    from_station: string
    from_station_name: string
    to_station: string
    to_station_name: string
    two_seat: string
    one_seat: string
    special_seat: string
    children?: Train[]
    hasChildren?: boolean
}

export interface Stations {
    value: string
    label: string
}

export interface Pass {
    station: string
    station_name: string
    arrive_time: string
    start_time: string

    two_seat: string
    one_seat: string
    special_seat: string
}

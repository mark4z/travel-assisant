import axiosInstance from './axios-instance';

const mapperUrl = "/otn/resources/js/framework/station_name.js?station_version=1.9274"

let m = <Record<string, string>>{}

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

// 定义缓存过期时间（单位：毫秒）
const cacheExpirationTime = 24 * 60 * 60 * 1000; // 24小时


export async function init(): Promise<Stations[]> {
    m = await mapper()
    //convert to Stations[]
    let res = <Stations[]>[]
    for (let k in m) {
        // if k start with A-Z
        if (/^[A-Z]/.test(k)) {
            res.push({
                value: k,
                label: m[k]
            })
        }
    }
    return res
}

// 保存映射到LocalStorage中
async function mapper(): Promise<Record<string, string>> {
    const localStorageKey = 'mapperData';
    const storedData = localStorage.getItem(localStorageKey);

    // 检查是否有缓存数据
    if (storedData) {
        const {data, timestamp} = JSON.parse(storedData);

        // 检查数据是否仍然有效
        if (Date.now() - timestamp <= cacheExpirationTime) {
            return data;
        }
    }
    try {
        const response = await get<string>(mapperUrl);

        const city = response.split('@').slice(1);
        const temp: Record<string, string> = {};

        for (const c of city) {
            const fields = c.split('|');
            temp[fields[1]] = fields[2];
            temp[fields[2]] = fields[1];
        }
        // 将数据存储在LocalStorage中，并记录时间戳
        const dataToStore = {data: temp, timestamp: Date.now()};
        localStorage.setItem(localStorageKey, JSON.stringify(dataToStore));
        // 清理过期的缓存数据
        cleanupCache();
        return temp;
    } catch (error) {
        // 处理HTTP请求错误
        throw error;
    }
}

// 清理过期的缓存数据
function cleanupCache() {
    const localStorageKey = 'mapperData';
    const storedData = localStorage.getItem(localStorageKey);

    if (storedData) {
        const {timestamp} = JSON.parse(storedData);

        // 检查缓存数据是否过期，如果过期则删除
        if (Date.now() - timestamp > cacheExpirationTime) {
            localStorage.removeItem(localStorageKey);
        }
    }
}

import axiosInstance from './axios-instance';

let initialized = false;
const mapperUrl = "/otn/resources/js/framework/station_name.js?station_version=1.9274"
const trainUrl = '/otn/leftTicket/queryZ';
const indexUrl = "/otn/leftTicket/init"

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
        localStorage.removeItem(localStorageKey);
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
        return temp;
    } catch (error) {
        // 处理HTTP请求错误
        throw error;
    }
}

export async function originalSearch(from: string, to: string, date: string): Promise<Train[]> {
    return findAllTrain(from, to, date)
}

async function findAllTrain(from: string, to: string, date: string): Promise<Train[]> {
    await index(from, to, date)
    return get<any>(`${trainUrl}`, {
        'leftTicketDTO.train_date': date,
        'leftTicketDTO.from_station': from,
        'leftTicketDTO.to_station': to,
        purpose_codes: 'ADULT',
    })
        .then(response => {
            const res: Train[] = [];
            for (const t of response.data.result) {
                const trainRes = decode(t);
                if (trainRes.train_no.startsWith('G') || trainRes.train_no.startsWith('D') || trainRes.train_no.startsWith('C')) {
                    res.push(trainRes);
                }
            }
            if (res.length === 0) {
                throw new Error(`Cannot find target ${from}-${to} ${date}`);
            }
            return res;
        })
        .catch((error) => {
            throw error;
        });
}

async function index(from: string, to: string, date: string): Promise<void> {
    if (initialized) {
        return;
    }
    const req = {
        linktypeid: 'dc',
        fs: m[from] + ',' + from,
        ts: m[to] + ',' + to,
        date: date,
        flag: 'N,N,Y'
    }
    get(indexUrl, req)
        .then(response => {
            console.info("index success")
        })
    initialized = true;
}

function decode(info: string): Train {
    const fields = info.split('|');
    return {
        train_code: fields[2],
        train_no: fields[3],
        special_seat: fields[32],
        one_seat: fields[31],
        two_seat: fields[30],
        start_station: fields[4],
        start_station_name: m[fields[4]], // 请确保 m 对象在当前作用域中可用
        end_station: fields[5],
        end_station_name: m[fields[5]], // 请确保 m 对象在当前作用域中可用
        from_station: fields[6],
        from_station_name: m[fields[6]], // 请确保 m 对象在当前作用域中可用
        to_station: fields[7],
        to_station_name: m[fields[7]], // 请确保 m 对象在当前作用域中可用
        start_time: fields[8],
        end_time: fields[9],
    };
}
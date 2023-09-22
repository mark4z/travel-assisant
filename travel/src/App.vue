<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {get} from '@/api';

const date = ref(new Date('2023-10-01'))
const from = ref('UUH')
const to = ref('TNV')
const no = ref('G3133')
const trains = ref<Train[]>([])

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

const tableRowClassName = ({
                             row,
                             rowIndex,
                           }: {
  row: Train
  rowIndex: number
}) => {
  if (row.two_seat === '有') {
    return 'success-row'
  } else {
    return 'warning-row'
  }
}

const options = ref<Stations[]>([])

onMounted(() => {
  get('/api/stations').then(res => {
    options.value = (res as Stations[])
  })
})

function walk() {
  get("/api/walk", {
    from: from.value,
    to: to.value,
    // YYYY-MM-DD
    date: date.value.toISOString().slice(0, 10),
    no: no.value
  }).then(res => {
    trains.value.push((res as Train))
  })
}

function fullWalk() {
  get("/api/fullWalk", {
    from: from.value,
    to: to.value,
    // YYYY-MM-DD
    date: date.value.toISOString().slice(0, 10),
  }).then(res => {
    trains.value = res as Train[]
  })
}

</script>

<template>
  <el-container>
    <el-header>
      <el-row :gutter="20">
        <el-col :span="3">
          <el-select filterable v-model="from" class="m-2" placeholder="从">
            <el-option
                v-for="item in options"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-col>
        <el-col :span="3">
          <el-select filterable v-model="to" class="m-2" placeholder="至">
            <el-option
                v-for="item in options"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-col>
        <el-col :span="5">
          <el-date-picker
              v-model="date"
              type="date"
              placeholder="日期"
          />
        </el-col>
        <el-col :span="3">
          <el-input v-model="no" placeholder="车次"/>
        </el-col>
        <el-col :span="4">
          <el-button-group>
            <el-button type="success" round @click=walk>Walk</el-button>
            <el-button type="primary" round @click=fullWalk>Search</el-button>
          </el-button-group>
        </el-col>
      </el-row>
    </el-header>
    <el-main>
      <!--show all train-->
      <el-table :data="trains" style="width: 100%" :row-class-name="tableRowClassName" border>
        <el-table-column prop="train_no" label="Train"/>
        <el-table-column prop="start_station_name" label="Start"/>
        <el-table-column prop="end_station_name" label="End"/>
        <el-table-column prop="from_station_name" label="From"/>
        <el-table-column prop="to_station_name" label="To"/>
        <el-table-column prop="start_time" label="Arrive" sortable/>
        <el-table-column prop="end_time" label="Start" sortable />
        <el-table-column prop="two_seat" label="Two"/>
        <el-table-column prop="one_seat" label="One"/>
        <el-table-column prop="special_seat" label="VIP"/>
        <el-table-column label="Operations">
          <template #default="scope">
            <el-button size="small">walk</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-main>
  </el-container>
</template>

<style>
.el-table .warning-row {
  --el-table-tr-bg-color: var(--el-color-warning-light-9);
}

.el-table .success-row {
  --el-table-tr-bg-color: var(--el-color-success-light-9);
}
</style>


<template>
  <el-container>
    <el-header height="100">
      <el-row :gutter="20">
        <el-col :span="3" :xs="12">
          <el-select filterable v-model="from" placeholder="从">
            <el-option
                v-for="item in options"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-col>
        <el-col :span="3" :xs="12">
          <el-select filterable v-model="to" class="m-2" placeholder="至">
            <el-option
                v-for="item in options"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-col>
        <el-col :span="3" :xs="24">
          <el-date-picker
              v-model="date"
              type="date"
              placeholder="日期"
              value-format="YYYY-MM-DD"
              :style="{width: '96%'}"
              :clearable="false"
          />
        </el-col>
        <el-col :span="2" :xs="12">
          <el-input v-model="no" placeholder="车次"/>
        </el-col>
        <el-col :span="4" :xs="12">
          <el-select
              v-model="chosenTypes"
              multiple
              collapse-tags
              collapse-tags-tooltip
              placeholder="车型"
              :max-collapse-tags="2"
          >
            <el-option
                v-for="item in types"
                :key="item"
                :label="item"
                :value="item"
            />
          </el-select>
        </el-col>
        <el-col :span="4" :xs="11">
          <el-slider v-model="delay" :min="1000" :max="6000"/>
        </el-col>
        <el-col :span="5" :xs="13">
          <el-button-group>
            <el-button type="primary" round @click=search :loading="searchLoading">Search</el-button>
            <el-button type="success" round @click=fullWalk :loading="fullWalkLoading">FullWalk</el-button>
          </el-button-group>
        </el-col>
      </el-row>
    </el-header>
    <el-main>
      <!--show all train-->
      <el-table
          :data="trains"
          :row-class-name="tableRowClassName"
          border
          table-layout="auto"
      >
        <el-table-column prop="train_no" label="Train" fixed/>
        <el-table-column prop="start_station_name" label="Start"/>
        <el-table-column prop="end_station_name" label="End"/>
        <el-table-column prop="from_station_name" label="From"/>
        <el-table-column prop="to_station_name" label="To"/>
        <el-table-column prop="start_time" label="Arrive" sortable/>
        <el-table-column prop="end_time" label="Start" sortable/>
        <el-table-column prop="two_seat" label="Two"/>
        <el-table-column prop="one_seat" label="One"/>
        <el-table-column prop="special_seat" label="VIP"/>
        <el-table-column label="Operations" :width="200">
          <template #default="scope">
            <el-button-group>
              <el-button size="default" type="success" round @click="inspect(scope.row, false)">Inspect</el-button>
              <el-button size="default" type="success" round @click="inspect(scope.row, true)">Reverse</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-main>
  </el-container>

  <el-dialog
      v-model="dialogVisible"
      title="Train Pass Stations"
      width="80%"
      :before-close="handleClose"
  >
    <el-timeline>
      <el-timeline-item
          v-for="(p, index) in pass"
          :key="index"
          :timestamp="p.arrive_time"
          hide-timestamp
      >
        <el-row>
          <el-col :span="4" :xs="6">
            <el-button :size="'small'" text>{{ p.station_name }}</el-button>
          </el-col>
          <el-col :span="6" :xs="8">
            <el-button :size="'small'" text>{{ p.arrive_time }}-{{ p.start_time }}</el-button>
          </el-col>
          <el-col :span="14" :xs="24">
            <el-badge :value="p.two_seat" :type="p.two_seat!='无'? 'success':'warning'">
              <el-button :size="'small'" round>Two</el-button>
            </el-badge>
            <el-badge :value="p.one_seat" :type="p.one_seat!='无'? 'success':'warning'">
              <el-button :size="'small'" round>One</el-button>
            </el-badge>
            <el-badge :value="p.special_seat" :type="p.special_seat!='无'? 'success':'warning'">
              <el-button :size="'small'" round>Spe</el-button>
            </el-badge>
          </el-col>
        </el-row>
      </el-timeline-item>
    </el-timeline>
  </el-dialog>
</template>

<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import type {Pass, Stations, Train} from "@/api";
import {init, originalPass, originalSearch} from '@/api';
import {ElNotification} from 'element-plus'

const dialogVisible = ref(false)

const date = ref('')
const from = ref('')
const to = ref('')
const no = ref('')

const types = ref(['G', 'D', 'K', 'Z', 'C'])
const chosenTypes = ref<string[]>([])
const trains = ref<Train[]>([])
const pass = ref<Pass[]>([])
const delay = ref<number>(2000)

const searchLoading = ref(false)
const fullWalkLoading = ref(false)


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
  from.value = window.localStorage.getItem("from") as string;
  to.value = window.localStorage.getItem("to") as string;
  date.value = window.localStorage.getItem("date") as string;
  no.value = window.localStorage.getItem("no") == null ? '' : window.localStorage.getItem("no") as string;
  chosenTypes.value = window.localStorage.getItem("types") == null ? ['G', 'D'] : window.localStorage.getItem("types")!.split(',');

  init().then(res => options.value = res)
})

function search() {
  searchLoading.value = true
  const req =
      {
        from: from.value,
        to: to.value,
        // YYYY-MM-DD
        date: date.value,
        no: no.value,
        types: chosenTypes.value
      }
  originalSearch(req.from, req.to, req.date, req.no, req.types)
      .then(res => {
        trains.value = res as Train[]
      })
      .catch(err => handelError(err))
  searchLoading.value = false
  //save req to localstorage
  window.localStorage.setItem("from", req.from)
  window.localStorage.setItem("to", req.to)
  window.localStorage.setItem("date", req.date)
  window.localStorage.setItem("no", req.no)
  window.localStorage.setItem("types", req.types.join(','))
}

async function fullWalk() {
  fullWalkLoading.value = true
  for (let i = 0; i < trains.value.length; i++) {
    const t = trains.value[i]
    originalSearch(
        t.start_station,
        t.end_station,
        // YYYY-MM-DD UTC+8
        date.value,
        t.train_code,
        types.value
    ).then(res => {
      var re = (res as Train[])[0];
      re.from_station_name = '@' + t.from_station_name
      re.to_station_name = '@' + t.to_station_name
      re.start_time = '@' + t.start_time
      re.end_time = '@' + t.end_time
      trains.value[i] = re
    }).catch(err => handelError(err))
    await new Promise(resolve => setTimeout(resolve, delay.value)); // 1000毫秒 = 1秒
  }
  fullWalkLoading.value = false
}

function inspect(t: Train, reverse: boolean) {
  dialogVisible.value = true
  originalPass(
      t.from_station,
      t.to_station,
      date.value,
      t.train_code
  ).then(async res => {
    pass.value = res
    for (let i = 1; i < pass.value.length; i++) {
      let from = pass.value[0].station
      let to = pass.value[i].station
      if (reverse) {
        from = pass.value[i - 1].station
        to = pass.value[pass.value.length - 1].station
      }
      originalSearch(
          from,
          to,
          // YYYY-MM-DD
          date.value,
          t.train_code,
          types.value
      ).then(res => {
        const re = (res as Train[])[0];
        let target = i;
        if (reverse) {
          target--
        }
        pass.value[target].two_seat = re.two_seat
        pass.value[target].one_seat = re.one_seat
        pass.value[target].special_seat = re.special_seat
      })
          .catch(err => handelError(err))
      await new Promise(resolve => setTimeout(resolve, delay.value)); // 1000毫秒 = 1秒
    }
  }).catch(err => handelError(err))
}

function handelError(err: any) {
  ElNotification({
    message: err as string,
    type: 'error',
  })
}


const handleClose = (done: () => void) => {
  pass.value = []
  done()
}
</script>


<style>
.el-table .warning-row {
  --el-table-tr-bg-color: var(--el-color-warning-light-8);
}

.el-table .success-row {
  --el-table-tr-bg-color: var(--el-color-success-light-8);
}
</style>


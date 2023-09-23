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
              value-format="YYYY-MM-DD"
          />
        </el-col>
        <el-col :span="3">
          <el-input v-model="no" placeholder="车次"/>
        </el-col>
        <el-col :span="6">
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
          style="width: 100%"
          :row-class-name="tableRowClassName"
          border
          size="default"
      >
        <el-table-column prop="train_no" label="Train"/>
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
      width="60%"
      :before-close="handleClose"
  >
    <el-timeline>
      <el-timeline-item
          v-for="(p, index) in pass"
          :key="index"
          :timestamp="p.arrive_time"
          hide-timestamp
          center
      >
        <el-row>
          <el-col :span="4">
            <el-button :size="'small'">{{ p.station_name }}</el-button>
          </el-col>
          <el-col :span="6">
            <el-button :size="'small'">{{ p.arrive_time }}-{{ p.start_time }}</el-button>
          </el-col>
          <el-col :span="12">
            <el-badge :value="p.two_seat" :type="p.two_seat!='无'? 'success':'warning'">
              <el-button :size="'small'">Two</el-button>
            </el-badge>
            <el-badge :value="p.one_seat" :type="p.one_seat!='无'? 'success':'warning'">
              <el-button :size="'small'">One</el-button>
            </el-badge>
            <el-badge :value="p.special_seat" :type="p.special_seat!='无'? 'success':'warning'">
              <el-button :size="'small'">Spe</el-button>
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
import {get} from '@/api';

const dialogVisible = ref(false)

const date = ref('2023-10-01')
const from = ref('UUH')
const to = ref('TNV')
const no = ref('G3133')
const trains = ref<Train[]>([])
const pass = ref<Pass[]>([])

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
  get('/stations').then(res => {
    options.value = (res as Stations[])
  })
})

function search() {
  searchLoading.value = true
  get("/search", {
    from: from.value,
    to: to.value,
    // YYYY-MM-DD
    date: date.value,
    no: no.value
  }).then(res => {
    trains.value = res as Train[]
  })
  searchLoading.value = false
}

async function fullWalk() {
  fullWalkLoading.value = true
  for (let i = 0; i < trains.value.length; i++) {
    const t = trains.value[i]
    get("/search", {
      from: t.start_station,
      to: t.end_station,
      // YYYY-MM-DD UTC+8
      date: date.value,
      no: t.train_code
    }).then(res => {
      var re = (res as Train[])[0];
      re.from_station_name = '@' + t.from_station_name
      re.to_station_name = '@' + t.to_station_name
      re.start_time = '@' + t.start_time
      re.end_time = '@' + t.end_time
      trains.value[i] = re
    })
    await new Promise(resolve => setTimeout(resolve, 1000)); // 1000毫秒 = 1秒
  }
  fullWalkLoading.value = false
}


function inspect(t: Train, reverse: boolean) {
  dialogVisible.value = true
  get("/pass", {
    from: t.from_station,
    to: t.to_station,
    // YYYY-MM-DD
    date: date.value,
    no: t.train_code
  }).then(async res => {
    pass.value = res as Pass[]
    for (let i = 1; i < pass.value.length; i++) {
      let from = pass.value[0].station
      let to = pass.value[i].station
      if (reverse) {
        from = pass.value[i - 1].station
        to = pass.value[pass.value.length - 1].station
      }
      get("/search", {
        from: from,
        to: to,
        // YYYY-MM-DD
        date: date.value,
        no: t.train_code
      }).then(res => {
        const re = (res as Train[])[0];
        let target = i;
        if (reverse){
          target--
        }
        pass.value[target].two_seat = re.two_seat
        pass.value[target].one_seat = re.one_seat
        pass.value[target].special_seat = re.special_seat
      })
      await new Promise(resolve => setTimeout(resolve, 1000)); // 1000毫秒 = 1秒
    }
  })
}

const handleClose = (done: () => void) => {
  pass.value = []
  done()
}
</script>


<style>
.el-table .warning-row {
  --el-table-tr-bg-color: var(--el-color-warning-light-9);
}

.el-table .success-row {
  --el-table-tr-bg-color: var(--el-color-success-light-9);
}
</style>


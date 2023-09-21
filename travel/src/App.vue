<script lang="ts" setup>
import {onMounted, ref} from 'vue'
import {get} from '@/api';

const date = ref()
const from = ref('')
const to = ref('')
const no = ref('')

export interface Stations {
  value: string
  label: string
}

const options = ref<Stations[]>([])

onMounted(() => {
  get('/api/stations').then(res => {
    options.value = (res as Stations[])
  })
})

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
        <el-col :span="2">
          <el-button type="success" round>Walk</el-button>
        </el-col>
        <el-col :span="2">
          <el-button type="primary" round>Search</el-button>
        </el-col>
      </el-row>
    </el-header>
    <el-main>
    </el-main>
  </el-container>
</template>

<style scoped>
</style>


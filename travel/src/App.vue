<script lang="ts" setup>
import {ref} from 'vue'

const date = ref()
const from = ref('')
const to = ref('')
const no = ref('')

const options = [
  {
    value: 'Option1',
    label: 'Option1',
  },
  {
    value: 'Option2',
    label: 'Option2',
  }
]
</script>

<template>
  <el-container>
    <el-header>
      <el-row :gutter="20">
        <el-col :span="4">
          <el-select v-model="from" class="m-2" placeholder="从">
            <el-option
                v-for="item in options"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="to" class="m-2" placeholder="至">
            <el-option
                v-for="item in options"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-date-picker
              v-model="date"
              type="date"
              placeholder="日期"
          />
        </el-col>
        <el-col :span="4">
          <el-input v-model="no" placeholder="车次"/>
        </el-col>
        <el-col :span="4">
          <el-button type="success" round>Walk</el-button>
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


<template>
  <div class="wrapper">
    <el-table
      :data="tableData"
      stripe
      style="width: 100%">
      <el-table-column
        type="index"
        width="64px"
      >
        <template slot-scope="scope">
          <span>{{ formatIndex(scope.$index) }}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="Node"
      >
        <template slot-scope="scope">
          <el-tooltip :content="scope.row.pubkey" :open-delay="500">
            <span class="node-id">{{ formatAddr(scope.row.pubkey) }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column
        prop="height"
        label="Height"
      >
        <template slot-scope="scope">
          <span class="node-height">{{formatNum(scope.row.height)}}</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="latency"
        label="Latency"
      >
        <template slot-scope="scope">
          <span :class="{slowly: scope.row.latency > 500}">{{scope.row.latency}}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="LastOnline"
      >
        <template slot-scope="scope">
          <span>{{formatTime(scope.row.lastOnlineTime)}}</span>
        </template>
      </el-table-column>
      <el-table-column
        label="Syncing"
        width="80"
        align="center"
      >
      <template slot-scope="scope">
        <i class="dot" :class="{'is-active': scope.row.healthy ||!scope.row.syncing }"></i>
      </template>
      </el-table-column>
    </el-table>
    <div class="folder"  v-if="datas.length >= 100" @click="toggleFolder">
      <i class="icon" :class="{expanded: showAll}"></i>
      <span>{{ folderText }}</span>
    </div>
  </div>
</template>

<script>
module.exports = {
  name: 'StatTable',
  props: {
    datas: {
      type: Array,
      default: () => []
    },
    indexLength: {
      type: Number,
      default: 2
    }
  },
  data() {
    return {
      showAll: false
    }
  },
  computed: {
    length() {
      return this.showAll ? this.datas.length : 100
    },
    tableData() {
      return this.datas.slice(0, this.length)
    },
    folderText() {
      return this.showAll ? 'Less' : 'More'
    }
  },
  methods: {
    formatNum(num) {
      return num.toLocaleString();
    },
    formatTime(str) {
      return dayjs(str).format('HH:mm:ss') 
    },
    formatAddr(str) {
      const id = str
      if (id.length > 10) {
        return id.substring(0, 5) + '...' + id.substring(id.length - 5);
      }
      return id || '';
    },
    formatIndex(index) {
      let idx = index + 1
      if (this.indexLength === 2) {
        if (idx < 10) {
          return '0' + idx
        } else {
          return idx
        }
      } else {
        if (idx < 10) {
          return '00' + idx
        } else if (idx >= 10 && idx < 100){
          return '0' + idx
        } else {
          return idx
        }
      }
    },
    toggleFolder() {
      this.showAll = !this.showAll
    }
  }
}
</script>

<style>
.wrapper {
  position: relative;
  width: 100%;
}
.folder {
  cursor: pointer;
  width: 100%;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #43e5fa;
  background-color: rgba(0, 0, 0, 0.7);
  position: relative;
}
.folder .icon {
  display: inline-block;
  width: 20px;
  height: 20px;
  margin-right: 16px;
  transition: all 0.3s;
  background: url('../img/arrow.svg') no-repeat center;
  background-size: 100%;
}
.folder .icon.expanded {
  transform: rotate(180deg);
}
.stat-table {
  margin-bottom: 70px;
}
.stat-table .description{
  margin-bottom: 30px;
  font-size: 30px;
  color: #ffffff;
}
.stat-table .el-table {
  background: none;
}
.stat-table .el-table::before {
  display: none;
}
.stat-table .el-table thead {
  background: url('../img/thead.png') no-repeat center;
  background-size: 100% 100%;
}
.stat-table .el-table .el-table__cell {
  border: none !important;
}
.stat-table .el-table th,tr{
  background: none !important;
  color: #5d687c;
}
.stat-table .el-table .el-table__body-wrapper {
  border: 1px solid #00d18c;
  border-radius: 4px;
  background-color: #0f1218;
}
.stat-table .el-table .el-table__body-wrapper * {
  font-family: Electrolize !important;
}
.stat-table .el-table--enable-row-hover .el-table__body tr:hover>td.el-table__cell {
  background-color: inherit !important;
}
.stat-table .el-table--striped .el-table__body tr.el-table__row--striped td.el-table__cell {
  background-color: #14181f;
}
.stat-table .el-table .slowly {
  color: #E6A23C;
}
.stat-table .el-table .dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
  position: relative;
  background: #F56C6C;
  margin-right: 5px;
}
.stat-table .el-table .dot.is-active {
  background: #67C23A;
}

.stat-table .el-table .node-id {
  color: #987abd;
}
.stat-table .el-table .node-height {
  color: #1b946c;
  font-family: Aldrich !important;
}
</style>
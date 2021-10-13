<template>
  <div class="stat-chart">
    <ul class="info-list">
      <li class="info-item">
        <div class="label">Requests count</div>
        <div class="value">{{formatNum(total)}}</div>
      </li>
      <li class="info-item">
        <div class="label">Requests per hour (24H Average)</div>
        <div class="value">{{formatNum(average)}}</div>
      </li>
      <li class="info-item">
        <div class="label">Requests history</div>
        <div class="tags">
          <div class="tag-item" v-for="item in tags" :key="item.value" :class="{'is-active': curr === item.value}" @click="setActive(item.value)">
            {{item.label}}
          </div>
        </div>
      </li>
    </ul>
    <div class="chart" ref="chart" />
  </div>
</template>

<script>
module.exports = {
  name: 'StatChart',
  props: {
    datas: {
      type: Object,
      default: () => ({})
    },
    total: {
      type: Number,
      default: 0
    },
    average: {
      type: Number,
      default: 0
    }
  },
  data() {
    return {
      tags: [{
        label: '24h',
        value: 24
      }],
      curr: 24,
      chart: null
    }
  },
  computed: {
    chartData() {
      return Object.entries(this.datas).map(item => [dayjs(item[0] * 1000).format('HH:mm'), item[1]]);
      
    },
    option() {
      return {
        color: '#8f40f0',
        tooltip: {
          axisPointer: { type: 'none' },
          backgroundColor: '#06090f',
          borderColor: '#06090f',
          textStyle: {
            fontSize: 14,
            color: '#ffffff',
            align: 'left',
            fontFamily: 'Aldrich'
          }
        },
        grid: {
          left: '5%',
          right: '5%',
          bottom: '4%',
          top: '10%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          axisTick: { show: false },
          axisLabel: { show: false },
          axisLine: { show: false },
          splitLine: { show: false }
        },
        yAxis: {
          type: 'value',
          axisLine: { show: false },
          minInterval: 1,
          axisLabel: { show: true, color: '#5d687c' },
          axisTick: { show: false },
          splitLine: { show: false }
        },
        series: [
          {
            name: 'Requests',
            type: 'bar',
            barMaxWidth: '16',
            itemStyle: { borderRadius: [4, 4, 0, 0] },
            emphasis: { itemStyle: { color: '#6f25c9' } },
            data: this.chartData
          }
        ]
      }
    }
  },
  watch: {
    option: {
      handler() {
        this.handleChart()
      },
      deep: true
    }
  },
  mounted() {
    window.addEventListener('resize', this.resizeHandler)
    this.handleChart();
    this.$once('hook:beforeDestroy', () => {
      window.removeEventListener('resize', this.resizeHandler)
    })
  },
  methods: {
    formatNum(num) {
      return num.toLocaleString();
    },
    resizeHandler() {
      this.chart && this.chart.resize();
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
    setActive(val) {
      this.curr = val;
    },
    handleChart() {
      if (this.chart) {
        this.setChart()
      } else {
        this.chart = echarts.init(this.$refs.chart)
        this.setChart()
      }
    },
    setChart() {
      this.chart.setOption(this.option)
    }
  }
}
</script>

<style>
.stat-chart {
  background-color: #0f1218;
	border-radius: 4px;
	border: solid 1px #00d18c;
}
.info-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 50px;
  font-size: 16px;
  color: #5d687c;
}
.info-item:nth-of-type(2) {
  background: #14181f;
}
.info-item .value {
  color: #ffffff;
}
.info-item .tags {
  display: flex;
  align-items: center;
}
.info-item .tag-item {
  padding: 6px 10px;
  margin: 0 2px;
  cursor: pointer;
}
.info-item .tag-item.is-active {
  background-color: #00d18c;
	border-radius: 4px;
  color: #0f1218;
}
.stat-chart .chart {
  width: 100%;
  height: 300px;
}
</style>
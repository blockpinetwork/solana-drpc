ELEMENT.locale(ELEMENT.lang.en);
document.title = window.customTitle;
const StatTable = httpVueLoader('./components/StatTable.vue')
const StatChart = httpVueLoader('./components/StatChart.vue')
const CustomFoot = httpVueLoader('./components/CustomFoot.vue')
var app = new Vue({
  el:"#app",
  components: {
    StatTable,
    StatChart,
    CustomFoot
  },
  data() {
    return{
      tableData: [],
      selectedData: [],
      total_peers: 0,
      total_rpc_backends: 0,
      total_rpc_nodes: 0,
      timer: null,
      copyText: 'copy',
      total_request: 0,
      average_request: 0,
      requests: {}
    }
  },
  computed: {
    endpoint() {
      return window.solana_rpc_url
    }
  },
  mounted() {
    this.getStat()
    this.timer = setInterval(()=>{
      this.getStat()
    }, 10000)
    this.$once('hook:beforeDestroy', () => {
      clearInterval(this.timer);
      this.timer = null;
    });
  },
  methods: {
    formatNum(num) {
      return num.toLocaleString();
    },
    getStat() {
      request.get('/api/status').then((res)=> {
        this.tableData = res.data.all_rpc_nodes || []
        this.selectedData = res.data.backend_rpc_nodes || []
        this.total_peers = res.data.total_peers || 0
        this.total_rpc_backends = res.data.total_rpc_backends || 0
        this.total_rpc_nodes = res.data.total_rpc_nodes || 0
        this.total_request = res.data.total_requests || 0
        this.average_request = res.data.hourly_avg_requests_24h || 0
        this.requests = res.data.requests || {}
      }).catch(() => {
        this.tableData = []
        this.selectedData = []
        this.total_peers = 0
        this.total_rpc_backends = 0
        this.total_rpc_nodes = 0
        this.total_request = 0
        this.average_request = 0
        this.requests = {}
      })
    },
    handleCopy() {
      navigator.clipboard.writeText(this.endpoint).then(() => {
        this.copyText = 'copied';
        setTimeout(() => {
          this.copyText = 'copy';
        }, 1000);
      });
    }
  }
})


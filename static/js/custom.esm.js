import * as echarts from 'echarts';

let chartDomElems = document.querySelectorAll('#graph');

Array.from(chartDomElems).forEach((elem) => {
  var chart = echarts.init(elem);
  var option;

  option = {
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        data: [820, 932, 901, 934, 1290, 1330, 1320],
        type: 'line',
        areaStyle: {}
      }
    ]
  };
  window.addEventListener('resize', function() {
    chart.resize();
  });
  option && chart.setOption(option);
});

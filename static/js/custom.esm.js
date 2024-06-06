import * as echarts from 'echarts';
let chartDomElems = document.querySelectorAll('#graph');

Array.from(chartDomElems).forEach((elem) => {
  let chart = echarts.init(elem);
  fetch("/getData").then(async (response) => {
    let option = await response.json();
    console.log("Option:", option);
    option && chart.setOption(option) && chart.resize();
  }).catch((error) => {
    console.error('Error:', error);
  });
  window.addEventListener('resize', function() {
    chart.resize();
  });
});

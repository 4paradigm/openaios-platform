/*
 * @Author: liyuying
 * @Date: 2021-05-07 16:27:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-05-08 15:58:27
 * @Description: file content
 */

import React, { useEffect, useRef } from 'react';
import { Chart, XChartType, XConfiguration, XDataSerial, XDataTuple } from 'smart-chart';
import { useSelector, IHomeState } from 'umi';
const colorList = ['#28a080', '#622cd2', '#32aaff', '#825a28', '#ff8250', '#e8a858', '#3278ff'];
function ResourceChart() {
  const { taskInfo }: IHomeState = useSelector((state: any) => state.home);
  const ref = useRef(null);
  let chart: Chart;

  const initBartChart = () => {
    let serialsData: XDataTuple[] = [];
    if (taskInfo && taskInfo.task_list) {
      taskInfo.task_list.forEach((item: any) => {
        serialsData.push([item.compute_unit, item.number * item.price]);
      });
    }
    const serials: XDataSerial[] = [
      {
        name: '每分钟消耗',
        data: serialsData,
      },
    ];
    if (!chart) {
      chart = new Chart((ref.current as unknown) as HTMLDivElement);
    }
    const configuration: XConfiguration = {
      type: XChartType.pie,
      serials: serials,
      optionConfiguration: {
        colors: colorList,
        showLegend: true,
      },
    };
    chart.render(configuration);
  };

  useEffect(() => {
    initBartChart();
  }, [taskInfo]);

  return <div style={{ height: 200 }} ref={ref}></div>;
}
export default ResourceChart;

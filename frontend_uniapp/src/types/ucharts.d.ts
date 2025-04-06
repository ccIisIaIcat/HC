declare module '@qiun/ucharts/u-charts' {
  export default class uCharts {
    constructor(options: any);
    touchLegend(e: any): void;
    showToolTip(e: any): void;
  }
}

declare module '@/components/u-charts/u-charts.js' {
  const uCharts: any;
  export default uCharts;
} 
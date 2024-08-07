export default {
  common: {
    lang: {
      zh: '简体中文',
      en: 'English'
    },
    title: {
      operation: 'Operation',
      createTime: 'CreateTime',
      creator: 'Creator',
      startTime: 'StartTime',
      endTime: 'EndTime',
      to: 'To',
      baseInfo: 'Basic Information',
      metrics: 'Metrics',
      testIndicator: 'Test Indicator',
      data: 'Data',
      parameters: 'Parameters',
      dataset: 'Dataset',
      noContrast: 'No Contrast',
      customContrast: 'Specified Contrast Period',
      noRefresh: 'No Auto Refresh',
      customRefresh: 'Specified Refresh Period',
      time: 'Time',
      failure: 'Failure',
      causes: 'Causes',
      unknown: 'Unknown'
    },
    action: {
      delete: 'Delete',
      operationLog: 'Operation Log',
      cancel: 'Cancel',
      refresh: 'Refresh'
    },
    inst: {
      deleteTips: 'Are you sure you want to delete?'
    },
    time: {
      lastHour: 'Last Hour',
      lastDay: 'Last Day',
      lastWeek: 'Last Week',
      lastMonth: 'Last Month',
      hour: 'Hour',
      minute: 'Minute',
      second: 'Second',
      day: 'Day',
      week: 'Week'
    }
  },
  menu: {
    service: 'LLM Service',
    record: 'Injection Records'
  },
  instance: {
    title: {
      llmInstance: 'LLM Service',
      instanceName: 'Service Name',
      deployStatus: 'Deployment Status',
      deploying: 'Deploying',
      deploySuccess: 'Deployment Success',
      deployFail: 'Deployment Fail',
      finished: 'Finished',
      instanceDetail: 'LLM Deployment Instance Detail',
      testConfig: 'Test Configuration',
      startupOptions: 'Startup Parameter',
      instanceSpec: 'Instance Specifications',
      instanceImage: 'Instance Image',
      area: 'Area',
      memory: 'Memory',
      systemDisk: 'System Disk',
      dataDisk: 'Data Disk',
      imageSource: 'Image Source',
      image: 'Image',
      injectTime: 'Traffic Injection Time',
      distribution: 'Distribution',
      poisson: 'Poisson Distribution',
      traffic: 'Gaussian Distribution',
      otherParams: 'Other Parameters',
      sysDataset: 'System Dataset',
      userDataset: 'Custom Dataset',
      configSuggest: 'Configuration Suggestions',
      currentVal: 'Current Value',
      suggestVal: 'Recommended Value'
    },
    action: {
      inject: 'Request Injection',
      startTest: 'Start Test'
    },
    inst: {
      searchPlaceholder: 'Instance Name, Deployment Status',
      datasetPlaceholder: 'Please Select The Dataset',
      paramsPlaceholder: 'Please Select The Parameters',
      chartSearch: 'Indicator Name、Indicator Type',
      durationTip: 'Duration {0}'
    }
  },
  experiment: {
    title: {
      testRecord: 'Test Records',
      testResult: 'Test Result',
      testLog: 'Test Log',
      testTime: 'Test Time',
      testId: 'Test ID',
      testInstance: 'Test Instance',
      testDataset: 'Test Dataset',
      testStatus: 'Test Status',
      prompt: 'Average Prompt Throughput Per Second',
      generation: 'Average Generation Throughput Per Second',
      testUser: 'Test User',
      testing: 'Testing',
      testSuccess: 'Test Success',
      testFail: 'Test Fail',
      testInit: 'Init',
      testUnknown: 'Unknown',
      finished: 'Test Finished',
      testResultOverview: 'Overview of Test Results',
      testConfig: 'Test Configuration',
      requestNum: 'Total Number of Requests Sent',
      requestSuccessNum: 'Number of Successful Requests Executed',
      requestSuccessRate: 'Request Execution Success Rate',
      avgTime: 'Average Execution Time of Requests'
    },
    action: {},
    inst: {
      searchPlaceholder: 'Test Instance, Test Dataset',
      successTip: 'Request Injection Successfully!'
    }
  },
  datepicker: {
    placeholder: 'Select Date and Time',
    to: 'to',
    startDate: 'Start Date',
    endDate: 'End Date',
    today: 'Today',
    lastWeek: 'Last Week',
    lastMonth: 'Last 30 Days',
    lastThreeMonths: 'Last 3 Months',
    lastTwelveHours: 'Last 12 Hours',
    lastOneHour: 'Last Hour',
    lastThreeHours: 'Last 3 Hours',
    lastTwentyFourHours: 'Last 24 Hours',
    lastTwoDays: 'Last 2 Days',
    lastSevenDays: 'Last 7 Days',
    day: 'day',
    hour: 'hour',
    minute: 'minute'
  },
  chart: {
    title: {
      requestNum: 'Number of requests received per second',
      requestExecuteNum: 'Number of requests returned per second',
      requestSuccessNum: 'Number of successful requests executed per second',
      requestSuccessRate: 'Success rate of request execution per second',
      activeRequestNum: 'The number of requests being processed',
      requestDuration: 'Average request execution time per second',
      responseSize: 'Response size',
      requestSize: 'Request size',
      promptThroughput: 'avg_prompt_throughput',
      generationThroughput: 'avg_generation_throughput',
      runningRequests: 'running_requests',
      pendingRequests: 'pending_requests',
      gpuKv: 'gpu_kv_cache_usage',
      timeToFirst: 'time_to_first_token_seconds',
      timePerOutput: 'time_per_output_token_seconds',
      gpuRate: 'DCGM_FI_DEV_GPU_UTIL',
      memRate: 'DCGM_FI_DEV_MEM_COPY_UTIL',
      fbNum: 'DCGM_FI_DEV_FB_USED',
      smClock: 'DCGM_FI_DEV_SM_CLOCK',
      memClock: 'DCGM_FI_DEV_MEM_CLOCK',
      memTemp: 'DCGM_FI_DEV_MEMORY_TEMP',
      gpuTemp: 'DCGM_FI_DEV_GPU_TEMP',
      power: 'DCGM_FI_DEV_POWER_USAGE',
      grEnginActive: 'DCGM_FI_PROF_GR_ENGINE_ACTIVE',
      tensor: 'DCGM_FI_PROF_PIPE_TENSOR_ACTIVE',
      memCopyUtil: 'DCGM_FI_PROF_DRAM_ACTIVE',
      pcie: 'DCGM_FI_PROF_PCIE_TX_BYTES DCGM_FI_PROF_PCIE_RX_BYTES',
      nvLink: 'DCGM_FI_PROF_NVLINK_RX_BYTES DCGM_FI_PROF_NVLINK_TX_BYTES',
      noData: 'There is currently no data available'
    }
  }
}

<script setup>
import { reactive } from 'vue';
let state = reactive({
  orders: [
    { num: '10001', title: '集成 Android SDK 报错', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10002', title: '集成 iOS SDK 报错', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10003', title: '问题怎么看', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10004', title: '问题怎么看', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10005', title: '问题怎么看', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10006', title: '问题怎么看', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10007', title: '问题怎么看', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10008', title: '问题怎么看', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10009', title: '集成 Android SDK 报错', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10010', title: '集成 Android SDK 报错', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10011', title: '集成 Android SDK 报错', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10012', title: '集成 Android SDK 报错', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10013', title: '集成 Android SDK 报错', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10014', title: '集成 Android SDK 报错', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10015', title: '集成 Android SDK 报错', status: '处理中', time: '2023-10-10 23:03' },
    { num: '10016', title: '集成 Android SDK 报错', status: '处理中', time: '2023-10-10 23:03' },
  ]
});

</script>
<template>
  <div class="mb-4">
    <table class="table cim-table">
      <thead>
        <tr>
          <th scope="col" class="cim-td-c">工单号</th>
          <th scope="col">标题</th>
          <th scope="col" class="cim-td-c">状态</th>
          <th scope="col" class="cim-td-c">最后回复时间</th>
          <th scope="col">操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="order in state.orders">
          <td class="cim-td-c">{{ order.num }}</td>
          <td>{{ order.title }}</td>
          <td class="cim-td-c">{{ order.status }}</td>
          <td class="cim-td-c">{{ order.time }}</td>
          <td>
            <a class="btn-link cim-btn-link" type="button">查看</a>
          </td>
        </tr>
      </tbody>
    </table>

    <nav class="cim-navigation">
      <ul class="pagination">
        <li class="page-item"><a class="page-link" href="#" aria-label="Previous"><span aria-hidden="true">«</span></a>
        </li>
        <li class="page-item"><a class="page-link" href="#">1</a></li>
        <li class="page-item"><a class="page-link" href="#">2</a></li>
        <li class="page-item"><a class="page-link" href="#">3</a></li>
        <li class="page-item"><a class="page-link" href="#" aria-label="Next"><span aria-hidden="true">»</span></a></li>
      </ul>
    </nav>
  </div>
</template>

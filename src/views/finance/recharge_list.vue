<script setup>
import { reactive } from 'vue';
let state = reactive({
  orders: [
  { id: 1, title: '北京未来科技有限公司', type: 1, amount: 25000, status: '已寄出', time: '2023-10-10 23:03' },
  { id: 1, title: '北京未来科技有限公司', type: 1, amount: 25000, status: '已寄出', time: '2023-10-10 23:03' },
  { id: 1, title: '北京未来科技有限公司', type: 2, amount: 25000, status: '已寄出', time: '2023-10-10 23:03' },
  { id: 1, title: '北京未来科技有限公司', type: 2, amount: 25000, status: '已寄出', time: '2023-10-10 23:03' },
  { id: 1, title: '北京未来科技有限公司', type: 1, amount: 25000, status: '已寄出', time: '2023-10-10 23:03' },
  ]
});
</script>
<template>
   <div class="mb-4">
    <ul class="nav nav-underline-border">
        <li class="nav-item"><a class="nav-link active cicon cicon-list">充值记录</a></li>
      </ul>
    <table class="table cim-table">
      <thead>
        <tr>
          <th scope="col" class="cim-td-c">开票抬头</th>
          <th scope="col">开票金额</th>
          <th scope="col" class="cim-td-c">发票类型</th>
          <th scope="col" class="cim-td-c">状态</th>
          <th scope="col" class="cim-td-c">开票时间</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="order in state.orders">
          <td class="cim-td-c">{{ order.title }}</td>
          <td class="cim-td-c">{{ order.amount }}</td>
          <td>{{ order.type == 1 ? '普通发票' : '专用发票' }}</td>
          <td class="cim-td-c">{{ order.status }}</td>
          <td class="cim-td-c">{{ order.time }}</td>
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

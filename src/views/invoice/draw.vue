<script setup>
import { reactive } from 'vue';
import { INVOICE_TYPE } from "../../common/enum";
import utils from '../../common/utils';
let state = reactive({
  records: [
    { id: 1, title: '北京未来科技有限公司', checked: false, recharge: 25000, invoice: 25000, type: '微信', time: '2023-10-10 23:03' },
    { id: 2, title: '北京未来科技有限公司', checked: false, recharge: 25000, invoice: 25000, type: '微信', time: '2023-10-10 23:03' },
    { id: 3, title: '北京未来科技有限公司', checked: false, recharge: 25000, invoice: 25000, type: '微信', time: '2023-10-10 23:03' },
    { id: 4, title: '北京未来科技有限公司', checked: false, recharge: 25000, invoice: 25000, type: '微信', time: '2023-10-10 23:03' },
    { id: 5, title: '北京未来科技有限公司', checked: false, recharge: 25000, invoice: 25000, type: '微信', time: '2023-10-10 23:03' },
    { id: 6, title: '北京未来科技有限公司', checked: false, recharge: 25000, invoice: 25000, type: '微信', time: '2023-10-10 23:03' },
    { id: 7, title: '北京未来科技有限公司', checked: false, recharge: 25000, invoice: 25000, type: '微信', time: '2023-10-10 23:03' },
    { id: 8, title: '北京未来科技有限公司', checked: false, recharge: 25000, invoice: 25000, type: '微信', time: '2023-10-10 23:03' },
    { id: 9, title: '北京未来科技有限公司', checked: false, recharge: 25000, invoice: 25000, type: '微信', time: '2023-10-10 23:03' },
    { id: 10, title: '北京未来科技有限公司', checked: false, recharge: 25000, invoice: 25000, type: '微信', time: '2023-10-10 23:03' },
  ],
  radios: [
    { name: 'type', value: INVOICE_TYPE.ONLINE, label: '增值税普通发票（电子）' },
    { name: 'type', value: INVOICE_TYPE.PAPER, label: '增值税普通发票（纸质）' },
  ],
  total: 0,
  invoiceType: INVOICE_TYPE.ONLINE,
  checkAll: false
});
function onRadieChanged(type){
  state.invoiceType = type;
}
function onCheckboxChanged(item){
  state.records.map((record) => {
    if(utils.isEqual(record.id, item.id)){
      record.checked = !item.checked;
    }
    return record;
  });

  let total = 0;
  utils.each(state.records, (record) => {
    if(record.checked){
      total += record.invoice
    }
  });
  state.total = utils.numberWithCommas(total);
}
function onCheckAll(){
  state.checkAll = !state.checkAll;
  let total = 0;
  state.records.map((record) => {
    record.checked = state.checkAll;
    if(state.checkAll){
      total += record.invoice
    }
    return record;
  });
  state.total = utils.numberWithCommas(total);
}
</script>
<template>
  <div class="mb-4">
    <div class="card-body">
      <ul class="nav nav-underline-border">
        <li class="nav-item"><a class="nav-link active cicon cicon-list">待开票记录</a></li>
      </ul>
      <div class="tab-content rounded-bottom">
        <div class="tab-pane p-3 active cim-wait-list">
          <table class="table cim-table">
            <thead>
              <tr>
                <th class="cim-td-c">
                  <input class="form-check-input cim-td-select" type="checkbox" @change="onCheckAll()">
                </th>
                <th>充值主体</th>
                <th class="cim-td-c">充值金额</th>
                <th class="cim-td-c">可开票金额</th>
                <th class="cim-td-c">充值方式</th>
                <th class="cim-td-c">充值时间</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="record in state.records">
                <td class="cim-td-c">
                  <input class="form-check-input cim-td-select" :checked="record.checked"  type="checkbox" @change="onCheckboxChanged(record)">
                </td>
                <td>{{ record.title }}</td>
                <td class="cim-td-c">{{ record.recharge }}</td>
                <td class="cim-td-c">{{ record.invoice }}</td>
                <td class="cim-td-c">{{ record.type }}</td>
                <td class="cim-td-c">{{ record.time }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="tab-pane p-3 active cim-tab-line">
          <div class="form-check form-check-inline" v-for="radio in state.radios">
            <input class="form-check-input" type="radio" name="radio.name" :value="radio.value" v-model="state.invoiceType" @change="onRadieChanged(radio.value)">
            <label class="form-check-label">{{ radio.label }}</label>
          </div>
          <div class="form-check form-check-inline">
            <label class="form-check-label cim-invoice-label">开票金额合计：</label>
            <label class="cim-invoice-label cim-invoice-amount"> {{ state.total }} </label>
          </div>
        </div>

        <div class="row g-2 cim-row">
          <div class="col-md">
            <div class="form-floating">
              <input class="form-control" type="email" placeholder="发票抬头">
              <label>发票抬头</label>
            </div>
          </div>
          <div class="col-md">
            <div class="form-floating">
              <input class="form-control" type="email" placeholder="纳税人识别号">
              <label>纳税人识别号</label>
            </div>
          </div>
        </div>

        <div class="row g-2 cim-row">
          <div class="col-md">
            <div class="form-floating">
              <input class="form-control" type="email" placeholder="接收人">
              <label>接收人</label>
            </div>
          </div>
          <div class="col-md">
            <div class="form-floating">
              <input class="form-control" type="email" placeholder="手机号">
              <label>手机号</label>
            </div>
          </div>
        </div>

        <div class="row g-2 cim-row">
          <div class="col-md">
            <div class="form-floating">
              <input class="form-control" type="email" placeholder="邮箱地址">
              <label>{{ state.invoiceType == INVOICE_TYPE.ONLINE ? '邮箱地址' : '收件地址' }}</label>
            </div>
          </div>
          <div class="col-md">
            <div class="form-floating">
              <input class="form-control" type="email" placeholder="开票备注">
              <label>开票备注</label>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

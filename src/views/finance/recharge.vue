<script setup>
import { reactive } from 'vue';
import utils from '../../common/utils';

let TAB_NAME = {
  ONLINE: 'online',
  PUBLIC: 'public'
};
let state = reactive({
  activeTab: TAB_NAME.ONLINE 
});

function onTab(name){
  utils.extend(state, { activeTab: name });
}
</script>
<template>
  <div class="mb-4">
    <ul class="nav nav-underline-border">
      <li class="nav-item"><a class="nav-link cicon cicon-money1" @click="onTab(TAB_NAME.ONLINE)" :class="{'active': utils.isEqual(state.activeTab, TAB_NAME.ONLINE)}">在线充值</a></li>
      <li class="nav-item"><a class="nav-link cicon cicon-money2" @click="onTab(TAB_NAME.PUBLIC)" :class="{'active': utils.isEqual(state.activeTab, TAB_NAME.PUBLIC)}">对公汇款</a></li>
    </ul>
    <div class="tab-content">
      <div class="tab-pane p-3" :class="{'active': utils.isEqual(state.activeTab, TAB_NAME.ONLINE)}">
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">充值账号</label>
          <div class="col-sm-4">
            <div class="form-control-plaintext">10929298490</div>
          </div>
          <div class="col-sm-4"></div>
        </div>
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">账户余额</label>
          <div class="col-sm-4">
            <div class="form-control-plaintext">30293.00 元</div>
          </div>
          <div class="col-sm-4"></div>
        </div>
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">充值金额</label>
          <div class="col-sm-2">
            <div class="form-control-plaintext cim-ra-input">
              <input class="form-control" type="number">
              <span>元</span>
            </div>
          </div>
          <div class="col-sm-4"></div>
        </div>

        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">支付类型</label>
          <div class="col-sm-6">
            <div class="form-control-plaintext cim-ra-paies">
              <div class="cim-ra-pay cicon cicon-wechatpay">微信支付</div>
              <div class="cim-ra-pay cicon cicon-alipay  cim-ra-pay-selected">支付宝</div>
            </div>
          </div>
          <div class="col-sm-4"></div>
        </div>
        
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label"></label>
          <div class="col-sm-6">
            <div class="cim-ra-qr"></div>
          </div>
          <div class="col-sm-4"></div>
        </div>
      </div>

      <div class="tab-pane p-3" :class="{'active': utils.isEqual(state.activeTab, TAB_NAME.PUBLIC)}">
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">收款名称</label>
          <div class="col-sm-4">
            <div class="form-control-plaintext">北京追风未来科技有限公司</div>
          </div>
          <div class="col-sm-4"></div>
        </div>
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">收款账号</label>
          <div class="col-sm-4">
            <div class="form-control-plaintext">110293840383</div>
          </div>
          <div class="col-sm-4"></div>
        </div>
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">收款银行</label>
          <div class="col-sm-4">
            <div class="form-control-plaintext">北京工商银行人民大学支行</div>
          </div>
          <div class="col-sm-4"></div>
        </div>
      </div>
    </div>
  </div>
</template>

import { request } from './request';
import SERVER_PATH from './api';
import utils from '../common/utils';

function getTranslate(params = {}){
  let { app_key } = params;
  let url = `${SERVER_PATH.R3D_TRANSLTE_GET}?app_key=${app_key}`;
  return request(url, {
    method: 'GET'
  });
}

function setTranslate(data){
  return request(SERVER_PATH.R3D_TRANSLTE_SET, {
    method: 'POST',
    body: utils.toJSON(data)
  });
}

function getRTC(params = {}){
  let { app_key } = params;
  let url = `${SERVER_PATH.R3D_RTC_GET}?app_key=${app_key}`;
  return request(url, {
    method: 'GET'
  });
}

function setRTC(data){
  return request(SERVER_PATH.R3D_RTC_SET, {
    method: 'POST',
    body: utils.toJSON(data)
  });
}

export default {
  getTranslate,
  setTranslate,
  getRTC,
  setRTC,
}
import { request } from './request';
import SERVER_PATH from './api';
import utils from '../common/utils';
function ajax(data){
  return request(SERVER_PATH.MISC_REQUEST, {
    method: 'POST',
    body: utils.toJSON(data)
  });
}

export default {
  ajax
}
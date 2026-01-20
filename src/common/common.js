import utils from "./utils";

function getRangeDate(num){
  let dayMs = num * 24 * 60 * 60 * 1000;
  let current = utils.formatTime(Date.now(), 'yyyy-MM-dd');
  let time = new Date(`${current} 00:00:00`).getTime() - dayMs;
  let date = utils.formatTime(time);
  let start = new Date(date).getTime();
  let end = new Date().getTime() - 1 * 24 * 60 * 60 * 1000
  return { start, end }
}

function calcYesterday(result){
  let { upMsgs, downMsgs, disMsgs } = result;
  let upMsg = getStatInfo(upMsgs);
  let downMsg = getStatInfo(downMsgs);
  let disMsg = getStatInfo(disMsgs);
  return { upMsg, downMsg, disMsg };
}

function getStatInfo(items){
  let item = items[0] || { count: 0 }
  let qianitem = items[1] || { count: 1 }
  let percent = (item.count - qianitem.count) / qianitem.count * 100;
  let isUp = false;
  if(percent >= 0){
    isUp = true;
  }
  return { 
    isUp: isUp,
    percent: Math.abs(percent.toFixed(2)),
    count: utils.numberWithCommas(item.count)
  };
}

function formatChatData(result){
  let { upMsgs, downMsgs, disMsgs } = result;
  let dates = utils.map(upMsgs, (item) => {
    return utils.formatTime(item.time_mark, 'yyyy-MM-dd');
  }).reverse();
  
  upMsgs = utils.map(upMsgs, (item) => {
    return item.count;
  }).reverse();
  
  downMsgs = utils.map(downMsgs, (item) => {
    return item.count;
  }).reverse();
  
  disMsgs = utils.map(disMsgs, (item) => {
    return item.count;
  }).reverse();
  return { dates, upMsgs, downMsgs, disMsgs };
}

function formatDauChat(result){
  let { items } = result;
  let dates = utils.map(items, (item) => {
    return utils.formatTime(item.time_mark, 'yyyy-MM-dd');
  });
  
  let daus = utils.map(items, (item) => {
    return item.count;
  });
 
  return { dates, daus };
}
function isValidateUrl(url) {
  if (typeof url !== 'string'){
    return false;
  };
  let strictUrlRegex = /^https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/=]*)$/;
  return strictUrlRegex.test(url);
}
export default {
  getRangeDate,
  calcYesterday,
  formatDauChat,
  formatChatData,
  isValidateUrl,
}
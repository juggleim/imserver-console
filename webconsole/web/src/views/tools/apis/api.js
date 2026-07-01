import { users } from "./user";
import { messages } from './message';
import { conversations } from './conversation';
import { groups } from './group';
import { histories } from './history';
import { sensitivewords } from './sensitive';
import { usertags } from './user-tag';
import { chatrooms } from './chatroom';

let apis = users.concat(messages).concat(conversations)
          .concat(groups)
          .concat(histories)
          .concat(sensitivewords)
          .concat(usertags)
          .concat(chatrooms);
export { apis };
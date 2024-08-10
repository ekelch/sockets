export type Message = {
  username: string;
  message: string;
  clients: Map<string, UserClient>;
};

export type DirectMessage = {
  fromUser: string;
  toUser: string;
  message: string;
};

export type UserClient = {
  connection: any;
  username: string;
};

export const DefaultMessages = {
  CONNECT: "has connected to the chat",
  DISCONNECT: "has disconnected from the chat",
};

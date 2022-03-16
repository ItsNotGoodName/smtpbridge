
const API_URL = import.meta.env.VITE_API_URL
  ? import.meta.env.VITE_API_URL
  : "";

const jsonResponse = (req: Promise<Response>) => req.then((res) => res.json());

export interface IResponse<T> {
  ok: boolean;
  status: number;
  data?: T;
  error?: string;
}

export interface ICursor {
  ascending?: Boolean
  limit?: Number
  cursor?: Number
}

export interface INextCursor {
  next_cursor: Number
  has_more: Boolean
}

export interface IVersion {
  version: String
  commit: String
  date: String
  built_by: String
}

export interface IInfo {
  events_count: Number
  messages_count: Number
  attachments_count: Number
  attachments_size: Number
}

export interface IMessage {
  id: number
  from: string
  to: [string]
  subject: string
  text: string
  attachment: [IAttachment]
  created_at: string
}

export interface IAttachment {
  id: number
  name: string
  file: string
  type: string
}

export interface IEvent {
  id: number
  name: string
  description: string
  created_at: string
}

export interface IMessages {
  cursor: INextCursor
  messages: [IMessage]
}

export interface IEvents {
  cursor: INextCursor
  events: [IEvent]
}

export default {
  getVersion(): Promise<IResponse<IVersion>> {
    return jsonResponse(fetch(API_URL + "/api/version"));
  },
  getInfo(): Promise<IResponse<IInfo>> {
    return jsonResponse(fetch(API_URL + "/api/info"));
  },
  deleteMessage(id: number): Promise<IResponse<undefined>> {
    return jsonResponse(fetch(API_URL + "/api/message/" + id, {
      method: "DELETE"
    }));
  },
  getMessage(id: Number): Promise<IResponse<IMessage>> {
    return jsonResponse(fetch(API_URL + "/api/message/" + id));
  },
  getMessageEvents(id: Number, cursor: ICursor): Promise<IResponse<IEvents>> {
    return jsonResponse(fetch(API_URL + "/api/message/" + id + new URLSearchParams(cursor as any)));
  },
  getMessages(cursor: ICursor): Promise<IResponse<IMessages>> {
    return jsonResponse(fetch(API_URL + "/api/messages" + new URLSearchParams(cursor as any)));
  },
  getEvents(cursor: ICursor): Promise<IResponse<IEvents>> {
    return jsonResponse(fetch(API_URL + "/api/events" + new URLSearchParams(cursor as any)));
  },
  attachmentUrl(attachment: IAttachment): string {
    return API_URL + "/attachment/" + attachment.file
  }
};
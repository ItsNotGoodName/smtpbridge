const API_URL = import.meta.env.VITE_API_URL
  ? import.meta.env.VITE_API_URL
  : "";

const jsonResponse = (req: Promise<Response>) => req.then((res) => res.json());

export interface IResponse<T> {
  ok: boolean;
  status: number;
  data?: T;
  error?: {
    message: string;
  };
}

export interface ICursor {
  ascending?: boolean
  limit?: Number
  cursor?: Number
}

export interface IPage {
  page: number
  limit: number
  ascending: boolean
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
  to: string[]
  subject: string
  text: string
  attachments: IAttachment[]
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
  messages: IMessage[]
  has_back: boolean
  back_cursor: number
  next_cursor: number
  has_next: boolean
}

export interface IEvents {
  max_page: number
  events: IEvent[]
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
  getMessageEvents(id: Number, page: IPage): Promise<IResponse<IEvents>> {
    return jsonResponse(fetch(API_URL + "/api/message/" + id + "/events?" + new URLSearchParams(page as any)));
  },
  getMessages(cursor: ICursor): Promise<IResponse<IMessages>> {
    return jsonResponse(fetch(API_URL + "/api/messages?" + new URLSearchParams(cursor as any)));
  },
  getEvents(page: IPage): Promise<IResponse<IEvents>> {
    return jsonResponse(fetch(API_URL + "/api/events?" + new URLSearchParams(page as any)));
  },
};
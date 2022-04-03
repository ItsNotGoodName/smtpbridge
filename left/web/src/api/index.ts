const API_URL = import.meta.env.VITE_API_URL
  ? import.meta.env.VITE_API_URL
  : "";

const jsonResponse = (req: Promise<Response>) => req.then((res) => res.json());

const request = (url: string, method = "GET") => { return { url: url, method } };

export interface IRequest<T> {
  url: string;
  method: string;
}

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
  limit?: number
  cursor?: number
}

export interface IPage {
  page: number
  limit: number
  ascending: boolean
}

export interface IVersion {
  version: string
  commit: string
  date: string
  built_by: string
}

export interface IInfo {
  events_count: number
  messages_count: number
  attachments_count: number
  attachments_size: number
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
  page: number,
  max_page: number
  max_count: number
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
  getMessageEvents(id: number, page: IPage): Promise<IResponse<IEvents>> {
    return jsonResponse(fetch(API_URL + "/api/message/" + id + "/events?" + new URLSearchParams(page as any)));
  },
  getMessages(cursor: ICursor): Promise<IResponse<IMessages>> {
    return jsonResponse(fetch(API_URL + "/api/messages?" + new URLSearchParams(cursor as any)));
  },
  getEvents(page: IPage): Promise<IResponse<IEvents>> {
    return jsonResponse(fetch(API_URL + "/api/events?" + new URLSearchParams(page as any)));
  },
  messageGet: (id: number | string): IRequest<IMessage> => request(API_URL + "/api/message/" + id),
  messageDelete: (id: number | string): IRequest<IMessage> => request(API_URL + "/api/message/" + id, "DELETE"),
  messageEventsGet: (id: number | string, page: IPage): IRequest<IEvents> => request(API_URL + "/api/message/" + id + "/events?" + new URLSearchParams(page as any)),
  infoGet: (): IRequest<IInfo> => request(API_URL + "/api/info"),
  versionGet: (): IRequest<IVersion> => request(API_URL + "/api/version"),
  eventsGet: (page: IPage): IRequest<IEvents> => request(API_URL + "/api/events?" + new URLSearchParams(page as any)),
};
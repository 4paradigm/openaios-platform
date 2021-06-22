interface ResponseList<T> {
  page: number;
  pageSize: number;
  total: number;
  rows: T;
}

interface ResponseBody<T> {
  content?: any;
  type?: string;
  message: T;
}

export { ResponseBody, ResponseList };

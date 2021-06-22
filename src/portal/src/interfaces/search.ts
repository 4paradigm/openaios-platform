enum SearchDirection {
  DESC = 'DESC',
  ASC = 'ASC',
}

interface SearchRequest {
  direction?: SearchDirection | null;
  orderBy?: string | null;
  page?: number; // 从0开始
  pageSize?: number;
  keyword?: string | null;
  condation?: any;
  type?: string;
}

export { SearchDirection, SearchRequest };

export interface Board {
  id: number;
  slug: string;
  description: string;
}

export interface Topic {
  id: number;
  board_id: number;
  title: string;
  author: string;
  pub_date: string;
  status: string;
  last_post_id?: number;
  post_count: number;
}

export interface Post {
  id: number;
  topic_id: number;
  author: string;
  content: string;
  pub_date: string;
}

export interface BoardWithRecent {
  id: number;
  slug: string;
  description: string;
  recent_topic?: Topic;
  recent_post?: Post;
}

export interface BoardsResponse {
  boards: BoardWithRecent[];
}

export interface PaginationMeta {
  page: number;
  per_page: number;
  total_pages: number;
  total: number;
}

export interface TopicsResponse {
  topics: Topic[];
  pagination: PaginationMeta;
}

export interface PostsResponse {
  posts: Post[];
  topic: Topic;
  pagination: PaginationMeta;
}

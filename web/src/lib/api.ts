import type { BoardsResponse, TopicsResponse, PostsResponse } from "./types";

const API_BASE = "/api";

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string = API_BASE) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {},
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;

    const config: RequestInit = {
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
      ...options,
    };

    const response = await fetch(url, config);

    if (!response.ok) {
      throw new Error(`API Error: ${response.status} ${response.statusText}`);
    }

    return response.json();
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "GET" });
  }

  async post<T>(endpoint: string, data?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: "POST",
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T>(endpoint: string, data?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: "PUT",
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "DELETE" });
  }

  // Health check
  async health() {
    return this.get<{ status: string; message: string }>("/health");
  }

  // Boards
  async getBoards() {
    return this.get<BoardsResponse>("/boards");
  }

  // Topics
  async getTopics(boardSlug: string, page: number = 1, perPage: number = 50) {
    const params = new URLSearchParams({
      page: page.toString(),
      per_page: perPage.toString(),
    });
    return this.get<TopicsResponse>(`/boards/${boardSlug}/topics?${params}`);
  }

  // Posts
  async getPosts(topicId: number, page: number = 1, perPage: number = 50) {
    const params = new URLSearchParams({
      page: page.toString(),
      per_page: perPage.toString(),
    });
    return this.get<PostsResponse>(`/topics/${topicId}/posts?${params}`);
  }
}

export const apiClient = new ApiClient();

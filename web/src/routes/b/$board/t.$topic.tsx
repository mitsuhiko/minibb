import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { apiClient } from "../../../lib/api";
import type { Post, PaginationMeta } from "../../../lib/types";

interface TopicSearchParams {
  page?: number;
}

export const Route = createFileRoute("/b/$board/t/$topic")({
  validateSearch: (search: Record<string, unknown>): TopicSearchParams => {
    return {
      page: Number(search?.page) || 1,
    };
  },
  component: TopicPage,
});

function TopicPage() {
  const { board, topic } = Route.useParams();
  const { page = 1 } = Route.useSearch();
  const navigate = useNavigate();
  const [jumpToPage, setJumpToPage] = useState("");

  const { data, isLoading, error } = useQuery({
    queryKey: ["posts", parseInt(topic), page],
    queryFn: () => apiClient.getPosts(parseInt(topic), page),
  });

  const handlePageJump = (e: React.FormEvent) => {
    e.preventDefault();
    const pageNum = parseInt(jumpToPage);
    if (pageNum > 0 && pageNum <= (data?.pagination.total_pages || 1)) {
      navigate({
        to: "/b/$board/t/$topic",
        params: { board, topic },
        search: { page: pageNum },
      });
      setJumpToPage("");
    }
  };

  const navigateToPage = (newPage: number) => {
    navigate({
      to: "/b/$board/t/$topic",
      params: { board, topic },
      search: { page: newPage },
    });
  };

  if (isLoading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">
          <p className="text-gray-600">Loading posts...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">
          <p className="text-red-600">Error loading posts</p>
        </div>
      </div>
    );
  }

  const posts = data?.posts || [];
  const topicData = data?.topic;
  const pagination = data?.pagination;

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <Link
              to="/b/$board"
              params={{ board }}
              className="text-blue-600 hover:text-blue-800 text-sm"
            >
              ← Back to /{board}/
            </Link>
            {topicData && (
              <h2 className="text-2xl font-bold text-gray-900 mt-2">
                {topicData.title}
              </h2>
            )}
            {topicData && (
              <div className="text-sm text-gray-600 mt-1">
                by {topicData.author} •{" "}
                {new Date(topicData.pub_date).toLocaleDateString()}
              </div>
            )}
          </div>
        </div>

        {pagination && <PaginationInfo pagination={pagination} />}

        {posts.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-6 text-center">
            <p className="text-gray-500">No posts in this topic yet.</p>
          </div>
        ) : (
          <div className="space-y-4">
            {posts.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        )}

        {pagination && pagination.total_pages > 1 && (
          <PaginationControls
            pagination={pagination}
            onNavigate={navigateToPage}
            jumpToPage={jumpToPage}
            setJumpToPage={setJumpToPage}
            onPageJump={handlePageJump}
          />
        )}
      </div>
    </div>
  );
}

function PostCard({ post }: { post: Post }) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  return (
    <div className="bg-white rounded-lg shadow hover:shadow-md transition-shadow p-6">
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center space-x-2">
          <span className="font-medium text-gray-900">{post.author}</span>
          <span className="text-gray-500">•</span>
          <span className="text-sm text-gray-500">
            {formatDate(post.pub_date)}
          </span>
        </div>
        <span className="text-xs text-gray-400">#{post.id}</span>
      </div>
      <div className="prose prose-sm max-w-none">
        <div className="whitespace-pre-wrap text-gray-700">{post.content}</div>
      </div>
    </div>
  );
}

function PaginationInfo({ pagination }: { pagination: PaginationMeta }) {
  return (
    <div className="bg-blue-50 rounded-lg p-4">
      <div className="text-sm text-blue-700">
        Page {pagination.page} of {pagination.total_pages} • {pagination.total}{" "}
        total posts
      </div>
    </div>
  );
}

interface PaginationControlsProps {
  pagination: PaginationMeta;
  onNavigate: (page: number) => void;
  jumpToPage: string;
  setJumpToPage: (page: string) => void;
  onPageJump: (e: React.FormEvent) => void;
}

function PaginationControls({
  pagination,
  onNavigate,
  jumpToPage,
  setJumpToPage,
  onPageJump,
}: PaginationControlsProps) {
  const { page, total_pages } = pagination;

  return (
    <div className="bg-white rounded-lg shadow p-4">
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div className="flex items-center space-x-2">
          <button
            onClick={() => onNavigate(page - 1)}
            disabled={page <= 1}
            className="px-3 py-2 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            ← Previous
          </button>
          <button
            onClick={() => onNavigate(page + 1)}
            disabled={page >= total_pages}
            className="px-3 py-2 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next →
          </button>
        </div>

        <div className="text-sm text-gray-600">
          Page {page} of {total_pages}
        </div>

        <form onSubmit={onPageJump} className="flex items-center space-x-2">
          <label htmlFor="jump-to-page" className="text-sm text-gray-600">
            Jump to:
          </label>
          <input
            id="jump-to-page"
            type="number"
            min="1"
            max={total_pages}
            value={jumpToPage}
            onChange={(e) => setJumpToPage(e.target.value)}
            placeholder="Page"
            className="w-20 px-2 py-1 text-sm border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            type="submit"
            className="px-3 py-1 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Go
          </button>
        </form>
      </div>
    </div>
  );
}

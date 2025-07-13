import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { apiClient } from "../../../lib/api";
import type { Topic, PaginationMeta } from "../../../lib/types";

interface BoardSearchParams {
  page?: number;
}

export const Route = createFileRoute("/b/$board/")({
  validateSearch: (search: Record<string, unknown>): BoardSearchParams => {
    return {
      page: Number(search?.page) || 1,
    };
  },
  component: BoardPage,
});

function BoardPage() {
  const { board } = Route.useParams();
  const { page = 1 } = Route.useSearch();
  const navigate = useNavigate();
  const [jumpToPage, setJumpToPage] = useState("");

  const { data, isLoading, error } = useQuery({
    queryKey: ["topics", board, page],
    queryFn: () => apiClient.getTopics(board, page),
  });

  const handlePageJump = (e: React.FormEvent) => {
    e.preventDefault();
    const pageNum = parseInt(jumpToPage);
    if (pageNum > 0 && pageNum <= (data?.pagination.total_pages || 1)) {
      navigate({
        to: "/b/$board",
        params: { board },
        search: { page: pageNum },
      });
      setJumpToPage("");
    }
  };

  const navigateToPage = (newPage: number) => {
    navigate({
      to: "/b/$board",
      params: { board },
      search: { page: newPage },
    });
  };

  if (isLoading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">
          <p className="text-gray-600">Loading topics...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">
          <p className="text-red-600">Error loading topics</p>
        </div>
      </div>
    );
  }

  const topics = data?.topics || [];
  const pagination = data?.pagination;

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <Link to="/" className="text-blue-600 hover:text-blue-800 text-sm">
              ← Back to boards
            </Link>
            <h2 className="text-2xl font-bold text-gray-900 mt-2">/{board}/</h2>
          </div>
        </div>

        {pagination && <PaginationInfo pagination={pagination} />}

        {topics.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-6 text-center">
            <p className="text-gray-500">No topics in this board yet.</p>
          </div>
        ) : (
          <div className="space-y-2">
            {topics.map((topic) => (
              <TopicRow key={topic.id} topic={topic} />
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

function TopicRow({ topic }: { topic: Topic }) {
  const { board } = Route.useParams();

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  return (
    <Link
      to="/b/$board/t/$topic"
      params={{ board, topic: topic.id.toString() }}
      className="block bg-white rounded border hover:bg-gray-50 transition-colors p-4"
    >
      <div className="flex justify-between items-center">
        <div className="flex-1">
          <h3 className="text-lg font-medium text-blue-600 hover:text-blue-800">
            {topic.title}
          </h3>
          <div className="text-sm text-gray-600 mt-1">
            by {topic.author} • {formatDate(topic.pub_date)}
          </div>
        </div>
        <div className="text-right text-sm text-gray-500">
          {topic.post_count} posts
        </div>
      </div>
    </Link>
  );
}

function PaginationInfo({ pagination }: { pagination: PaginationMeta }) {
  return (
    <div className="bg-blue-50 rounded-lg p-4">
      <div className="text-sm text-blue-700">
        Page {pagination.page} of {pagination.total_pages} • {pagination.total}{" "}
        total topics
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

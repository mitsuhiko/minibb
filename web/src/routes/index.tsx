import { createFileRoute, Link } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "../lib/api";
import type { BoardWithRecent } from "../lib/types";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const {
    data: boardsData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["boards"],
    queryFn: () => apiClient.getBoards(),
  });

  if (isLoading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">
          <p className="text-gray-600">Loading boards...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="text-center">
          <p className="text-red-600">Error loading boards</p>
        </div>
      </div>
    );
  }

  const boards = boardsData?.boards || [];

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="space-y-6">
        <h3 className="text-xl font-semibold text-gray-900">Boards</h3>

        {boards.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-6 text-center">
            <p className="text-gray-500">
              No boards available yet. Check back later!
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {boards.map((board) => (
              <BoardCard key={board.id} board={board} />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function BoardCard({ board }: { board: BoardWithRecent }) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  return (
    <div className="bg-white rounded-lg shadow hover:shadow-md transition-shadow p-6">
      <div className="flex justify-between items-start">
        <div className="flex-1">
          <Link
            to="/b/$board"
            params={{ board: board.slug }}
            className="text-lg font-semibold text-blue-600 hover:text-blue-800 block"
          >
            /{board.slug}/
          </Link>
          <p className="text-gray-600 text-sm mt-1">{board.description}</p>
        </div>

        {board.recent_topic ? (
          <div className="ml-6 text-right">
            <p className="text-sm font-medium text-gray-800 truncate max-w-xs">
              {board.recent_topic.title}
            </p>
            <div className="text-xs text-gray-500 mt-1">
              by {board.recent_topic.author}
            </div>
            <div className="text-xs text-gray-500">
              {formatDate(board.recent_topic.pub_date)} •{" "}
              {board.recent_topic.post_count} posts
            </div>
          </div>
        ) : (
          <div className="ml-6 text-right text-gray-500 text-sm italic">
            No topics yet
          </div>
        )}
      </div>
    </div>
  );
}

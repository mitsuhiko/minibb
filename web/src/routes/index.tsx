import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="text-center">
        <h2 className="text-2xl font-semibold text-gray-900 mb-4">
          Welcome to MiniBB
        </h2>
        <p className="text-gray-600 mb-8">A simple bulletin board system</p>
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">Boards</h3>
          <p className="text-gray-500">
            No boards available yet. Check back later!
          </p>
        </div>
      </div>
    </div>
  );
}

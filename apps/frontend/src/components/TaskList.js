import React from 'react';

const TaskList = ({ tasks, loading }) => {
    if (loading) {
        return <div className="text-center py-4">Loading tasks...</div>;
    }

    if (!tasks || tasks.length === 0) {
        return <div className="text-center py-4 text-gray-500">No tasks found. Create one above!</div>;
    }

    const getStatusColor = (status) => {
        switch (status) {
            case 'COMPLETED': return 'bg-green-100 text-green-800';
            case 'PROCESSING': return 'bg-blue-100 text-blue-800';
            case 'FAILED': return 'bg-red-100 text-red-800';
            default: return 'bg-yellow-100 text-yellow-800';
        }
    };

    return (
        <div className="bg-white shadow overflow-hidden sm:rounded-md">
            <ul className="divide-y divide-gray-200">
                {tasks.map((task) => (
                    <li key={task.id || task._id}>
                        <div className="px-4 py-4 sm:px-6">
                            <div className="flex items-center justify-between">
                                <div className="text-sm font-medium text-indigo-600 truncate">
                                    Task ID: {task.id || task._id}
                                </div>
                                <div className="ml-2 flex-shrink-0 flex">
                                    <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getStatusColor(task.status)}`}>
                                        {task.status}
                                    </span>
                                </div>
                            </div>
                            <div className="mt-2 sm:flex sm:justify-between">
                                <div className="sm:flex">
                                    <p className="flex items-center text-sm text-gray-500">
                                        Original: <a href={task.image_url} target="_blank" rel="noopener noreferrer" className="ml-1 text-indigo-500 hover:text-indigo-700 truncate max-w-xs">{task.image_url}</a>
                                    </p>
                                </div>
                                <div className="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
                                    <p>
                                        Created: {new Date(task.created_at).toLocaleString()}
                                    </p>
                                </div>
                            </div>
                            {task.processed_url && (
                                <div className="mt-2">
                                    <p className="text-sm text-gray-500">
                                        Processed: <a href={task.processed_url} target="_blank" rel="noopener noreferrer" className="text-green-600 hover:text-green-800 font-medium">{task.processed_url}</a>
                                    </p>
                                </div>
                            )}
                        </div>
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default TaskList;

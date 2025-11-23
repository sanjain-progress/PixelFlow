import React, { useState } from 'react';
import { taskService } from '../services/api';

const UploadForm = ({ onTaskCreated }) => {
    const [imageUrl, setImageUrl] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setError('');
        setSuccess('');

        try {
            await taskService.upload(imageUrl);
            setSuccess('Task created successfully!');
            setImageUrl('');
            if (onTaskCreated) onTaskCreated();
        } catch (err) {
            setError('Failed to create task. Please try again.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="bg-white shadow sm:rounded-lg p-6 mb-6">
            <h3 className="text-lg leading-6 font-medium text-gray-900 mb-4">
                Create New Task
            </h3>
            <form onSubmit={handleSubmit} className="space-y-4">
                {error && (
                    <div className="text-red-600 text-sm">{error}</div>
                )}
                {success && (
                    <div className="text-green-600 text-sm">{success}</div>
                )}
                <div>
                    <label htmlFor="image_url" className="block text-sm font-medium text-gray-700">
                        Image URL
                    </label>
                    <div className="mt-1 flex rounded-md shadow-sm">
                        <input
                            type="url"
                            name="image_url"
                            id="image_url"
                            className="flex-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full min-w-0 rounded-md sm:text-sm border-gray-300 p-2 border"
                            placeholder="https://example.com/image.jpg"
                            value={imageUrl}
                            onChange={(e) => setImageUrl(e.target.value)}
                            required
                        />
                        <button
                            type="submit"
                            disabled={loading}
                            className="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
                        >
                            {loading ? 'Creating...' : 'Create Task'}
                        </button>
                    </div>
                </div>
            </form>
        </div>
    );
};

export default UploadForm;

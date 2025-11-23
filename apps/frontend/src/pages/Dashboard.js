import React, { useState, useEffect } from 'react';
import UploadForm from '../components/UploadForm';
import TaskList from '../components/TaskList';
import { taskService } from '../services/api';

const Dashboard = () => {
    const [tasks, setTasks] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetchTasks = async () => {
        try {
            const data = await taskService.getAll();
            setTasks(data);
        } catch (error) {
            console.error('Failed to fetch tasks', error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchTasks();
        // Poll for updates every 3 seconds
        const interval = setInterval(fetchTasks, 3000);
        return () => clearInterval(interval);
    }, []);

    return (
        <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
            <div className="px-4 py-6 sm:px-0">
                <h1 className="text-3xl font-bold text-gray-900 mb-8">Dashboard</h1>

                <UploadForm onTaskCreated={fetchTasks} />

                <div className="mt-8">
                    <h3 className="text-lg leading-6 font-medium text-gray-900 mb-4">
                        Your Tasks
                    </h3>
                    <TaskList tasks={tasks} loading={loading} />
                </div>
            </div>
        </div>
    );
};

export default Dashboard;

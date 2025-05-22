// src/components/admin/UsersList.tsx
import { useEffect, useState } from 'react';
import { getUsers } from '../../api/client';
import { useAuth } from '../../context/AuthContext';
import { User } from '../../types'

export const UsersList = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const { token } = useAuth();

  useEffect(() => {
    if (!token) return;
    
    getUsers(token)
      .then(setUsers)
      .finally(() => setLoading(false));
  }, [token]);

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <h2>Users</h2>
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Role</th>
          </tr>
        </thead>
        <tbody>
          {users.map(user => (
            <tr key={user.id}>
              <td>{user.id}</td>
              <td>{user.name}</td>
              <td>{user.role}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
const API_URL = 'http://localhost:8080';

export const artistService = {
    async getById(id) {
        const response = await fetch(`${API_URL}/artists/${id}`);
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async getAll() {
        // According to instructions: GET /artists/all
        const response = await fetch(`${API_URL}/artists/all`);
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async getPaginated(params = { page: 1, limit: 10, name: '', genre: '', country: '' }) {
        const query = new URLSearchParams(params).toString();
        const response = await fetch(`${API_URL}/artists?${query}`);
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json(); // Returns paginated format
    },

    async create(artistData) {
        const response = await fetch(`${API_URL}/artists`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(artistData)
        });
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async update(id, artistData) {
        const response = await fetch(`${API_URL}/artists/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(artistData)
        });
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async delete(id) {
        const response = await fetch(`${API_URL}/artists/${id}`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Network response was not ok');
        /* Some APIs return 204 No Content for DELETE, so handle carefully */
        if (response.status === 204) return { data: { message: 'Deleted successfully' } };
        return response.json();
    }
};

const API_URL = 'http://localhost:8080';

export const albumService = {
    async getById(id) {
        const response = await fetch(`${API_URL}/albums/${id}`);
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async getByArtistId(artistId) {
        const response = await fetch(`${API_URL}/albums/artist/${artistId}`);
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async getPaginated(params = { page: 1, limit: 10, artist_name: '', title: '' }) {
        const query = new URLSearchParams(params).toString();
        const response = await fetch(`${API_URL}/albums?${query}`);
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async create(albumData) {
        const response = await fetch(`${API_URL}/albums`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(albumData)
        });
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async update(id, albumData) {
        const response = await fetch(`${API_URL}/albums/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(albumData)
        });
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async delete(id) {
        const response = await fetch(`${API_URL}/albums/${id}`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Network response was not ok');
        if (response.status === 204) return { data: { message: 'Deleted successfully' } };
        return response.json();
    },

    async addTrack(albumId, trackData) {
        const response = await fetch(`${API_URL}/albums/${albumId}/tracks`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(trackData)
        });
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async removeTrack(albumId, songId) {
        const response = await fetch(`${API_URL}/albums/${albumId}/tracks/${songId}`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Network response was not ok');
        if (response.status === 204) return { data: { message: 'Deleted successfully' } };
        return response.json();
    }
};

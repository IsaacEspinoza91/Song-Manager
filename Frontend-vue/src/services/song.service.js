const API_URL = 'http://localhost:8080';

export const songService = {
    async getById(id) {
        const response = await fetch(`${API_URL}/songs/${id}`);
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async getAll() {
        const response = await fetch(`${API_URL}/songs/all`);
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async getPaginated(params = { page: 1, limit: 10, title: '', artist_id: '', artist_name: '' }) {
        const query = new URLSearchParams(params).toString();
        const response = await fetch(`${API_URL}/songs?${query}`);
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json(); // Returns paginated format
    },

    async create(songData) {
        const response = await fetch(`${API_URL}/songs`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(songData)
        });
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async update(id, songData) {
        const response = await fetch(`${API_URL}/songs/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(songData)
        });
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async delete(id) {
        const response = await fetch(`${API_URL}/songs/${id}`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Network response was not ok');
        if (response.status === 204) return { data: { message: 'Deleted successfully' } };
        return response.json();
    },

    async addArtistToSong(songId, artistData) {
        // artistData format: { "artist_id": 3, "role": "main" }
        const response = await fetch(`${API_URL}/songs/${songId}/artist`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(artistData)
        });
        if (!response.ok) throw new Error('Network response was not ok');
        return response.json();
    },

    async removeArtistFromSong(songId, artistId) {
        const response = await fetch(`${API_URL}/songs/${songId}/artist/${artistId}`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Network response was not ok');
        if (response.status === 204) return { data: { message: 'Deleted successfully' } };
        return response.json();
    }
};

// mockData.js
// Provides static mock data matching the API structures.

export const mockArtists = [
    {
        "id": 1,
        "name": "Bad Bunny",
        "genre": "Reggaeton/Trap",
        "country": "Puerto Rico",
        "bio": "Conocido como el Conejo Malo, uno de los mayores exponentes del género urbano.",
        "image_url": "https://i.scdn.co/image/ab6761610000e5eb9e2f95ea7ee4fa3bc8b9a117",
        "created_at": "2026-02-24T20:55:29.816026Z",
        "updated_at": "2026-02-24T20:55:29.816026Z"
    },
    {
        "id": 2,
        "name": "Gustavo Cerati",
        "genre": "Rock en Español",
        "country": "Argentina",
        "bio": "Líder de Soda Stereo y uno de los músicos más influyentes del rock latinoamericano.",
        "image_url": "https://i.scdn.co/image/ab6761610000e5ebcb2aae276d3f2256da6c8577",
        "created_at": "2026-02-24T20:55:29.816026Z",
        "updated_at": "2026-02-24T20:55:29.816026Z"
    },
    {
        "id": 3,
        "name": "Mon Laferte",
        "genre": "Pop/Rock Alternativo",
        "country": "Chile",
        "bio": "Cantante y compositora chilena con éxito internacional.",
        "image_url": "https://i.scdn.co/image/ab6761610000e5ebb7af01f807df9e7ae71e4288",
        "created_at": "2026-02-24T20:55:29.816026Z",
        "updated_at": "2026-02-24T20:55:29.816026Z"
    },
    {
        "id": 7,
        "name": "WOS",
        "genre": "Rap/Hip Hop",
        "country": "Argentina",
        "bio": "Rapero y freestyler argentino.",
        "image_url": "https://i.scdn.co/image/ab6761610000e5ebeb1be79b2ced568c07e997f0",
        "created_at": "2026-02-24T20:55:29.816026Z",
        "updated_at": "2026-02-24T20:55:29.816026Z"
    }
];

export const mockSongs = [
    {
        "id": 1,
        "title": "Tití Me Preguntó",
        "duration": 243,
        "created_at": "2026-02-24T20:55:29.816026Z",
        "updated_at": "2026-02-24T20:55:29.816026Z",
        "artists": [
            { "id": 1, "name": "Bad Bunny", "role": "main" }
        ]
    },
    {
        "id": 2,
        "title": "Me Porto Bonito",
        "duration": 178,
        "created_at": "2026-02-24T20:55:29.816026Z",
        "updated_at": "2026-02-24T20:55:29.816026Z",
        "artists": [
            { "id": 1, "name": "Bad Bunny", "role": "main" }
        ]
    },
    {
        "id": 8,
        "title": "Arrancarmelo",
        "duration": 201,
        "created_at": "2026-02-24T21:00:28.736434Z",
        "updated_at": "2026-02-24T21:05:34.048721Z",
        "artists": [
            { "id": 3, "name": "Mon Laferte", "role": "ft" },
            { "id": 7, "name": "WOS", "role": "main" }
        ]
    }
];

export const mockAlbums = [
    {
        "id": 11,
        "title": "Debi tirar mas fotos",
        "release_date": "2025-01-05T00:00:00Z",
        "type": "LP",
        "cover_url": "https://i.scdn.co/image/ab67616d0000b273752e505ed61329596af42d87",
        "created_at": "2026-02-24T23:02:11.650518Z",
        "updated_at": "2026-02-24T23:02:11.650518Z",
        "artists": [
            { "id": 1, "name": "Bad Bunny", "is_primary": true }
        ],
        "tracks": [
            { "track_number": 1, "song_id": 1, "title": "Tití Me Preguntó", "duration": 243 },
            { "track_number": 2, "song_id": 2, "title": "Me Porto Bonito", "duration": 178 }
        ]
    },
    {
        "id": 1,
        "title": "Un Verano Sin Ti",
        "release_date": "2022-05-06T00:00:00Z",
        "type": "Album",
        "cover_url": "https://i.scdn.co/image/ab67616d0000b27349d694203245f241a1bcaa72",
        "created_at": "2026-02-24T20:55:29.816026Z",
        "updated_at": "2026-02-24T20:55:29.816026Z",
        "artists": [
            { "id": 1, "name": "Bad Bunny", "is_primary": true }
        ],
        "tracks": []
    }
];

import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: HomeView
        },
        {
            path: '/artists',
            name: 'artists',
            component: () => import('../views/artists/ArtistListView.vue')
        },
        // {
        //   path: '/artists/:id',
        //   name: 'artist-detail',
        //   component: () => import('../views/artists/ArtistDetailView.vue')
        // },
        {
            path: '/songs',
            name: 'songs',
            component: () => import('../views/songs/SongListView.vue')
        },
        {
            path: '/albums',
            name: 'albums',
            component: () => import('../views/albums/AlbumListView.vue')
        }
    ]
})

export default router

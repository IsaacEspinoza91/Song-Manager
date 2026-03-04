import { ref } from 'vue';

const toasts = ref([]);
let nextId = 0;

export function useToast() {
    const addToast = (message, type = 'success', duration = 5000) => {
        const id = nextId++;
        toasts.value.push({
            id,
            message,
            type
        });

        if (duration > 0) {
            setTimeout(() => {
                removeToast(id);
            }, duration);
        }
        return id;
    };

    const removeToast = (id) => {
        toasts.value = toasts.value.filter(t => t.id !== id);
    };

    const success = (message, duration = 3000) => {
        addToast(message, 'success', duration);
    };

    const error = (message, duration = 5000) => {
        addToast(message, 'error', duration);
    };

    const handleApiError = (err, fallbackMessage = 'Ha ocurrido un error inesperado.') => {
        console.error(err);
        if (err.response && err.response.data) {
            const data = err.response.data;
            let errorMessage = data.message || fallbackMessage;

            if (data.details && typeof data.details === 'object') {
                const detailsList = Object.values(data.details).join('. ');
                if (detailsList) {
                    errorMessage += `: ${detailsList}`;
                }
            }

            error(errorMessage);
        } else {
            error(fallbackMessage);
        }
    };

    return {
        toasts,
        addToast,
        removeToast,
        success,
        error,
        handleApiError
    };
}

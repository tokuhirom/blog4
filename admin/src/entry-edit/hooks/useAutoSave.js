import { useRef, useCallback } from 'preact/hooks';

export function useAutoSave(saveFn, delay = 500) {
    const timerRef = useRef(null);

    const trigger = useCallback((...args) => {
        if (timerRef.current) {
            clearTimeout(timerRef.current);
        }
        timerRef.current = setTimeout(() => {
            timerRef.current = null;
            saveFn(...args);
        }, delay);
    }, [saveFn, delay]);

    return trigger;
}

<script lang="ts">
    import {onMount} from "svelte";

    onMount(() => {
        fetch(import.meta.env.VITE_API_BASE_URL + "/entries")
            .then((res) => res.json())
            .then((data) => {
                console.log(data);
            });
    });
</script>

<div>
</div>
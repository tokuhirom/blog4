<script lang="ts">
    import { onMount } from "svelte";
    import {createAdminApiClient} from "./admin_api";
    import type {GetLatestEntriesRow} from "./generated-client";

    let { path } = $props();
    const api = createAdminApiClient();
    let entry: GetLatestEntriesRow = $state({ });

    onMount(async () => {
        entry = await api.getEntryByDynamicPath({
            path: path
        });
    });
</script>

<div>
    <div>{entry.title}</div>
    <div>{entry.body}</div>
    <div>{entry.format}</div>
    <!-- Additional content and logic for displaying the entry -->
</div>

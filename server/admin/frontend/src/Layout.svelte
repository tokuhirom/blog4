<script>
import AdminHeader from "./components/AdminHeader.svelte";
</script>

<div>
    <AdminHeader />

    <main>
        <!-- This is where the routed components will be rendered -->
        <slot></slot>
    </main>
</div>

<style>
    main {
        margin-top: 4rem;
        margin-bottom: 4rem; /* to avoid overlap with footer */
    }
</style>


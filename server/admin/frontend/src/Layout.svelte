<script>
import {Link} from "svelte-routing";
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
    /* Add any styles specific to your layout here */
    header {
        background-color: #f8f9fa;
        padding: 1rem;
        text-align: center;
    }

    main {
        padding: 4rem;
        margin-bottom: 4rem; /* to avoid overlap with footer */
    }
</style>


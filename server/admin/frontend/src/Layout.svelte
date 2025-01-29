<script>
import {Link} from "svelte-routing";
</script>

<div>
    <header>
        <h1>tokuhirom's blog admin</h1>
        <nav>
            <a href="/">Home</a> | <Link to="/admin/help">Help</Link>
        </nav>
    </header>

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
        padding: 2rem;
        margin-bottom: 4rem; /* to avoid overlap with footer */
    }
</style>


<script lang="ts">
    import {Link, Route, Router} from "svelte-routing";
    import Layout from "./Layout.svelte";
    import NotFound from "./NotFound.svelte";
</script>

<Router>
    <Layout>
        <Route path="/admin/">
            top...
            <Link to="/admin/help">Posts</Link>
        </Route>
        <Route path="/admin/help">
            help...
        </Route>
        <Route path="/*" component={NotFound} />
    </Layout>
</Router>

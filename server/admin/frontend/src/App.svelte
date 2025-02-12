<script lang="ts">
    import {Link, Route, Router} from "svelte-routing";
    import Layout from "./Layout.svelte";
    import NotFound from "./NotFound.svelte";
    import TopPage from "./TopPage.svelte";
</script>

<Router>
    <Layout>
        <Route path="/admin/">
            top...
            {import.meta.env.VITE_API_BASE_URL}
            <Link to="/admin/help">Posts</Link>
            <TopPage />
        </Route>
        <Route path="/admin/help">
            help...
        </Route>
        <Route path="/*" component={NotFound} />
    </Layout>
</Router>

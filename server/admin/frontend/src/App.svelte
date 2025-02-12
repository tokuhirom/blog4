<script lang="ts">
    import {Link, Route, Router} from "svelte-routing";
    import Layout from "./Layout.svelte";
    import NotFound from "./NotFound.svelte";
    import TopPage from "./TopPage.svelte";
</script>

<Router>
    <Layout>
        <Route path="/admin/">
            <Link to="/admin/help">Posts</Link>
            <TopPage />
        </Route>
        <Route path="/admin/help">
            help...
        </Route>
        <Route path="/*" component={NotFound} />
    </Layout>
</Router>

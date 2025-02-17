<script lang="ts">
import { createAdminApiClient } from "../admin_api";

const api = createAdminApiClient();

async function handleNewEntry() {
	try {
		const data = await api.createEntry({ createEntryRequest: {} });
		console.log(`New entry created: ${data.path}`);
		location.href = `/admin/entry/${data.path}`;
	} catch (e) {
		console.error("Error creating new entry:", e);
		alert("Failed to create new entry");
	}
}
</script>

<header>
	<div class="container">
		<!-- トップページへのリンク -->
		<a href="/admin/" class="text-xl">Blog Admin</a>

		<!-- ナビゲーションリンク -->
		<nav>
			<button onclick={handleNewEntry}>New Entry</button>
		</nav>
	</div>
</header>

<style>
	header {
		position: fixed;
		left: 0;
		top: 0;
		z-index: 10;
		width: 100%;
		background-color: #d92706;
		color: white;
		height: 62px;
		vertical-align: middle;

		font-family: 'Hiragino Kaku Gothic ProN', 'Meiryo', sans-serif;
	}

	.container {
		max-width: 1200px;
		margin: 0 auto;
	}

	a {
		display: block;
		float: left;
		text-decoration: none;
		color: white;
		padding: 8px;
	}

	a:hover {
		text-decoration: underline;
		cursor: pointer;
	}

	.text-xl {
		font-size: 1.25rem;
		font-weight: bold;
	}

	nav {
		float: right;
	}

	button {
		background: none;
		border: none;
		color: white;
		text-decoration: none;
		cursor: pointer;
		font: inherit;
	}

	button:hover {
		text-decoration: underline;
	}
</style>

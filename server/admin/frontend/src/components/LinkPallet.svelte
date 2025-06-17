<script lang="ts">
import { createAdminApiClient } from "../admin_api";
import type { LinkPalletData } from "../generated-client/model";
import AdminEntryCardItem from "./AdminEntryCardItem.svelte";
import CardItem from "./CardItem.svelte";

const { linkPallet }: { linkPallet: LinkPalletData } = $props();

const api = createAdminApiClient();

function createNewEntry(title: string) {
	api
		.createEntry({
			createEntryRequest: { title },
		})
		.then((data) => {
			location.href = `/admin/entry/${data.Path}`;
		})
		.catch((err) => {
			console.error("Error creating new entry:", err);
			alert("Failed to create new entry");
		});
}
</script>

<div class="link-pallet">
	<div class="one-hop-link">
		{#each linkPallet.links as link}
			<AdminEntryCardItem entry={link} />
		{/each}
	</div>
	{#each linkPallet.twohops as twohops}
		<div class="two-hop-link">
			{#if twohops.src.Title}
				<AdminEntryCardItem entry={twohops.src} backgroundColor={'yellowgreen'} />
			{:else}
				<CardItem
					onClick={() => createNewEntry(twohops.src.dstTitle)}
					title={twohops.src.dstTitle}
					content=""
					backgroundColor="#c0f6f6"
					color="gray"
				/>
			{/if}
			{#each twohops.links as link}
				<AdminEntryCardItem entry={link} />
			{/each}
		</div>
	{/each}
	{#if linkPallet.newLinks.length > 0}
		<div class="one-hop-link">
			<CardItem onClick={() => false} title="New Item" content="" backgroundColor="darkgoldenrod" />
			{#each linkPallet.newLinks as title}
				<CardItem onClick={() => createNewEntry(title)} {title} content="" color="gray" />
			{/each}
		</div>
	{/if}
</div>

<style>
	.one-hop-link {
		display: flex;
		flex-wrap: wrap;
		gap: 1rem;
		clear: both;
	}
	.two-hop-link {
		display: flex;
		flex-wrap: wrap;
		gap: 1rem;
		clear: both;
		margin-top: 1rem;
	}
</style>

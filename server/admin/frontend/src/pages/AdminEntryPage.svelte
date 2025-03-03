<script lang="ts">
import { onMount } from "svelte";
import MarkdownEditor from "../components/MarkdownEditor.svelte";

import { createAdminApiClient } from "../admin_api";
import {
	type GetLatestEntriesRow,
	type LinkPalletData,
	ResponseError,
} from "../generated-client";
import { extractLinks } from "../extractLinks";
import { debounce } from "../utils";
import LinkPallet from "../components/LinkPallet.svelte";

const { path } = $props();
const api = createAdminApiClient();
let entry: GetLatestEntriesRow = $state({});

let links: { [key: string]: string | null } = $state({});

let title: string = $state("");
let body: string = $state("");
let visibility = $state("private");
let currentLinks: string[] = $state([]);

let linkPallet: LinkPalletData = $state({
	links: [],
	twohops: [],
	newLinks: [],
});

onMount(async () => {
	try {
		entry = await api.getEntryByDynamicPath({
			path: path,
		});

		title = entry.title;
		body = entry.body;
		visibility = entry.visibility;
		currentLinks = extractLinks(entry.body);
	} catch (e) {
		console.error("Failed to get entry:", e);
		if (e instanceof ResponseError) {
			if (e.response.status === 404) {
				// maybe entry was deleted
				location.href = "/admin/";
			}
		}
	}

	links = await api.getLinkedEntryPaths({ path });
	loadLinks();
});

let pageTitles: string[] = $state([]);

let isDirty = false;

function loadLinks() {
	api
		.getLinkPallet({ path })
		.then((data) => {
			console.log("Got link pallet data", data);
			linkPallet = data;
		})
		.catch((error) => {
			console.error("Failed to get links:", error);
		});
}

let message = $state("");
let messageType: "success" | "error" | "" = $state("");

let updatedMessage = $state("");

function showUpdatedMessage(text: string) {
	updatedMessage = text;
	setTimeout(() => {
		updatedMessage = "";
	}, 1000);
}

function clearMessage() {
	message = "";
	messageType = "";
}

function showMessage(type: "success" | "error", text: string) {
	messageType = type;
	message = text;
	setTimeout(() => {
		message = "";
		messageType = "";
	}, 5000); // Hide message after 5 seconds
}

async function handleDelete(event: Event) {
	event.preventDefault();

	const confirmed = confirm(
		`Are you sure you want to delete the entry "${title}"?`,
	);
	if (confirmed) {
		clearMessage();

		try {
			await api.deleteEntry({
				path: entry.path,
			});
			showMessage("success", "Entry deleted successfully");
			location.href = "/admin/";
		} catch (e) {
			console.log(e);
			showMessage("error", "Failed to delete entry");
		}
	}
}

async function handleRegenerateEntryImage(event: Event) {
    event.preventDefault();

    clearMessage();

    try {
        await api.regenerateEntryImage({
            "path": entry.path,
        });
        showMessage("success", "Entry image regenerated successfully");
        location.href = "/admin/";
    } catch (e) {
        console.log(e);
        showMessage("error", "Failed to regenerate entry image");
    }
}

async function handleUpdateBody() {
	clearMessage();

	if (body === "") {
		showMessage("error", "Body cannot be empty");
		return;
	}

	try {
		await api.updateEntryBody({
			path: path,
			updateEntryBodyRequest: {
				body: body,
			},
		});

		showUpdatedMessage("Updated");
		isDirty = false; // Reset dirty flag on successful update
	} catch (e) {
		showMessage("error", "Failed to update entry body");
		console.error("Failed to update entry body:", e);
	}
}

async function handleUpdateTitle() {
	clearMessage();
	if (title === "") {
		showMessage("error", "Title cannot be empty");
		return;
	}

	try {
		await api.updateEntryTitle({
			path: path,
			updateEntryTitleRequest: {
				title,
			},
		});

		showMessage("success", "Entry updated successfully");
		isDirty = false; // Reset dirty flag on successful update
	} catch (e) {
		showMessage("error", "Failed to update entry title");
		console.error("Failed to update entry title:", e);
	}
}

async function createNewEntry(title: string): Promise<void> {
	clearMessage();

	try {
		const data = await api.createEntry({
			createEntryRequest: {
				title,
			},
		});
		location.href = `/admin/entry/${data.path}`;
	} catch (error) {
		console.error("Failed to create new entry:", error);
		showMessage(
			"error",
			`Failed to create new entry: ${error instanceof Error ? error.message : error}`,
		);
	}
}

// デバウンスした自動保存関数
const debouncedUpdateBody = debounce(() => {
	handleUpdateBody();
}, 800);

// title related.
const debouncedTitleUpdate = debounce(() => {
	handleUpdateTitle();
}, 500);

function handleTitleInput() {
	isDirty = true;
	debouncedTitleUpdate();
}

// 入力イベントや変更イベントにデバウンスされた関数をバインド
function handleInputBody() {
	isDirty = true;
	debouncedUpdateBody();

	const newLinks = extractLinks(body);
	if (currentLinks !== newLinks) {
		currentLinks = newLinks;
		loadLinks();
	}
}

function toggleVisibility(event: Event) {
	event.preventDefault();
	event.stopPropagation();

	const newVisibility = visibility === "private" ? "public" : "private";

	if (
		!confirm("Are you sure you want to change the visibility of this entry?")
	) {
		return;
	}

	console.log("Updating visibility to", newVisibility);

	api
		.updateEntryVisibility({
			path: entry.path,
			updateVisibilityRequest: {
				visibility: newVisibility,
			},
		})
		.then((data) => {
			visibility = data.visibility;
		})
		.catch((error) => {
			console.error("Failed to update visibility:", error);
			showMessage("error", `Failed to update visibility: ${error.message}`);
		});
}

// TODO: 未保存の変更がある場合に警告を表示する
// beforeNavigate(({ cancel }) => {
//     if (isDirty && !confirm('You have unsaved changes. Are you sure you want to leave?')) {
//         cancel();
//     }
// });

function getEditDistance(a: string, b: string): number {
	const dp = Array.from({ length: a.length + 1 }, () =>
		Array(b.length + 1).fill(0),
	);

	for (let i = 0; i <= a.length; i++) dp[i][0] = i;
	for (let j = 0; j <= b.length; j++) dp[0][j] = j;

	for (let i = 1; i <= a.length; i++) {
		for (let j = 1; j <= b.length; j++) {
			if (a[i - 1] === b[j - 1]) {
				dp[i][j] = dp[i - 1][j - 1];
			} else {
				dp[i][j] =
					Math.min(
						dp[i - 1][j], // 削除
						dp[i][j - 1], // 挿入
						dp[i - 1][j - 1], // 置換
					) + 1;
			}
		}
	}

	return dp[a.length][b.length];
}

// 定期的に本文情報を再取得する。
// 他のユーザーが大幅に変更していた場合は警告を表示し、リロードを促す。
// TODO: 編集不可状態とする。
function checkOtherUsersUpdate() {
	api
		.getEntryByDynamicPath({
			path: entry.path,
		})
		.then((data) => {
			// 本文が短いときは消えてもダメージ少ないので無視
			if (data.body && data.body.length > 100 && !isDirty) {
				const threshold = Math.max(body.length, data.body.length) * 0.1; // 10%以上の変更で判定
				const editDistance = getEditDistance(body, data.body);
				if (editDistance > threshold) {
					if (
						confirm(
							`他のユーザーが大幅に変更しました。リロードしてください。 ${editDistance} > ${threshold}`,
						)
					) {
						location.reload();
					}
				}
			}
		});
}

onMount(() => {
	// get page titles
	setTimeout(async () => {
		pageTitles = await api.getAllEntryTitles();
	}, 0);

	document.addEventListener("visibilitychange", () => {
		checkOtherUsersUpdate();
	});
});

function selectIfPlaceholder(target: HTMLInputElement) {
	if (/^2\d+$/.test(target.value)) {
		target.select(); // 入力値を全選択
	}
}
</script>

<div class="parent">
    <div class="container {entry.visibility === 'private' ? 'private' : ''}">
        <div class="left-pane">
            <form class="form">
                <div class="title-container">
                    <input
                            id="title"
                            name="title"
                            type="text"
                            class="input"
                            bind:value={title}
                            onfocus={(event) => selectIfPlaceholder(event.target as HTMLInputElement)}
                            oninput={handleTitleInput}
                            required
                    />
                    <button class="visibility-icon" onclick={(event) => toggleVisibility(event)}
                    >{visibility === 'private' ? '🔒️' : '🌍'}</button
                    >
                </div>

                <div class="editor">
                    <input type="hidden" name="body" bind:value={body}/>
                    <MarkdownEditor
                            initialContent=""
                            content={body}
                            onUpdateText={(content) => {
                                if (content !== body) {
                                    body = content;
                                    handleInputBody(); // エディタ更新時にデバウンスされた更新をトリガー
                                }
                            }}
                            existsEntryByTitle={(title) => {
                                return !!links[title.toLowerCase()];
                            }}
                            onClickEntry={(title) => {
                                if (links[title.toLowerCase()]) {
                                    location.href = '/admin/entry/' + links[title.toLowerCase()];
                                } else {
                                    createNewEntry(title);
                                }
                            }}
                            onSave={() => {
                                handleUpdateBody();
                            }}
                            {pageTitles}
                    />
                </div>
            </form>
        </div>

        <div class="right-pane">
            <div class="button-container">
                <button type="submit" class="delete-button" onclick={handleDelete}> Delete</button>
                <button type="submit" class="regenerate-button" onclick={handleRegenerateEntryImage}> Regenerate entry_image</button>
            </div>

            <!-- link to the user side page -->
            {#if visibility === 'public'}
                <div class="link-container">
                    <a href="/entry/{entry.path}" class="link">Go to User Side Page</a>
                </div>
            {/if}

            {#if updatedMessage !== ''}
                <div class="updated-message">
                    {updatedMessage}
                </div>
            {/if}
        </div>
    </div>

    {#if message}
        <div class="popup-message {messageType}">
            <p>{message}</p>
        </div>
    {/if}

    <div class="link-pallet">
        <LinkPallet {linkPallet} />
    </div>
</div>

<style>
    .link-pallet {
        margin: auto;
        max-width: 1200px;
    }

    .container {
        display: flex;
        flex-wrap: wrap;
        max-width: 1200px;
        margin: auto;
    }

    .container.private {
        background-color: #e5e7eb;
    }

    .left-pane {
        flex: 1;
        min-width: 300px;
        max-width: 800px;
    }

    .right-pane {
        margin-left: 1rem; /* Add some space between the panes */
        max-width: 300px;

        .link-container {
            margin-top: 80px;
        }

        button, a {
            display: block;
            border-radius: 0.375rem;
            border: 1px solid antiquewhite;
            padding: 0.5rem 1rem;
            margin-top: 4px;
        }

        .delete-button {
            background-color: #ef4444;
            color: white;

            &:hover {
                background-color: #dc2626;
            }
        }

        .link {
            background-color: #10b981;
            color: white;
            text-decoration: none;

            &:hover {
                text-decoration: underline;
            }
        }
    }

    .form {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        padding: 1rem;
        max-width: 800px;
    }

    .input {
        width: calc(100% - 2rem); /* Adjust width to make space for the icon */
        border-radius: 0.375rem;
        border: 1px solid #d1d5db;
        padding: 0.5rem;
    }

    .title-container {
        display: flex;
        align-items: center;
    }

    .visibility-icon {
        border: 0;
        background-color: transparent;
        margin-left: 0.5rem;
        font-size: 1.5rem;
        cursor: pointer;
    }

    .editor {
        border: 1px solid #d1d5db;
        border-radius: 0.25rem;
        height: 400px;
        overflow-y: scroll;
    }

    .popup-message {
        position: fixed;
        bottom: 1rem;
        right: 1rem;
        padding: 1rem;
        border-radius: 0.375rem;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        z-index: 1000;
    }

    .popup-message.success {
        background-color: #10b981;
        color: white;
    }

    .popup-message.error {
        background-color: #ef4444;
        color: white;
    }

    @media (max-width: 600px) {
        .container {
            flex-direction: column;
        }

        .right-pane {
            margin-left: 0;
            margin-top: 1rem; /* Add some space between the panes */
        }
    }

    .updated-message {
        color: #10b981;
    }
</style>

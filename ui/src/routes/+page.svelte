<script lang="ts">
  import { faGear } from "@fortawesome/free-solid-svg-icons";
  import Icon from "$lib/icon.svelte";
  import ActionButton from "$lib/actionButton.svelte";
  import type { ViewPatient, ViewSetting } from "$lib/types";
  import Fuse, { type FuseResult } from "fuse.js";
  import {
    settingsFolders,
    updatePatient,
    folders,
    type Folder,
    patientsStore,
    updatePatients,
  } from "$lib/DataService";
  import FolderList from "./Folder.svelte";
  import Patient from "./Patient.svelte";
  import { dndzone } from "svelte-dnd-action";
  import { flip } from "svelte/animate";
  let patients: Array<ViewPatient> = [];
  $: {
    patients = $patientsStore.filter((p: ViewPatient) => p.folder == -1);
  }
  // var filteredPatient: any[] = [];
  var fuse: any = null;
  var searchQuery = "";

  function handleDnd(e: any) {
    patients = e.detail.items;
    patients.forEach((p: ViewPatient, index: number) => {
      p.folder = -1;
      p.indexFolder = index;
    });
  }

  function handleDndFinalize(e: any) {
    handleDnd(e);
    updatePatients(patients);
  }

  // Function to wrap matched parts in <mark> tags for highlighting
  function generateHighlightedString(text: any, indices: any) {
    let highlighted = "";
    let lastIndex = 0;

    for (const [start, end] of indices) {
      // Add text before the match
      highlighted += text.slice(lastIndex, start);

      // Wrap the matched text in <mark> tags
      highlighted += `<mark>${text.slice(start, end + 1)}</mark>`;

      // Update the last index
      lastIndex = end + 1;
    }

    // Add remaining text after the final match
    highlighted += text.slice(lastIndex);

    return highlighted;
  }
</script>

<header class="max-w-7xl mx-auto bg-white shadow flex items-center p-2">
  <input
    name="search"
    type="search"
    placeholder="Search"
    class="border rounded px-3 mr-4 py-2 flex-grow focus:outline-none focus:ring"
    bind:value={searchQuery}
  />
  <a
    href="/settings"
    class="text-gray-600 hover:bg-blue-100 hover:rounded-lg p-2"
  >
    <Icon icon={faGear} />
  </a>
</header>

<ActionButton href="/patient/new/edit" />
<div class="max-w-7xl mx-auto">
  <section
    class="min-h-8"
    use:dndzone={{ items: patients, flipDurationMs: 300 }}
    on:consider={handleDnd}
    on:finalize={handleDndFinalize}
  >
    {#each patients as p (p.id)}
      <!-- Patient List -->
      <div animate:flip={{ duration: 300 }}>
        <Patient patient={p} />
      </div>
    {/each}
  </section>
  {#each $folders as folder}
    <FolderList {folder} />
  {/each}
</div>

<style>
  :global(mark) {
    background-color: transparent;
    color: blue;
    font-weight: bold;
  }
</style>

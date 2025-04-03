<script lang="ts">
  import {
    faAngleDown,
    faAngleRight,
    faFolder,
    faGear,
  } from "@fortawesome/free-solid-svg-icons";
  import Icon from "$lib/icon.svelte";
  import ActionButton from "$lib/actionButton.svelte";
  import type { ViewPatient, ViewSetting } from "$lib/types";
  import Fuse, { type FuseResult } from "fuse.js";
  import { dndzone } from "svelte-dnd-action";
  import {
    patientsStore,
    settingsFolders,
    updatePatient,
  } from "$lib/DataService";
  import { flip } from "svelte/animate";
  import Patient from "./Patient.svelte";
  import { writable, type Writable } from "svelte/store";

  var patients = patientsStore;
  // var filteredPatient: any[] = [];
  var fuse: any = null;
  var searchQuery = "";

  interface Folder {
    collapsed: boolean;
    section: any | null;
    folder: ViewSetting;
    patients: Array<ViewPatient>;
  }
  var folders = writable(new Array<Folder>());
  settingsFolders.subscribe(($items) => {
    folders.set(
      $items.map((item) => {
        return {
          collapsed: false,
          section: null,
          folder: item,
          patients: new Array<ViewPatient>(),
        };
      }),
    );
  });

  function handleDnd(e: any) {
    e.detail.items.forEach((p: ViewPatient, i: number) => {
      p.indexFolder = i;
    });
    patients.set(e.detail.items);
  }

  function handleDndFolder(folder: Folder): (e: any) => void {
    return (e: any) => {
      e.detail.items.forEach((p: ViewPatient, i: number) => {
        p.indexFolder = i;
      });
      folder.patients = e.detail.items;
      folders.update((f) => {
        return f.map((fItem) => {
          if (fItem.folder.id === folder.folder.id) {
            return folder;
          }
          return fItem;
        });
      });
      // console.log(folder.folder.value, e.detail.items);
    };
  }

  // Function to generate highlighted matches
  function highlightMatches(result: any) {
    const highlightedFields: any = {};
    for (const match of result.matches) {
      const { key, value, indices } = match;

      if (!highlightedFields[key]) {
        highlightedFields[key] = generateHighlightedString(value, indices);
      }
    }
    return highlightedFields;
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
    use:dndzone={{ items: $patients, flipDurationMs: 50 }}
    on:consider={handleDnd}
    on:finalize={handleDnd}
  >
    <!-- Patient List -->
    {#each $patients as p (p.id)}
      <div animate:flip={{ duration: 200 }}>
        <Patient patient={p} />
      </div>
    {/each}
  </section>
  {#each $folders as folder}
    <div class="m-2 p-2 flex items-center">
      <Icon class="m-4" icon={faFolder} />{folder.folder.value}
    </div>
    <section
      class="ml-6 min-h-8"
      use:dndzone={{
        items: folder.patients,
        flipDurationMs: 50,
      }}
      on:consider={handleDndFolder(folder)}
      on:finalize={handleDndFolder(folder)}
    >
      {#each folder.patients as p (p.id)}
        <div animate:flip={{ duration: 200 }}>
          <Patient patient={p} />
        </div>
      {/each}
      {#if folder.patients.length === 0}
        <div class="flex items-center space-x-4 border-b border-gray-400">
          <p class="p-4 text-gray-500">No patients in this folder</p>
        </div>
      {/if}
    </section>
  {/each}
</div>

<style>
  :global(mark) {
    background-color: transparent;
    color: blue;
    font-weight: bold;
  }
</style>

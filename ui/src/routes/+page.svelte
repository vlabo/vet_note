<script lang="ts">
  import {
    faGear,
    // faMicrochip,
    // faPhone,
    faUser,
    faUserSlash,
  } from "@fortawesome/free-solid-svg-icons";
  import Icon from "$lib/icon.svelte";
  import ActionButton from "$lib/actionButton.svelte";
  import type { ViewPatient } from "$lib/types";
  import Fuse, { type FuseResult } from "fuse.js";
  import { patientsStore } from "$lib/DataService";
  import { onMount } from "svelte";

  var patients: ViewPatient[];
  var filteredPatient: any[] = [];
  var fuse: any = null;
  var searchQuery = "";
  onMount(() => {
    patientsStore.subscribe((storePatients) => {
      fuse = new Fuse(storePatients, {
        keys: ["name", "owner", "ownerPhone", "chipId"],
        includeMatches: true,
      });
      filteredPatient = storePatients;
      patients = storePatients;
    });
  });

  function onSearch() {
    if (!searchQuery) {
      filteredPatient = patients;
    } else {
      const result = fuse.search(searchQuery);

      // Map over the search results to include the highlighted parts
      filteredPatient = result.map((p: FuseResult<ViewPatient>) => {
        return {
          ...p.item, // Original patient data
          highlightedFields: highlightMatches(p), // Extract highlighted matches
        };
      });
    }
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
    on:input={onSearch}
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
  {#each filteredPatient as p}
    <a
      class="max-w-7xl mx-auto bg-white mb-2 cursor-pointer overflow-hidden"
      href={"/patient/" + p.id}
    >
      <div class="p-4 flex items-center space-x-4 border-b border-gray-400">
        <!-- Type -->
        <div class="w-1/5 font-medium text-gray-700">{p.type}</div>
        <div class="flex-1 flex items-center text-gray-800">
          <div class="ml-2">
            {@html p.highlightedFields?.name || p.name}
          </div>
        </div>
        <!-- Patient Owner -->
        {#if p.owner}
          <div class="flex-1 flex items-center text-gray-800">
            <Icon icon={faUser} />
            <div class="ml-2">
              {@html p.highlightedFields?.owner || p.owner}
            </div>
          </div>
        {:else}
          <div class="flex-1 flex items-center text-gray-800">
            <Icon icon={faUserSlash} />
          </div>
        {/if}
        <!-- Optional Chip ID
        {#if p.chipId}
          <div class="flex-1 flex items-center text-gray-800">
            <Icon icon={faMicrochip} />
            <div class="ml-2">
              {@html p.highlightedFields?.chipId || p.chipId}
            </div>
          </div>
        {/if} -->
        <!-- Optional Phone
        {#if p.ownerPhone}
          <div class="flex-1 flex items-center text-gray-800">
            <Icon icon={faPhone} />
            <div class="ml-2">
              {@html p.highlightedFields?.ownerPhone || p.ownerPhone}
            </div>
          </div>
        {/if} -->
      </div>
    </a>
  {/each}
</div>

<style>
  :global(mark) {
    background-color: transparent;
    color: blue;
    font-weight: bold;
  }
</style>

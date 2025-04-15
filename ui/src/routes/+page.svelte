<script lang="ts">
  import { faGear } from "@fortawesome/free-solid-svg-icons";
  import Icon from "$lib/icon.svelte";
  import ActionButton from "$lib/actionButton.svelte";
  import type { ViewPatient, ViewSetting } from "$lib/types";
  import Fuse, { type FuseResult } from "fuse.js";
  import { folders, patientsStore, updatePatients } from "$lib/DataService";
  import FolderList from "./Folder.svelte";
  import Patient from "./Patient.svelte";
  import { dndzone } from "svelte-dnd-action";
  import { flip } from "svelte/animate";
  import type { ViewPatientHighlighted } from "$lib/Utils";
  let patients: Array<ViewPatient> | Array<ViewPatientHighlighted> = [];
  var fuse: Fuse<ViewPatient> | null = null;
  let filteredPatients: Array<any> = [];

  $: {
    patients = $patientsStore.filter((p: ViewPatient) => p.folder == -1);
    fuse = new Fuse(patients, {
      includeMatches: true,
      keys: ["name", "owner"],
    });
  }

  $: displayPatients = searchQuery.length ? filteredPatients : patients;
  // var filteredPatient: any[] = [];
  var searchQuery = "";

  function handleDnd(e: any) {
    displayPatients = e.detail.items;
    displayPatients.forEach((p: ViewPatient, index: number) => {
      p.folder = -1;
      p.indexFolder = index;
    });
  }

  function handleDndFinalize(e: any) {
    handleDnd(e);
    updatePatients(displayPatients);
  }

  function onSearch() {
    if (searchQuery.length == 0) {
      displayPatients = patients;
      displayPatients.forEach((p: any) => {
        p.highlightedFields = {};
      });
      return;
    }
    let result = fuse!.search(searchQuery);
    filteredPatients = result.map((r) => {
      let item = r.item;
      if (r.matches) {
        let patient: ViewPatientHighlighted = {
          ...item,
          highlightedFields: {},
        };

        for (const match of r.matches) {
          if (match.key === "name" && match.indices) {
            patient.highlightedFields.name = generateHighlightedString(
              patient.name,
              match.indices,
            );
          }
          if (match.key === "owner" && match.indices) {
            patient.highlightedFields.owner = generateHighlightedString(
              patient.owner,
              match.indices,
            );
          }
        }
        return patient;
      }
    });
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
  <section
    class="min-h-8"
    use:dndzone={{ items: displayPatients, flipDurationMs: 100 }}
    on:consider={handleDnd}
    on:finalize={handleDndFinalize}
  >
    {#each displayPatients as p (p.id)}
      <!-- Patient List -->
      <div animate:flip={{ duration: 100 }}>
        <Patient patient={p} />
      </div>
    {/each}
  </section>
  {#each $folders as folder}
    <FolderList {folder} {searchQuery} />
  {/each}
</div>

<style>
  :global(mark) {
    background-color: transparent;
    color: blue;
    font-weight: bold;
  }
</style>

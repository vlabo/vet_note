<script lang="ts">
  import Icon from "$lib/icon.svelte";
  import {
    faAngleRight,
    faAngleDown,
    faFolder,
  } from "@fortawesome/free-solid-svg-icons";
  import Patient from "./Patient.svelte";
  import { updatePatients, type Folder } from "$lib/DataService";
  import { dndzone, SOURCES, TRIGGERS } from "svelte-dnd-action";
  import { flip } from "svelte/animate";
  import type { ViewPatient } from "$lib/types";
  import Fuse from "fuse.js";
  import type { ViewPatientHighlighted } from "$lib/Utils";

  export var folder: Folder;
  export var collapsed: boolean = true;
  export var searchQuery: string = "";

  let patients = new Array<ViewPatient>();
  let filteredPatients: Array<any> = [];
  var fuse: Fuse<ViewPatient> | null = null;
  var displayPatients: Array<any> = patients;
  let dragDisabled: boolean = true;

  $: {
    patients = folder.patients;
    fuse = new Fuse(patients, {
      includeMatches: true,
      keys: ["name", "owner"],
    });
    displayPatients = searchQuery.length ? filteredPatients : patients;
  }

  function onSearch(searchQuery: string) {
    if (searchQuery.length == 0) {
      displayPatients.forEach((p: any) => {
        p.highlightedFields = {};
      });
      return;
    }

    // Search with Fuse.js
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
    displayPatients = filteredPatients;
  }

  $: onSearch(searchQuery);

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

  function handleDnd(e: any) {
    const {
      items: _,
      info: { source, trigger },
    } = e.detail;

    displayPatients = e.detail.items;
    displayPatients.forEach((p: ViewPatient, index: number) => {
      p.folder = folder.setting.id!;
      p.indexFolder = index;
    });
    if (source === SOURCES.KEYBOARD && trigger === TRIGGERS.DRAG_STOPPED) {
      dragDisabled = true;
    }
  }

  function handleDndFinalize(e: any) {
    const {
      items: _,
      info: { source },
    } = e.detail;
    handleDnd(e);
    updatePatients(displayPatients);
    if (source === SOURCES.POINTER) {
      dragDisabled = true;
    }
  }

  function startDrag(e: any) {
    // preventing default to prevent lag on touch devices (because of the browser checking for screen scrolling)
    e.preventDefault();
    dragDisabled = false;
  }

  function handleKeyDown(e: any) {
    if ((e.key === "Enter" || e.key === " ") && dragDisabled) {
      dragDisabled = false;
    }
  }
</script>

<div class="mr-2 pr-2 flex items-center">
  <button on:click={() => (collapsed = !collapsed)}>
    {#if collapsed}
      <Icon class="m-4" icon={faAngleRight} />
    {:else}
      <Icon class="m-4" icon={faAngleDown} />
    {/if}
  </button>
  <Icon class="mr-4" icon={faFolder} />{folder.setting.value}
</div>
<div>
  <section
    class="ml-6 min-h-8"
    use:dndzone={{ items: displayPatients, dragDisabled, flipDurationMs: 100 }}
    on:consider={handleDnd}
    on:finalize={handleDndFinalize}
    hidden={collapsed}
  >
    {#each displayPatients as p (p.id)}
      <!-- Patient List -->
      <div animate:flip={{ duration: 100 }}>
        <Patient patient={p} {dragDisabled} {startDrag} {handleKeyDown} />
      </div>
    {/each}
    {#if displayPatients.length == 0}
      <div class="ml-6">
        <p class="text-gray-500">No patients in this folder.</p>
      </div>
    {/if}
  </section>
</div>

<script lang="ts">
  import Icon from "$lib/icon.svelte";
  import {
    faAngleRight,
    faAngleDown,
    faFolder,
  } from "@fortawesome/free-solid-svg-icons";
  import Patient from "./Patient.svelte";
  import { updatePatients, type Folder } from "$lib/DataService";
  import { dndzone } from "svelte-dnd-action";
  import { flip } from "svelte/animate";
  import type { ViewPatient } from "$lib/types";

  export var folder: Folder;
  export var collapsed: boolean = true;
  let patients = new Array<ViewPatient>();
  $: {
    patients = folder.patients;
    // patients.sort((a, b) => a.indexFolder! - b.indexFolder!);
  }

  function handleDnd(e: any) {
    console.log(e.detail);
    // patients = e.detail.items;
    patients = e.detail.items.map((p: ViewPatient, index: number) => {
      p.folder = folder.setting.id;
      p.indexFolder = index;
      console.log(index, p);
      return p;
    });
    console.log(patients);
  }

  function handleDndFinalize(e: any) {
    handleDnd(e);
    updatePatients(patients);
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
    use:dndzone={{ items: patients, flipDurationMs: 300 }}
    on:consider={handleDnd}
    on:finalize={handleDndFinalize}
    hidden={collapsed}
  >
    {#each patients as p (p.id)}
      <!-- Patient List -->
      <div animate:flip={{ duration: 300 }}>
        <Patient patient={p} />
      </div>
    {/each}
    {#if patients.length == 0}
      <div class="ml-6">
        <p class="text-gray-500">No patients in this folder.</p>
      </div>
    {/if}
  </section>
</div>

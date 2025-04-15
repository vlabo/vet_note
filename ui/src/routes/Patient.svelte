<script lang="ts">
  import Icon from "$lib/icon.svelte";
  import {
    faMicrochip,
    faPhone,
    faUser,
    faUserSlash,
  } from "@fortawesome/free-solid-svg-icons";
  import { onMount } from "svelte";

  export var patient: any;
  let elem: HTMLElement;
  onMount(() => {
    if (elem) {
      elem.dataset.patientId = String(patient.id);
    }
  });
</script>

<div bind:this={elem}>
  <a
    class="max-w-7xl mx-auto mb-2 cursor-pointer overflow-hidden"
    href={"/patient/" + patient.id}
  >
    <div class="p-4 flex items-center space-x-4 border-b border-gray-400">
      <!-- Type -->
      <div class="w-1/5 font-medium text-gray-700">{patient.type}</div>
      <div class="flex-1 flex items-center text-gray-800">
        <div class="ml-2">
          {@html patient.highlightedFields?.name
            ? patient.highlightedFields.name
            : patient.name}
        </div>
      </div>
      <!-- Patient Owner -->
      {#if patient.owner}
        <div class="flex-1 flex items-center text-gray-800">
          <Icon icon={faUser} />
          <div class="ml-2">
            {@html patient.highlightedFields?.owner
              ? patient.highlightedFields.owner
              : patient.owner}
          </div>
        </div>
      {:else}
        <div class="flex-1 flex items-center text-gray-800">
          <Icon icon={faUserSlash} />
        </div>
      {/if}
    </div>
  </a>
</div>

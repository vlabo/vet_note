<script lang="ts">
  import Icon from "$lib/icon.svelte";
  import {
    faBars,
    faMicrochip,
    faPhone,
    faUser,
    faUserSlash,
  } from "@fortawesome/free-solid-svg-icons";
  import { onMount } from "svelte";

  export var patient: any;
  export var dragDisabled: boolean = true;
  export var startDrag = (_: any) => {};
  export var handleKeyDown = (_: any) => {};
</script>

<div class="w-full flex items-center border-b border-gray-400">
  <div
    tabindex="0"
    role="button"
    aria-label="drag-handle"
    style={dragDisabled ? "cursor: grab" : "cursor: grabbing"}
    on:mousedown={startDrag}
    on:touchstart={startDrag}
    on:keydown={handleKeyDown}
  >
    <Icon icon={faBars} class="py-3 pl-2 pr-5 drag-handle" />
  </div>
  <a
    class="flex-1 mx-auto max-w-7xl mb-2 cursor-pointer overflow-hidden"
    href={"/patient/" + patient.id}
  >
    <div class="p-4 flex items-center space-x-4">
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
